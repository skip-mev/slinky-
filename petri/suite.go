package petri

import (
	"go.uber.org/zap"
	"github.com/stretchr/testify/suite"
	petritypes "github.com/skip-mev/petri/core/v2/types"
	"context"
)

const (
	envKeepAlive = "PETRI_LOAD_TEST_KEEP_ALIVE"
	marketsPath = "./fixtures/markets.json"
)

// SlinkyIntegrationSuite is a test-suite used to spin up load-tests of arbitrary size for dydx nodes
type SlinkyIntegrationSuite struct {
	suite.Suite

	logger *zap.Logger

	spec *petritypes.ChainConfig

	chain petritypes.ChainI
}

func NewSlinkyIntegrationSuite(spec *petritypes.ChainConfig) *SlinkyIntegrationSuite {
	return &SlinkyIntegrationSuite{
		spec: spec,
	}
}

func (s *SlinkyIntegrationSuite) SetupSuite() {
	// create the logger
	var err error
	s.logger, err = zap.NewDevelopment()
	s.Require().NoError(err)

	// create the chain
	s.chain, err = GetChain(context.Background(), s.logger, *s.spec)
	s.Require().NoError(err)

	//initialize the chain
	err = s.chain.Init(context.Background())
	s.Require().NoError(err)
}
