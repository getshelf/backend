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

### Транзитивные зависимости

Эти модули подтягиваются основными библиотеками автоматически и напрямую в коде не используются:

- `go.uber.org/dig`, `go.uber.org/zap`, `go.uber.org/multierr`, `go.uber.org/atomic` — внутренние зависимости Uber Fx.
- `github.com/cespare/xxhash/v2` — зависимость Redis клиента.
- `go.yaml.in/yaml/v3` — зависимость Testify для YAML-related helpers.
- `golang.org/x/sys` — низкоуровневая зависимость Go crypto/tooling packages.

Актуальный список версий находится в [go.mod](go.mod), а контрольные суммы — в [go.sum](go.sum). Обновить зависимости можно командами:

```bash
go get <module>@<version>
go mod tidy
```

## 1. Подготовить `.env`

Скопируйте шаблон конфигурации:

```bash
cp .env.example .env
```

Не добавляйте `.env` в git: файл содержит локальные credentials.

Для Docker Compose используются такие параметры:

```env
DB_HOST=db
DB_PORT=5432
DB_USERNAME=getshelf_user
DB_PASSWORD=getshelf_password
DB_NAME=getshelf_db
DB_SSLMODE=disable

HTTP_HOST=0.0.0.0
HTTP_PORT=8080
SWAGGER_PORT=8081

REDIS_HOST=redis
REDIS_PORT=6379

JWT_SECRET=change-me-in-production
JWT_ACCESS_TTL_MINUTES=15
JWT_REFRESH_TTL_HOURS=24
```

## 2. Запуск через Docker

Запустите PostgreSQL, Redis, API и Swagger UI:

```bash
make docker-up
```

Адреса после запуска:

- API: http://localhost:8080
- Swagger UI: http://localhost:8081
- OpenAPI файл: [api/openapi.yaml](api/openapi.yaml)

Миграции применяются автоматически при старте API через `go-migrate`.

Остановить контейнеры:

```bash
make docker-down
```

Если PostgreSQL был создан со старыми credentials, пересоздайте dev volume:

```bash
make docker-reset
```

Внимание: `docker-reset` удаляет volume `postgres_data` и локальные данные PostgreSQL.

## 3. Создание миграций

Новая миграция создаётся командой:

```bash
make migration-create name=create_books
```

Команда создаст файлы в `internal/infrastructure/persistence/postgres/migrations`:

```text
000002_create_books.up.sql
000002_create_books.down.sql
```

Добавьте SQL-команды в `up`-файл и обратные команды в `down`-файл. После запуска API новая миграция применится автоматически.

Для каждой новой миграции используйте новое имя и номер:

```bash
make migration-create name=add_user_avatar
```

## 4. Линтеры, чекеры и тесты

Форматирование Go-файлов:

```bash
gofmt -w $(find cmd internal test -name '*.go')
```

Проверка, какие файлы не отформатированы:

```bash
test -z "$(gofmt -l $(find cmd internal test -name '*.go'))"
```

Статический анализ:

```bash
make vet
```

Все unit и integration тесты:

```bash
make test
```

Только unit-тесты:

```bash
make test-unit
```

Только integration-тесты:

```bash
make test-integration
```

Проверка с race detector:

```bash
go test -race ./...
```

Сборка приложения:

```bash
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
