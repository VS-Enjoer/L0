package DB

import (
	"L00/DB/Model"
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Подключаемся к бд
func DbConnect() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=postgres password=1202 dbname=L00 host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatal("Ошибка при подключении к БД: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка при пинге БД: %v", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Ошибка при создании драйвера: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./DB/migration",
		"postgres", driver)
	if err != nil {
		log.Fatal("Ошибка при создании миграции: %v", err)

	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal("Ошибка при выполнении миграции: %v", err)
		}
	}

	return db, nil
}

// Берем данные из табличек бд и возвращаем заказ
func LoadOrdersFromDB(db *sql.DB) ([]Model.ModelOrder, error) {
	query := `
        SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, 
               o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, 
               o.oof_shard, 
               d.order_uid, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email, 
               p.order_uid, p.transaction, p.request_id, p.currency, p.provider, p.amount, 
               p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee, 
               i.order_uid, i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, 
               i.size, i.total_price, i.nm_id, i.brand, i.status
        FROM orders o
        LEFT JOIN order_delivery d ON o.order_uid = d.order_uid
        LEFT JOIN order_payment p ON o.order_uid = p.order_uid
        LEFT JOIN order_items i ON o.order_uid = i.order_uid
    `

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Ошибка при выполнении SQL-запроса: %v", err)
		return nil, err
	}
	defer rows.Close()

	var orders []Model.ModelOrder

	for rows.Next() {
		var order Model.ModelOrder
		var delivery Model.Delivery
		var payment Model.Payment
		var items Model.Items

		err := rows.Scan(
			&order.Order_uid,
			&order.Track_number,
			&order.Entry,
			&order.Locale,
			&order.Internal_signature,
			&order.Customer_id,
			&order.Delivery_service,
			&order.Shardkey,
			&order.Sm_id,
			&order.Date_created,
			&order.Oof_shard,
			&delivery.Order_uid,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
			&payment.Order_uid,
			&payment.Transaction,
			&payment.Request_id,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.Payment_dt,
			&payment.Bank,
			&payment.Delivery_cost,
			&payment.Goods_total,
			&payment.Custom_fee,
			&items.Order_uid,
			&items.Chrt_id,
			&items.Track_number,
			&items.Price,
			&items.Rid,
			&items.Name,
			&items.Sale,
			&items.Size,
			&items.Total_price,
			&items.Nm_id,
			&items.Brand,
			&items.Status,
		)
		if err != nil {
			log.Fatalf("Ошибка при сканировании результата: %v", err)
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Ошибка при обработке результатов: %v", err)
		return nil, err
	}

	return orders, nil
}
