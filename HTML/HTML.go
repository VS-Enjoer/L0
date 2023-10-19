package HTML

import (
	"L00/DB/Model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Обработчик для поиска заказа по ID в кеше
func searchOrderHandler(cache map[string]Model.ModelOrder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			orderID := r.FormValue("order_id")
			result, ok := cache[orderID]

			if ok {
				response := make(map[string]interface{})
				response["order_id"] = orderID
				response["result"] = result
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				log.Printf("Запрос на поиск заказа по ID: %s. Результат успешно отправлен.", orderID)
			} else {
				response := "Заказ не найден"
				w.Header().Set("Content-Type", "text/html")
				fmt.Fprint(w, response)
				log.Printf("Запрос на поиск заказа по ID: %s. Заказ не найден.", orderID)
			}
		}
	}
}

// Запуск HTTP-сервера и обработка запросов
func HtmlStart(cache map[string]Model.ModelOrder) {
	http.HandleFunc("/search_order", searchOrderHandler(cache))

	http.Handle("/", http.FileServer(http.Dir("HTML/View")))

	port := "8080"
	fmt.Printf("Сервер запущен. http://localhost:8080")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске HTTP-сервера: %v", err)
	}
}
