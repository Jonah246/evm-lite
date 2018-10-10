package babble

import (
	_babble "github.com/mosaicnetworks/babble/src/babble"
	_proxy "github.com/mosaicnetworks/babble/src/proxy"
	_inmem "github.com/mosaicnetworks/babble/src/proxy/inmem"
	"github.com/mosaicnetworks/evm-lite/src/config"
	"github.com/mosaicnetworks/evm-lite/src/service"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"
)

//InmemBabble implementes the Consensus interface.
//It uses an inmemory Babble node.
type InmemBabble struct {
	ethState   *state.State
	ethService *service.Service
	proxy      _proxy.AppProxy

	config *config.BabbleConfig
	babble *_babble.Babble

	logger *logrus.Logger
}

//NewInmemBabble instantiates a new InmemBabble consensus system
func NewInmemBabble(config *config.BabbleConfig, logger *logrus.Logger) *InmemBabble {
	return &InmemBabble{
		config: config,
		logger: logger,
	}
}

/*******************************************************************************
IMPLEMENT CONSENSUS INTERFACE
*******************************************************************************/

//Init instantiates a Babble inmemory node
func (b *InmemBabble) Init(state *state.State, service *service.Service) error {

	b.logger.Debug("INIT")

	b.ethState = state
	b.ethService = service
	b.proxy = _inmem.NewInmemProxy(NewHandler(state, service), b.logger)

	realConfig := b.config.ToRealBabbleConfig(b.logger)
	realConfig.Proxy = b.proxy

	babble := _babble.NewBabble(realConfig)

	err := babble.Init()
	if err != nil {
		return err
	}

	b.babble = babble

	return nil
}

//Run starts the Babble node and relay txs from SubmitCh to Babble
func (b *InmemBabble) Run() error {

	go b.babble.Run()

	serviceSubmitCh := b.ethService.GetSubmitCh()
	proxySubmitCh := b.proxy.SubmitCh()

	for {
		select {
		case t := <-serviceSubmitCh:
			proxySubmitCh <- t
		}
	}
}

//Info returns Babble stats
func (b *InmemBabble) Info() (map[string]string, error) {
	info := b.babble.Node.GetStats()
	info["type"] = "babble"
	return info, nil
}
