package DB

import (
	"L00/DB/Model"
	"database/sql"
	"log"
)

// Сохранение заказа в бд (вынес его в отдельный файл т.к. получился большой)
func SaveOrderToDB(db *sql.DB, order Model.ModelOrder) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Ошибка при начале транзакции: %v", err)
		return err
	}

	_, err = tx.Exec("INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.Order_uid, order.Track_number, order.Entry, order.Locale, order.Internal_signature, order.Customer_id, order.Delivery_service, order.Shardkey, order.Sm_id, order.Date_created, order.Oof_shard)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Ошибка при выполнении INSERT INTO orders: %v", err)
		return err
	}

	_, err = tx.Exec("INSERT INTO order_delivery (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		order.Order_uid, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Ошибка при выполнении INSERT INTO order_delivery: %v", err)
		return err
	}

	_, err = tx.Exec("INSERT INTO order_payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.Order_uid, order.Payment.Transaction, order.Payment.Request_id, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.Payment_dt, order.Payment.Bank, order.Payment.Delivery_cost, order.Payment.Goods_total, order.Payment.Custom_fee)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Ошибка при выполнении INSERT INTO order_payment: %v", err)
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec("INSERT INTO order_items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			order.Order_uid, item.Chrt_id, item.Track_number, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.Total_price, item.Nm_id, item.Brand, item.Status)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Ошибка при выполнении INSERT INTO order_items: %v", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Ошибка при коммите транзакции: %v", err)
		return err
	}

	return nil
}
