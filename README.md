# GetShelf Backend

Backend API на Go: domain, application, infrastructure и REST presentation layers.

## Требования

- Go 1.26+
- Docker и Docker Compose

## Сторонние библиотеки

Основные Go-зависимости проекта:

| Библиотека | Назначение |
| --- | --- |
| `github.com/go-chi/chi/v5` | HTTP router и вложенные route groups для REST API. |
| `github.com/golang-jwt/jwt/v5` | Создание и проверка JWT access/refresh tokens. |
| `github.com/golang-migrate/migrate/v4` | Создание и автоматическое применение PostgreSQL migrations. Используются source `iofs` и database driver `postgres`. |
| `github.com/jmoiron/sqlx` | Работа с PostgreSQL поверх `database/sql`: подключения, транзакции и расширенные SQL helpers. |
| `github.com/joho/godotenv` | Загрузка локальных переменных из `.env`. |
| `github.com/lib/pq` | PostgreSQL driver для `database/sql` и `sqlx`. |
| `github.com/redis/go-redis/v9` | Подключение к Redis и хранение отозванных JWT с TTL. |
| `golang.org/x/crypto` | bcrypt hashing и проверка паролей. |
| `go.uber.org/fx` | Dependency injection и lifecycle приложения. |
| `github.com/stretchr/testify` | `assert` и `require` в unit/integration тестах. |
| `github.com/DATA-DOG/go-sqlmock` | Mock SQL driver для integration-тестов без запущенного PostgreSQL. |

```bash
go mod tidy
```

## Запуск `.env`

```bash
# Copy .env file
cp .env.example .env
# Start backend
make docker-up
```

## Адресса:

- API: http://localhost:8080
- Swagger UI: http://localhost:8081
- OpenAPI файл: [api/openapi.yaml](api/openapi.yaml)

## Команды

```bash
make docker-down
```

```bash
make docker-reset
```

```bash
make migration-create name=create_books
```

```bash
make migration-create name=add_user_avatar
```

## 4. Линтеры, чекеры и тесты

```bash
# Форматирование Go-файлов:
gofmt -w $(find cmd internal test -name '*.go')

# Проверка, какие файлы не отформатированы:
test -z "$(gofmt -l $(find cmd internal test -name '*.go'))"

# Статический анализ:
make vet

# Все unit и integration тесты:
make test

# Только unit-тесты:
make test-unit

# Только integration-тесты:
make test-integration

# Проверка с race detector:
go test -race ./...

# Сборка приложения:
make build
```

Полный чек перед коммитом:

```bash
test -z "$(gofmt -l $(find cmd internal test -name '*.go'))" && make vet && make test && make build
```

Integration-тесты используют `sqlmock`, поэтому PostgreSQL и Redis для них не нужны.

Подробные схемы запросов и ответов доступны в Swagger UI.

## Структура проекта

- `cmd/api` — запуск приложения и Uber Fx DI
- `api` — OpenAPI спецификация
- `internal/config` — конфигурация из `.env`
- `internal/application` — application services
- `internal/domain/user` — entity, value objects, ports и domain services
- `internal/infrastructure` — PostgreSQL, sqlx, Redis, JWT, bcrypt и миграции
- `internal/presentation/rest` — schemas, chi routers и HTTP server
- `test/unit` — unit-тесты
- `test/integration` — integration-тесты
