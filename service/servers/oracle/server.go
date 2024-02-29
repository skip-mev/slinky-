package oracle

import (
	"context"
	"fmt"
	"github.com/CosmWasm/wasmd/pkg/sync"
	"github.com/CosmWasm/wasmd/service/servers/oracle/types"
	"github.com/CosmWasm/wasmd/x/slpp/service"
	"net/http"
	"strings"
	"time"

	gateway "github.com/cosmos/gogogateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const DefaultServerShutdownTimeout = 3 * time.Second

type OracleServer struct { //nolint
	service.UnimplementedOracleServiceServer

	// underlying grpc-server -- serves all grpc requests
	grpcSrv *grpc.Server

	// grpc-gateway mux -- serves all http grpc proxy requests
	gatewayMux *runtime.ServeMux

	// underlying http server
	httpSrv *http.Server

	// closer to handle graceful closures from multiple go-routines
	*sync.Closer

	// logger to log incoming requests
	logger *zap.Logger

	// vote extension data to return
	data []byte
}

// NewOracleServer returns a new instance of the OracleServer, given an implementation of the Oracle interface.
func NewOracleServer(logger *zap.Logger, data []byte) *OracleServer {
	logger = logger.With(zap.String("server", "oracle"))

	os := &OracleServer{
		logger: logger,
		data:   data,
	}
	os.Closer = sync.NewCloser().WithCallback(func() {
		// if the server has been started, close it
		if os.httpSrv != nil {
			ctx, cf := context.WithTimeout(context.Background(), DefaultServerShutdownTimeout)
			os.httpSrv.Shutdown(ctx) // close HTTP server backing GRPC-gateway
			os.grpcSrv.Stop()        // close GRPC server serving listeners that have been routed to GRPC server
			cf()
		}
	})

	return os
}

// routeRequest determines if the incoming http request is a grpc or http request and routes to the proper handler.
func (os *OracleServer) routeRequest(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.HasPrefix(
		r.Header.Get("Content-Type"), "application/grpc") {

		os.grpcSrv.ServeHTTP(w, r)
	} else {
		os.gatewayMux.ServeHTTP(w, r)
	}
}

// StartServer starts the oracle gRPC server on the given host and port. The server is killed on any errors from the listener, or if ctx is cancelled.
// This method returns an error via any failure from the listener. This is a blocking call, i.e. until the server is closed or the server errors,
// this method will block.
func (os *OracleServer) StartServer(ctx context.Context, host, port string) error {
	serverEndpoint := fmt.Sprintf("%s:%s", host, port)
	os.httpSrv = &http.Server{
		Addr:              serverEndpoint,
		ReadHeaderTimeout: DefaultServerShutdownTimeout,
	}
	// create grpc server
	os.grpcSrv = grpc.NewServer()
	// register oracle server
	service.RegisterOracleServiceServer(os.grpcSrv, os)

	// register the grpc-gateway
	// it handles the http request and dials the server endpoint with the grpc request
	os.gatewayMux = runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &gateway.JSONPb{
			EmitDefaults: true,
			Indent:       "",
			OrigName:     true,
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := types.RegisterOracleHandlerFromEndpoint(ctx, os.gatewayMux, serverEndpoint, opts)
	if err != nil {
		return err
	}

	router := http.NewServeMux()
	router.HandleFunc("/", os.routeRequest)

	os.httpSrv.Handler = h2c.NewHandler(router, &http2.Server{})

	eg, ctx := errgroup.WithContext(ctx)

	// listen for ctx cancellation
	eg.Go(func() error {
		// if the context is closed, close the server + oracle
		<-ctx.Done()
		os.logger.Info("context cancelled, closing oracle")

		_ = os.Close()
		return nil
	})

	// start the server
	eg.Go(func() error {
		// serve, and return any errors
		os.logger.Info(
			"starting grpc server",
			zap.String("host", host),
			zap.String("port", port),
		)

		err = os.httpSrv.ListenAndServe()
		if err != nil {
			return fmt.Errorf("[grpc server]: error serving: %w", err)
		}

		return nil
	})

	// wait for everything to finish
	return eg.Wait()
}

// Prices calls the underlying oracle's implementation of GetPrices. It defers to the ctx in the request, and errors if the context is cancelled
// for any reason, or if the oracle errors.
func (os *OracleServer) VoteExtensionData(ctx context.Context, req *service.VoteExtensionDataRequest) (*service.VoteExtensionDataResponse, error) {
	// check that the request is non-nil
	if req == nil {
		return nil, ErrNilRequest
	}

	return &service.VoteExtensionDataResponse{
		Data: os.data,
	}, nil
}

// Close closes the underlying oracle server, and blocks until all open requests have been satisfied.
func (os *OracleServer) Close() error {
	// close + close server if necessary
	os.Closer.Close()
	return nil
}

// Done returns a channel that is closed when the oracle server is closed.
func (os *OracleServer) Done() <-chan struct{} {
	return os.Closer.Done()
}
