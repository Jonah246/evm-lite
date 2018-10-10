package babble

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mosaicnetworks/babble/src/hashgraph"
	"github.com/mosaicnetworks/evm-lite/src/service"
	"github.com/mosaicnetworks/evm-lite/src/state"
)

//Handler implements Babble's ProxyHandler interface
type Handler struct {
	service *service.Service
	state   *state.State
}

//NewHandler create a new Handler
func NewHandler(state *state.State, service *service.Service) *Handler {
	return &Handler{
		service: service,
		state:   state,
	}
}

/*******************************************************************************
IMPLEMENT PROXYHANDLER
*******************************************************************************/

//CommitHandler applies the block's transactions to the state
func (h *Handler) CommitHandler(block hashgraph.Block) ([]byte, error) {

	blockHashBytes, err := block.Hash()
	blockHash := common.BytesToHash(blockHashBytes)

	for i, tx := range block.Transactions() {
		if err := h.state.ApplyTransaction(tx, i, blockHash); err != nil {
			return []byte{}, err
		}
	}

	hash, err := h.state.Commit()
	if err != nil {
		return []byte{}, err
	}

	return hash.Bytes(), nil
}

//SnapshotHandler does nothing at the moment
func (h *Handler) SnapshotHandler(blockIndex int) ([]byte, error) {
	return []byte{}, nil
}

//RestoreHandler does nothing at the moment
func (h *Handler) RestoreHandler(snapshot []byte) ([]byte, error) {
	return []byte{}, nil
}
