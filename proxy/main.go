package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	r := chi.NewRouter()
	hugoProxy := NewReverseProxy("hugo", "1313")
	// ...

	r.Use(hugoProxy.ReverseProxy)

	r.Get("/api/*", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from API")
	})

	go WorkerTest()

	log.Println("Proxy server started on :8080")
	http.ListenAndServe(":8080", r)
}

const contentTemplate = `Текущее время: %s
Счетчик: %d
`

func WorkerTest() {
	// Интервал обновления данных
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	counter := 0 // Начальное значение счётчика

	for {
		select {
		case <-ticker.C:
			// Формируем текущую дату и время
			currentTime := time.Now().Format("2006-01-02 15:04:05")

			// Формируем содержимое файла
			content := fmt.Sprintf(contentTemplate, currentTime, counter)

			// Проверяем наличие каталога и создаём его, если он отсутствует
			filePath := "/app/static/tasks/_index.md"
			dirPath := "/app/static/tasks"
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				err := os.MkdirAll(dirPath, 0755)
				if err != nil {
					log.Fatalf("Не удалось создать каталог: %v", err)
				}
			}

			// Записываем содержимое в файл
			err := os.WriteFile(filePath, []byte(content), 0644)
			if err != nil {
				log.Println("Ошибка записи файла:", err)
			} else {
				log.Printf("Файл обновлён: время %s, счётчик %d", currentTime, counter)
			}

			// Увеличиваем счётчик
			counter++
		}
	}
}
