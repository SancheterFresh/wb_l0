package database

import (
	"context"
	"fmt"
	"w0/data"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func PoolConnect() *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), "postgresql://postgres:password@localhost:5432/wbl0")

	if err != nil {
		logrus.Error(err)
	}
	return pool

}

func InsertOrder(order *data.Order, pool *pgxpool.Pool) error {
	orders_query := fmt.Sprintf("INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v', '%v')", order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	_, err := pool.Exec(context.Background(), orders_query)
	if err != nil {
		logrus.Error(err)
		return err
	}

	delivery_query := fmt.Sprintf("INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email) VALUES ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v')", order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	_, err = pool.Exec(context.Background(), delivery_query)
	if err != nil {
		logrus.Error(err)
		return err
	}

	payment_query := fmt.Sprintf("INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ('%v', '%v', '%v', '%v', '%v', %v, %v, '%v', %v, %v, %v)", order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	_, err = pool.Exec(context.Background(), payment_query)
	if err != nil {
		logrus.Error(err)
		return err
	}

	for _, v := range order.Items {
		item_query := fmt.Sprintf("INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ('%v', %v, '%v', %v , '%v', '%v', %v, '%v', %v, %v, '%v', %v)", order.OrderUID, v.ChrtID, v.TrackNumber, v.Price, v.Rid, v.Name, v.Sale, v.Size, v.TotalPrice, v.NmID, v.Brand, v.Status)
		_, err = pool.Exec(context.Background(), item_query)
		if err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}

func RecoverCash(pool *pgxpool.Pool) []data.Order {
	query := "SELECT * FROM orders"

	ordrs, err := pool.Query(context.Background(), query)
	if err != nil {
		logrus.Error(err)
	}
	var c_orders []data.Order
	for ordrs.Next() {
		var c_order data.Order

		err = ordrs.Scan(&c_order.OrderUID, &c_order.TrackNumber, &c_order.Entry, &c_order.Locale, &c_order.InternalSignature, &c_order.CustomerID, &c_order.DeliveryService, &c_order.Shardkey, &c_order.SmID, &c_order.DateCreated, &c_order.OofShard)

		if err != nil {
			logrus.Error(err)
		}

		delivery_query := fmt.Sprintf("SELECT * FROM deliveries WHERE order_uid='%v'", c_order.OrderUID)
		dlvrs, err := pool.Query(context.Background(), delivery_query)
		if err != nil {
			logrus.Error(err)
		}
		var c_delivery data.Delivery
		var c_uuid string
		dlvrs.Scan(c_uuid, &c_delivery.Name, &c_delivery.Phone, &c_delivery.Zip, &c_delivery.City, &c_delivery.Address, &c_delivery.Region, &c_delivery.Email)
		c_order.Delivery = c_delivery

		payment_query := fmt.Sprintf("SELECT * FROM payments WHERE order_uid='%v'", c_order.OrderUID)
		pmnts, err := pool.Query(context.Background(), payment_query)
		if err != nil {
			logrus.Error(err)
		}
		var c_payment data.Payment

		pmnts.Scan(c_uuid, &c_payment.Transaction, &c_payment.RequestID, &c_payment.Currency, &c_payment.Provider, &c_payment.Amount, &c_payment.PaymentDt, &c_payment.Bank, &c_payment.DeliveryCost, &c_payment.GoodsTotal, &c_payment.CustomFee)
		c_order.Payment = c_payment

		items_query := fmt.Sprintf("SELECT * FROM items WHERE order_uid='%v'", c_order.OrderUID)
		itms, err := pool.Query(context.Background(), items_query)
		if err != nil {
			logrus.Error(err)
		}
		var c_items []data.Item
		for itms.Next() {
			var c_item data.Item
			pmnts.Scan(c_uuid, &c_item.ChrtID, &c_item.TrackNumber, &c_item.Price, &c_item.Rid, &c_item.Name, &c_item.Sale, &c_item.Size, &c_item.TotalPrice, &c_item.NmID, &c_item.Brand, &c_item.Status)
			c_items = append(c_items, c_item)
		}
		c_order.Items = c_items
		c_orders = append(c_orders, c_order)

	}
	return c_orders

}
