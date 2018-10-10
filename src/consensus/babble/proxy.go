package babble

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mosaicnetworks/babble/src/hashgraph"
	"github.com/mosaicnetworks/babble/src/proxy/inmem"
	"github.com/mosaicnetworks/evm-lite/src/service"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"
)

//Proxy implements the Babble AppProxy interface by inheriting Babble's
//InmemProxy methods
type Proxy struct {
	*inmem.InmemProxy

	service *service.Service
	state   *state.State
	logger  *logrus.Logger
}

//NewProxy initializes and returns a new Proxy
func NewProxy(state *state.State,
	service *service.Service,
	logger *logrus.Logger) *Proxy {

	proxy := &Proxy{
		service: service,
		state:   state,
		logger:  logger,
	}

	proxy.InmemProxy = inmem.NewInmemProxy(
		proxy.getCommitHandler(),
		proxy.getSnapshotHandler(),
		proxy.getRestoreHandler(),
		logger)

	return proxy
}

func (p *Proxy) getCommitHandler() inmem.CommitHandler {
	var commitHandler inmem.CommitHandler = func(block hashgraph.Block) ([]byte, error) {
		p.logger.Debug("CommitBlock")

		blockHashBytes, err := block.Hash()
		blockHash := common.BytesToHash(blockHashBytes)

		for i, tx := range block.Transactions() {
			if err := p.state.ApplyTransaction(tx, i, blockHash); err != nil {
				return []byte{}, err
			}
		}

		hash, err := p.state.Commit()
		if err != nil {
			return []byte{}, err
		}

		return hash.Bytes(), nil
	}
	return commitHandler
}

func (p *Proxy) getSnapshotHandler() inmem.SnapshotHandler {
	var snapshotHandler inmem.SnapshotHandler = func(blockIndex int) ([]byte, error) {
		return []byte{}, nil
	}
	return snapshotHandler
}

func (p *Proxy) getRestoreHandler() inmem.RestoreHandler {
	var restoreHandler inmem.RestoreHandler = func(snapshot []byte) ([]byte, error) {
		return []byte{}, nil
	}
	return restoreHandler
}
