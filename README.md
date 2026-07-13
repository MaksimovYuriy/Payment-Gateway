# Payment Gateway

Payment Gateway - это Go-сервис для регистрации банков и мерчантов, создания платежей и проведения платежных попыток через абстракцию банковского процессора.

Проект построен в стиле небольшой clean-architecture структуры:

- `controller/restapi` обрабатывает HTTP-запросы, валидацию и JSON-ответы.
- `usecase` содержит сценарии приложения и бизнес-проверки.
- `repo` содержит контракты хранилища и PostgreSQL-реализации репозиториев.
- `entity` содержит доменные сущности и базовую доменную валидацию.
- `storage/postgres` содержит подключение к PostgreSQL и реализацию Unit of Work.

## Возможности

- Регистрация и получение списка банков.
- Регистрация и получение списка мерчантов.
- Генерация API-ключа для мерчанта. Сырой ключ возвращается один раз, в базе хранится только его хэш.
- Создание платежа и получение статуса платежа.
- Проведение платежа через mock-реализацию банковского процессора.
- Unit of Work для сценария создания платежа.
- Валидация входящих запросов и единый формат JSON-ошибок.
- Unit-тесты и тесты HTTP-хэндлеров.

## Быстрый старт

### Требования

- Go
- Docker и Docker Compose
- Goose CLI для миграций

Если Goose не установлен:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 1. Запустить PostgreSQL

```bash
docker compose up -d
```

PostgreSQL будет доступен по адресу:

```text
localhost:5434
```

Данные для подключения описаны в `docker-compose.yml` и `config/config.json`:

```text
database: payment_gateway
user: payment_gateway
password: payment_gateway
```

### 2. Применить миграции

```bash
make migrate-up
```

Проверить статус миграций:

```bash
make migrate-status
```

### 3. Запустить API

```bash
go run ./cmd/api
```

Или запустить PostgreSQL и API одной командой:

```bash
make build-up
```

API будет доступно по адресу:

```text
http://localhost:8081
```

### 4. Проверить healthcheck

```bash
curl http://localhost:8081/health
```

Ожидаемый ответ:

```json
{
  "service": "payment-gateway",
  "status": "ok"
}
```

## API

### Банки

```text
POST /admin/banks
GET  /admin/banks
```

### Мерчанты

```text
POST /admin/merchants
GET  /admin/merchants
```

При регистрации мерчанта API-ключ возвращается один раз. Его нужно сразу сохранить.

### Платежи

```text
POST /payments
GET  /payments
GET  /payments/{id}
```

Сейчас используется mock-реализация банковского процессора, поэтому платежи проводятся локально без обращения во внешний банк.

## Тесты

Запуск всех тестов:

```bash
go test ./...
```

Если в текущем окружении Go build cache недоступен для записи, используйте:

```bash
GOCACHE=/tmp/go-build go test ./...
```
