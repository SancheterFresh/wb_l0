package pub

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func StartPub() {
	url := "nats://localhost:4222"
	nc, err := nats.Connect(url)
	if err != nil {
		fmt.Printf("Publisher: ")
		logrus.Fatal(err)
	}
	defer nc.Close()

	for i := 0; i < 1000; i++ {
		order := getNewOrder()
		order_json, err := json.Marshal(order)
		if err != nil {
			logrus.Error(err)
			return
		}
		err = nc.Publish("orders", order_json)
		if err != nil {
			logrus.Error(err)
			return
		}

		time.Sleep(5 * time.Second)
	}

}
