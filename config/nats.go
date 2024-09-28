package config

import (
	"fmt"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

type NatsBroker interface {
	CoinReceive()
}

type NatsBrokerImpl struct {
	nc *nats.Conn
}

func (u NatsBrokerImpl) CoinReceive() {
	_, err := u.nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Checking: %s\n", string(m.Data))
	})
	if err != nil {
		return
	}
}

func (u NatsBrokerImpl) Init() {
	log.Infoln("Initialization nats...")
	u.CoinReceive()
	log.Infoln("Initialized nats success")
}

func NatsBrokerInit(nc *nats.Conn) *NatsBrokerImpl {
	broker := &NatsBrokerImpl{
		nc: nc,
	}
	broker.Init()
	return broker
}
