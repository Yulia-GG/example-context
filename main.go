package main

import(
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/data", handler) // при обращении к /data, вызываем обработчик handler
	fmt.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
// обработчик handler
func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Запрос получен")

	// создаем канал для результата/операции
	resultChan := make(chan string, 1)

	// запускаем гортину, эмулируем долгую операцию
	go func() {
		log.Println("Операция началась")

		// допустим, операция длится 5 секунд
		time.Sleep(5 *time.Second)
		resultChan <- "готовые данные"
	}()

	// select ждем или результат, или отмену ctx
	select {
	case res := <- resultChan:
		log.Println("Операция завершеноа успешно")
		w.Write([]byte(res))
	case <- ctx.Done():
		// клиент прервал соединение
		log.Println("Контекст отменен:", ctx.Err())
		http.Error(w, "Запрос отменен клиентом", http.StatusRequestTimeout)
	}
}

