package petri_test

import (
	"github.com/CosmWasm/wasmd/petri"
	"github.com/stretchr/testify/suite"
	"testing"
)

// runs the slinky integration tests
func TestSlinkyIntegration(t *testing.T) {
	chainCfg, err := petri.GetChainConfig()
	if err != nil {
		t.Fatal(err)
	}
	suite.Run(t, petri.NewSlinkyIntegrationSuite(&chainCfg))
}
