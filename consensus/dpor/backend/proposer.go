package backend

import (
	"bitbucket.org/cpchain/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

// ReceiveMinedPendingBlock receives a block to add to pending block channel
func (ph *Handler) ReceiveMinedPendingBlock(block *types.Block) error {
	select {
	case ph.pendingBlockCh <- block:
		err := ph.knownBlocks.AddBlock(block)
		if err != nil {
			return err
		}

		return nil
	}
}

// UpdateRemoteValidators updates handler.dialer.remoteValidators
// this is called if local peer is a future proposer
func (ph *Handler) UpdateRemoteValidators(term uint64, validators []common.Address) error {
	return ph.dialer.UpdateRemoteValidators(term, validators)
}

// DialAllRemoteValidators dials remote validators
// this is called if local peer is a future proposer
func (ph *Handler) DialAllRemoteValidators(term uint64) error {
	return ph.dialer.DialAllRemoteValidators(term)
}
