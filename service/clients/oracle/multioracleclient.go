package oracle

import (
	"context"
	"fmt"
	voteextensions "github.com/CosmWasm/wasmd/abci/vote_extensions"
	"github.com/CosmWasm/wasmd/x/slpp/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type multiOracleClientImpl struct {
	oracleByAVSID map[uint64]service.OracleServiceClient
}

func NewMultiOracleClientFromConfig(ctx context.Context, oracleConfigs []OracleClientConfig) (voteextensions.MultiOracleClient, error) {
	var oracleClients []service.OracleServiceClient
	var avsIDs []uint64
	for _, conf := range oracleConfigs {
		conn, err := grpc.DialContext(
			ctx,
			conf.OracleAddress,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, err
		}
		oracleClients = append(oracleClients, service.NewOracleServiceClient(conn))
		avsIDs = append(avsIDs, conf.AVSID)
	}
	return NewMultiOracleClient(oracleClients, avsIDs), nil
}

func NewMultiOracleClient(oracleClients []service.OracleServiceClient, avsIDs []uint64) voteextensions.MultiOracleClient {
	oracleByAVSID := make(map[uint64]service.OracleServiceClient)
	for i := range oracleClients {
		oracleByAVSID[avsIDs[i]] = oracleClients[i]
	}
	return &multiOracleClientImpl{oracleByAVSID: oracleByAVSID}
}

func (m *multiOracleClientImpl) VoteExtensionData(ctx context.Context, avsID uint64, req *service.VoteExtensionDataRequest) (*service.VoteExtensionDataResponse, error) {
	oracleClient, ok := m.oracleByAVSID[avsID]
	if !ok {
		return nil, fmt.Errorf("oracle client not found for avsID")
	}
	return oracleClient.VoteExtensionData(ctx, req)
}
