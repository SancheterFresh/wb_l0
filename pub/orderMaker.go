package pub

import (
	"time"
	"w0/data"

	"github.com/google/uuid"
)

func getNewOrder() *data.Order {
	id := uuid.New()

	order := data.Order{
		OrderUID:    id.String(),
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: data.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: data.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "reqid",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []data.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        21454,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "kz",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now().String(),
		OofShard:          "1",
	}

	return &order
}
