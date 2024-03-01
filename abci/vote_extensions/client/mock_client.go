package client

import (
	"context"
	"github.com/CosmWasm/wasmd/x/slpp/service"
)

type MultiOracleClientMock struct {}

func (mc MultiOracleClientMock) VoteExtensionData(ctx context.Context, avsID uint64, req *service.VoteExtensionDataRequest) (*service.VoteExtensionDataResponse, error) {
	return &service.VoteExtensionDataResponse{
		Data: []byte("mock data"),
	}, nil
}
