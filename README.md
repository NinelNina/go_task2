# Ledger & Gateway Service

Учебный проект на языке Go, состоящий из двух модулей в одном репозитории:

* **Gateway** — HTTP-сервер (REST API) для приёма клиентских запросов. Импортирует модуль `ledger` и вызывает его функции.
* **Ledger** — библиотека бизнес-логики: учёт транзакций, бюджеты по категориям и валидация. Данные хранятся в памяти.

## Структура проекта

```text
go_task2/
├── go.work                    # объединяет модули gateway и ledger
├── gateway/                   # HTTP-слой
│   ├── go.mod
│   ├── main.go                # запуск сервера на :8080
│   └── internal/api/
│       ├── dto.go             # DTO запросов/ответов и функции-преобразователи
│       ├── handlers.go        # обработчики эндпоинтов, writeJSON / writeError
│       ├── handlers_test.go   # интеграционные HTTP-тесты (httptest)
│       ├── middleware.go      # middleware логирования запросов
│       └── router.go          # NewRouter: маршруты /api/... + middleware
├── ledger/                    # бизнес-логика (библиотека)
│   ├── go.mod
│   ├── budget.go              # Budget, SetBudget, ListBudgets, LoadBudgets
│   ├── transaction.go         # Transaction, AddTransaction, ListTransactions, ErrBudgetExceeded
│   ├── validatable.go         # интерфейс Validatable и функция CheckValid
│   ├── reset.go               # Reset: очистка данных в памяти (для тестов)
│   ├── ledger_test.go         # unit-тесты домена
│   └── budgets.json           # пример файла с лимитами категорий
├── .gitignore
└── README.md
```

## Требования

* **Go** версии 1.22 или выше (используются шаблоны маршрутов с методами в `http.ServeMux`).
* Утилита `curl` для проверки эндпоинтов.

## Запуск

Модули связаны через `go.work`, поэтому дополнительная настройка `go.mod` не нужна.

```bash
cd gateway
go run .
```

Сервер запускается на порту `:8080`:

```text
2026/10/29 00:13:42 Gateway listening on :8080
```

## API

Все эндпоинты имеют префикс `/api`. Ответы возвращаются в формате JSON с заголовком `Content-Type: application/json; charset=utf-8`.

| Метод  | Путь                | Описание                     | Успешный ответ |
|--------|---------------------|------------------------------|----------------|
| `POST` | `/api/transactions` | Создать транзакцию           | `201 Created`  |
| `GET`  | `/api/transactions` | Список всех транзакций       | `200 OK`       |
| `POST` | `/api/budgets`      | Создать или обновить бюджет  | `201 Created`  |
| `GET`  | `/api/budgets`      | Список всех бюджетов         | `200 OK`       |

### Формат ошибок

Все ошибки возвращаются в едином формате:

```json
{"error": "описание ошибки"}
```

| Статус                      | Когда возникает                                                         |
|-----------------------------|-------------------------------------------------------------------------|
| `400 Bad Request`           | некорректный JSON, неверный формат даты, не пройдена валидация           |
| `409 Conflict`              | транзакция превышает бюджет категории (`{"error":"budget exceeded"}`)   |
| `500 Internal Server Error` | прочие внутренние ошибки                                                |

### POST /api/transactions

Тело запроса (дата в формате `YYYY-MM-DD`):

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount":450,"category":"еда","description":"ланч","date":"2025-09-10"}'
```

Ответ `201 Created`:

```json
{"id":1,"amount":450,"category":"еда","description":"ланч","date":"2025-09-10"}
```

Ошибка валидации, `400 Bad Request`:

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount":-5,"category":"еда","description":"ошибка","date":"2025-09-10"}'
```

```json
{"error":"amount must be greater than zero"}
```

### GET /api/transactions

```bash
curl http://localhost:8080/api/transactions
```

Ответ `200 OK` — массив транзакций. Если транзакций нет, возвращается пустой массив `[]`.

### POST /api/budgets

Повторный запрос с той же категорией обновляет лимит.

```bash
curl -X POST http://localhost:8080/api/budgets \
  -H "Content-Type: application/json" \
  -d '{"category":"транспорт","limit":1000}'
```

Ответ `201 Created`:

```json
{"category":"транспорт","limit":1000}
```

Превышение бюджета, `409 Conflict`:

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount":1500,"category":"транспорт","description":"такси","date":"2025-09-11"}'
```

```json
{"error":"budget exceeded"}
```

### GET /api/budgets

```bash
curl http://localhost:8080/api/budgets
```

Ответ `200 OK` — массив бюджетов (или `[]`). Порядок элементов не гарантирован, так как бюджеты хранятся в `map`.

### Проверка в Windows PowerShell

В PowerShell `curl` — это псевдоним `Invoke-WebRequest`, поэтому для примеров выше используйте `curl.exe` или Git Bash.

## Тесты

Тесты написаны на стандартной библиотеке `testing`, без сторонних зависимостей.

**Unit-тесты домена** (`ledger/ledger_test.go`):
- `TestTransactionValidate`, `TestBudgetValidate` — табличные тесты валидации (корректные данные, нулевые и отрицательные значения, пустая категория). Выполняются параллельно (`t.Parallel()`).
- `TestBudgetLimit` — сценарий «бюджет → транзакция в пределах лимита → превышение»: проверяется ошибка `ErrBudgetExceeded` и то, что отклонённая транзакция не сохранилась.

**Интеграционные тесты HTTP** (`gateway/internal/api/handlers_test.go`) прогоняют запросы через настоящий роутер (`NewRouter`) с помощью `httptest`:
- `TestBudgetsAPI` — подтесты `create`, `list`, `bad_json`, `invalid_limit`;
- `TestTransactionsAPI` — подтесты `ok`, `list`, `exceeded` (409), `bad_json` (400).

В каждом тесте проверяются статус, заголовок `Content-Type: application/json; charset=utf-8` и тело ответа. Перед тестами, меняющими данные, вызывается `ledger.Reset()`, а после — `t.Cleanup(ledger.Reset)`, чтобы сценарии не влияли друг на друга.

### Запуск

Из корня проекта (из-за `go.work` указываются имена модулей, а не `./...`):

```bash
go test gateway/... ledger/... -v
```

Покрытие и HTML-отчёт:

```bash
go test gateway/... ledger/... -cover -coverprofile=cover.out
go tool cover -func=cover.out
go tool cover -html=cover.out -o cover.html
```

> В PowerShell аргументы с `=` берите в кавычки: `"-coverprofile=cover.out"`, `"-func=cover.out"`, `"-html=cover.out"`.

**Покрытие тестами: ~77%** (`gateway/internal/api` — 87.1%, `ledger` — 62.2%).

## Логирование

Middleware `Logging` (`func(next http.Handler) http.Handler`) подключено ко всем маршрутам и выводит в консоль метод, путь и длительность обработки запроса:

```text
2026/10/29 00:13:50 GET /api/transactions 0s
2026/10/29 00:18:24 POST /api/transactions 0s
```

Запросы обрабатываются за микросекунды. На Windows разрешения таймера иногда не хватает, и длительность показывается как `0s`; на Linux/macOS выводится, например, `52µs`.

## Бизнес-логика Ledger

1. **Модели данных**:
   - `Transaction`: поля `ID`, `Amount`, `Category`, `Description`, `Date`.
   - `Budget`: поля `Category` и `Limit`.
2. **Хранение в памяти**:
   - Транзакции: срез `[]Transaction` с автоинкрементом идентификаторов.
   - Бюджеты: `map[string]Budget`, ключ — название категории.
3. **Контроль лимитов**: если для категории задан бюджет, `AddTransaction` суммирует траты с учётом новой операции. При превышении лимита возвращается сигнальная ошибка `ErrBudgetExceeded`, и транзакция не сохраняется. Gateway распознаёт её через `errors.Is` и отвечает `409 Conflict`.
4. **Импорт бюджетов (`LoadBudgets`)**: читает JSON-массив бюджетов из любого `io.Reader`, каждый бюджет проходит валидацию.
5. **Валидация (`Validatable`)**:

   ```go
   type Validatable interface {
       Validate() error
   }
   ```

   - `Transaction.Validate()`: `Amount > 0`, `Category` не пустая.
   - `Budget.Validate()`: `Limit > 0`, `Category` не пустая.
   - `CheckValid(v Validatable) error` — полиморфная проверка любого объекта, реализующего интерфейс.

### Формат файла `budgets.json`

```json
[
  { "category": "еда", "limit": 5000 },
  { "category": "развлечения", "limit": 2000 }
]
```
