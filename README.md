## Describe
Сервис предоставляющий API по созданию сокращённых ссылок, написанный на языке Golang. У него есть 2 эндпоинта:
1. Метод Post, который сохраняет оригинальный URL в базе данных и возвращает сокращённый.
2. Метод Get, который принимает сокращённый URL и возвращает оригинальный.

Сервер реализован с помощью технологий gRPC и gRPC-gateway, что позволяет обращаться к сервису и по HTTP. Для проверки работы функционала можно зайти в swagger, который расположен по адресу http://host:port/swagger

Так же имеются Unit и интеграционные тесты.


## Configuration
- Все параметры сервера задаются в файле `.env`
- Параметры для Postgres задаются в файле `docker-compose.service.yml` / `docker-compose.db.yml` / `docker-compose.test.yml`
- Для запуска тестов, обращающихся к базе данных Postgres необходимо запустить тестовую бд ``docker-compose.test.yml``
- Для запуска интеграционных тестов, необходимо запустить сервис и сконфегурировать `SHORTENER_SERVER_HOST` в файле `./internal/tests/integration/url_shortener_test.go`
- При запуске сервиса локально, необходимо изменить в `.env` поле `POSTGRES_HOST` с `url_shortener_postgres` на `localhost`


## Command
- Узнать все команды `make help`
- Docker
  - 
  - Создать docker образ сервиса с параметром database=in-memory/postgres `make df database=postgres` (по умолчанию значение database=in-memory)
  - Запуск сервиса и базы данных в Docker `make service_up`
  - Запуск только базы данных Postgres в Docker `make docker_db` (может понадобиться при запуске сервиса локально)
  - Запуск тестовой базы данных Postgres в Docker `make docker_test`
- Local
  - 
  - Генерация proto файлов и swagger `make proto`
  - Создание исполняющего файла сервиса `make build`
  - Запуск сервиса локально `make run`
- Test
  - 
  - Для запуска Unit тестов `make unit-test`
  - Для запуска интеграционных тестов `make integration-test`



