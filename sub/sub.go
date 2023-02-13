package sub

import (
	"encoding/json"
	"fmt"
	"sync"

	//"time"
	"w0/data"
	"w0/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func StertSub(pool *pgxpool.Pool) {
	url := "nats://localhost:4222"
	sub, err := nats.Connect(url)
	if err != nil {
		fmt.Printf("Subscriber: ")
		logrus.Fatal(err)

	}
	defer sub.Close()

	sub.Subscribe("orders", func(msg *nats.Msg) {
		var order data.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			logrus.Error(err)
			return
		}

		err = database.InsertOrder(&order, pool)
		if err != nil {
			logrus.Error(err)
			return
		}
	})

	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()

}
