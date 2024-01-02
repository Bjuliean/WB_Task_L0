package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"wbl0/WB_Task_L0/internal/config"
	"wbl0/WB_Task_L0/internal/logs"
	"wbl0/WB_Task_L0/internal/models"
	storagemanager "wbl0/WB_Task_L0/internal/storage_manager"

	"github.com/nats-io/stan.go"
)

type Broker struct {
	cn      stan.Conn
	storage *storagemanager.StorageManager
	logs    *logs.Logger
	cfg     *config.Config
}

func New(cfg *config.Config, storage *storagemanager.StorageManager, logsHandler *logs.Logger) *Broker {
	const ferr = "internal.broker.New"

	nwcn, err := stan.Connect(cfg.NatsStreaming.ClusterID, cfg.NatsStreaming.ClientID,
		stan.NatsURL(fmt.Sprintf("%s:%s", cfg.NatsStreaming.Host, cfg.NatsStreaming.Port)))
	if err != nil {
		log.Fatalf("%s: failed to connect nats-streaming: %s", ferr, err)
		return nil
	}

	return &Broker{
		cn:      nwcn,
		storage: storage,
		logs:    logsHandler,
		cfg:     cfg,
	}
}

func (b *Broker) CloseConnection() { // todo: close error
	b.cn.Close()
}

func (b *Broker) SubscribeAndHandle() error {
	const ferr = "internal.broker.SubscribeAndHandle"

	_, err := b.cn.Subscribe(b.cfg.NatsStreaming.SubscribeSubject,
		brokerMsgHandler(b.logs, b.storage), stan.SetManualAckMode())
	if err != nil {
		b.logs.WriteError(ferr, err.Error())
		return err
	}

	return nil
}

func brokerMsgHandler(logs *logs.Logger, storage *storagemanager.StorageManager) stan.MsgHandler {
	const ferr = "internal.broker.brokerMsgHandler"

	return func(msg *stan.Msg) {

		var order models.Order

		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			msg.Ack()
			logs.WriteError(ferr, err.Error())
			return
		}

		err = storage.AddOrder(order)
		if err != nil {
			msg.Ack()
			logs.WriteError(ferr, err.Error())
			return
		}

		if err = msg.Ack(); err != nil {
			logs.WriteError(ferr, err.Error())
			return
		}

		logs.WriteInfo(fmt.Sprintf("successfully saved order from broker: %v", order.OrderUID))
	}
}
