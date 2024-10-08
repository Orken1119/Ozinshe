# Сервис для просмотра мультсериалов
Описание: Проект представляет собой бэкэнд для сервиса для просмотра сериалов и фильмов. 

## Основные возможности:
• Добавление и удаление любимых мультсериалов.  

• Получение топ-50 самых просматриваемых фильмов.  

• Получение недавно просмотренных эпизодов.  

• Получение похожих фильмов на основе жанров фильма.  

### Стек используемых технологий:
• Go (Golang) — основная технология для написания серверной части.  

• [Gin](https://gin-gonic.com/docs/) — Web framework для организации маршрутов и обработки HTTP-запросов.  

• [swageer](https://github.com/swaggo/swag)Swag — интеграция со Swagger для автоматической генерации документации API.  

• [JWT-go](https://github.com/dgrijalva/jwt-go) — реализация токенов JSON Web Tokens (JWT) для аутентификации.  

PostgreSQL — база данных для хранения информации о фильмах, пользователях и оценках.  

### Установка и запуск проекта:
1. Установить Go
Если Go еще не установлен, его можно скачать с официального сайта: https://golang.org/dl/.
После установки убедитесь, что Go настроен правильно:
go version
2. Установить зависимости
В проекте все зависимости описаны в файлах go.mod и go.sum. Для их установки выполните:
go mod tidy
3. Запустить проект
Выполните команду для запуска приложения:
go run ./cmd
Приложение будет запущено на вашем локальном сервере.

### Доступ к Swagger UI:
Swagger-документация доступна по следующему адресу:
http://localhost:1136/swagger/index.html
