package Nats

import (
	"L00/Cache"
	"L00/DB"
	"L00/DB/Model"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/nats-io/stan.go"
)

// Подключение к Nats
func СonnectToNats() (stan.Conn, error) {
	clusterID := "My-Cluster-ID"
	clientID := "My-Client-ID"
	natsURL := "nats://localhost:4222"

	natsConn, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Ошибка при подключении к серверу NATS: %v", err)
		return nil, err
	}

	log.Printf("Успешно подключились к серверу NATS")

	return natsConn, nil
}

// Подписка на канал в NATS
func SubscribeToNats(natsConn stan.Conn, subject string, dbOpen *sql.DB, cache map[string]Model.ModelOrder) error {
	_, err := natsConn.Subscribe(subject, func(msg *stan.Msg) {
		MessageHandler(msg, dbOpen, cache)
	})
	if err != nil {
		log.Fatalf("Ошибка при подписке на канал в NATS: %v", err)
		return err
	}

	log.Printf("Успешно подписались на канал в NATS")

	return nil
}

// Обработчик сообщения из очереди и если с сообщением все в порядке, сохраняем в бд и кеш
func MessageHandler(msg *stan.Msg, dbOpen *sql.DB, cache map[string]Model.ModelOrder) {
	jsonData := msg.Data
	var order Model.ModelOrder
	err := json.Unmarshal(jsonData, &order)
	if err != nil {
		log.Printf("Данные не десериализованы (не правильный формат)")
		return
	}
	log.Printf("Данные успешно десериализованы.")

	if _, exists := cache[order.Order_uid]; exists {
		log.Printf("Заказ с uid %s уже существует в базе данных.", order.Order_uid)
		return
	}

	err = DB.SaveOrderToDB(dbOpen, order)
	if err != nil {
		log.Fatalf("Ошибка сохранения в БД после десериализации: %v", err)
	}
	Cache.AddToCache(order, cache)
}
