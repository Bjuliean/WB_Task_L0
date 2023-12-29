package broker

import (
	"fmt"
	"log"
	"wbl0/WB_Task_L0/internal/config"
	"wbl0/WB_Task_L0/internal/logs"

	"github.com/nats-io/stan.go"
)

type Broker struct {
	cn stan.Conn
	logs *logs.Logger
}

func New(cfg *config.Config, logsHandler *logs.Logger) *Broker {
	const ferr = "internal.broker.New"

	nwcn, err := stan.Connect(cfg.NatsStreaming.ClusterID, cfg.NatsStreaming.ClientID,
	stan.NatsURL(fmt.Sprintf("%s:%s", cfg.NatsStreaming.Host, cfg.NatsStreaming.Port)))
	if err != nil {
		log.Fatalf("%s: failed to connect nats-streaming: %s", ferr, err)
		return nil
	}

	return &Broker{
		cn: nwcn,
		logs: logsHandler,
	}
}

func (b *Broker) CloseConnection() { // todo: close error
	b.cn.Close()
}