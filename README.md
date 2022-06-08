# Тестовое задание

* Разворачивать Postgresql и Nats-streaming в Docker (имеется файл docker-compose)

* Запуск сервера: `go run ./cmd/main.go`
* Отправка публикаций в Nats с помощью: `go run ./pkg/publisher.go -p {path_to_files}`

По умолчанию http-serv работает на 8000 порту на localhost'e (config.yml)
