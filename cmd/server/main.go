package main

import (
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/handler"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/server"
)

// Пример запроса к серверу:
// POST /update/counter/someMetric/527 HTTP/1.1
// Host: localhost:8080
// Content-Length: 0
// Content-Type: text/plain

// Пример ответа от сервера:
// HTTP/1.1 200 OK
// Date: Tue, 21 Feb 2023 02:51:35 GMT
// Content-Length: 11
// Content-Type: text/plain; charset=utf-8

// func main() {
// 	fmt.Println("Hello from othercmd application!")
// 	package2.FunctionFromPackage2()
// }

// ADDRESS отвечает за адрес эндпоинта HTTP-сервера.

func main() {
	handler := handler.NewHandler()
	handler.InitStorage()

	server.InitRouter(*handler)
}
