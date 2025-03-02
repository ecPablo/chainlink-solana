package client

import (
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupLocalSolNode_SimultaneousNetworks(t *testing.T) {
	// run two networks
	network0 := SetupLocalSolNode(t)
	network1 := SetupLocalSolNode(t)

	account := solana.NewWallet()
	pubkey := account.PublicKey()

	// client configs
	requestTimeout := 5 * time.Second
	lggr := logger.Test(t)
	cfg := config.NewConfig(db.ChainCfg{}, lggr)

	// check & fund address
	checkFunded := func(t *testing.T, url string) {
		// create client
		c, err := NewClient(url, cfg, requestTimeout, lggr)
		require.NoError(t, err)

		// check init balance
		bal, err := c.Balance(pubkey)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), bal)

		FundTestAccounts(t, []solana.PublicKey{pubkey}, url)

		// check end balance
		bal, err = c.Balance(pubkey)
		assert.NoError(t, err)
		assert.Equal(t, uint64(100_000_000_000), bal) // once funds get sent to the system program it should be unrecoverable (so this number should remain > 0)
	}

	checkFunded(t, network0)
	checkFunded(t, network1)
}
