package Cache

import (
	"L00/DB"
	"L00/DB/Model"
	"database/sql"
	"log"
	"sync"
)

var mu sync.RWMutex

// Загрузка данных из бд в кеш
func LoadCacheFromDB(dbOpen *sql.DB, cache map[string]Model.ModelOrder) error {
	orders, err := DB.LoadOrdersFromDB(dbOpen)
	if err != nil {
		log.Fatalf("Ошибка при загрузке данных из БД в кеш: %v", err)
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	for _, order := range orders {
		cache[order.Order_uid] = order
	}

	log.Printf("Данные успешно загружены из БД в кеш")

	return nil
}

// Добавление данных в кеш
func AddToCache(order Model.ModelOrder, cache map[string]Model.ModelOrder) {
	mu.Lock()
	defer mu.Unlock()
	cache[order.Order_uid] = order
	log.Printf("Заказ успешно добавлен в кеш: %s", order.Order_uid)
}
