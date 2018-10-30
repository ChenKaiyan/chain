// Package dpor implements the dpor consensus engine.
package dpor

import (
	"sync"

	"bitbucket.org/cpchain/chain/configs"
	"bitbucket.org/cpchain/chain/consensus"
	"bitbucket.org/cpchain/chain/ethdb"
	"github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru"
)

const (
	inmemorySnapshots  = 1000 // Number of recent vote snapshots to keep in memory
	inmemorySignatures = 1000 // Number of recent block signatures to keep in memory

	pctA = 2
	pctB = 3 // only when n > 2/3 * N, accept the block
)

// Dpor is the proof-of-reputation consensus engine proposed to support the
// cpchain testnet.
type Dpor struct {
	dh     dporHelper
	db     ethdb.Database      // Database to store and retrieve Snapshot checkpoints
	config *configs.DporConfig // Consensus engine configuration parameters

	recents    *lru.ARCCache // Snapshots for recent block to speed up reorgs
	signatures *lru.ARCCache // Signatures of recent blocks to speed up mining

	signedBlocks map[uint64]common.Hash // record signed blocks

	signer common.Address // Ethereum address of the signing key
	signFn SignerFn       // Signer function to authorize hashes with

	contractCaller          *consensus.ContractCaller         // contractCaller is used to handle contract related calls
	committeeNetworkHandler consensus.CommitteeNetworkHandler // committeeHandler handles all related dials and calls

	lock sync.RWMutex // Protects the signer fields defined above
}

// New creates a Dpor proof-of-reputation consensus engine with the initial
// signers set to the ones provided by the user.
func New(config *configs.DporConfig, db ethdb.Database) *Dpor {

	// Set any missing consensus parameters to their defaults
	conf := *config
	if conf.Epoch == 0 {
		conf.Epoch = uint64(epochLength)
	}
	if conf.View == 0 {
		conf.View = uint64(viewLength)
	}

	// Allocate the Snapshot caches and create the engine
	recents, _ := lru.NewARC(inmemorySnapshots)
	signatures, _ := lru.NewARC(inmemorySignatures)
	signedBlocks := make(map[uint64]common.Hash)

	return &Dpor{
		dh:           &defaultDporHelper{&defaultDporUtil{}},
		config:       &conf,
		db:           db,
		recents:      recents,
		signatures:   signatures,
		signedBlocks: signedBlocks,
	}
}

// SetContractCaller sets dpor.contractCaller
func (d *Dpor) SetContractCaller(contractCaller *consensus.ContractCaller) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.contractCaller = contractCaller
	return nil
}

// SetCommitteeNetworkHandler sets dpor.committeeNetworkHandler
func (d *Dpor) SetCommitteeNetworkHandler(committeeNetworkHandler consensus.CommitteeNetworkHandler) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.committeeNetworkHandler = committeeNetworkHandler
	return nil
}
