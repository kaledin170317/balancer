# HTTP Load Balancer


## Запуск с нуля через Docker

```bash
# Клонировать проект
git clone https://github.com/kaledin170317/balancer.git
cd balancer

# Перейти в папку с Docker-файлами
cd build

# Поднять всю инфраструктуру
docker-compose up --build
```

## Структура проекта

```
.
├── build/                  # Dockerfile и docker-compose.yml
├── cmd/
│   └── balancer/           # main.go
├── config/                 # Загрузка конфигурации
├── docs/                   # Задание
├── internal/
│   ├── adapters/
│   │   ├── api/rest/controllers      # CRUD /clients
│   │   ├── api/rest/middleware       # Middleware
│   │   ├── api/rest/errors           # JSON ошибки
│   │   └── db/postgreSQL             # Хранилище клиентов
│   ├── balancer/
│   │   ├── algoritms                 # round-robin, random, least-connections
│   │   ├── checker                   # health-check
│   │   └── proxy                     # Reverse proxy
│   ├── config                        # Загрузка переменных из env/file/cli
│   ├── domain/
│   │   ├── models                    # Client, Backend
│   │   └── usecases                  # ClientUseCase + sync.Map
│   ├── logger                        # slog JSON логгер
│   └── ratelimit                     # TokenBucket
├── migrations/             # SQL-миграции
├── go.mod
├── go.sum
└── README.md
```

---

## Конфигурация через переменные окружения

```env
LISTEN_ADDR=:8080
BACKENDS=http://localhost:9001,http://localhost:9002
ALGORITHM=round-robin
RATE_CAPACITY=100
RATE_REFILL=10
HC_INTERVAL=5s
HC_TIMEOUT=2s
DB_DSN=postgres://postgres:password@postgres:5432/balancer?sslmode=disable
```

---


После запуска:

- API доступно на `http://localhost:8080`
- PostgreSQL работает на `localhost:5555`

---

## Примеры запросов

```bash
curl -X POST http://localhost:8080/clients -H "Content-Type:application/json" -d "{"clientId":"user123","capacity":100,"ratePerSec":10}"
curl http://localhost:8080/clients/user123
curl -X PUT http://localhost:8080/clients/user123 -H "Content-Type:application/json" -d "{"capacity":200,"ratePerSec":20}"
curl -X DELETE http://localhost:8080/clients/user123
curl http://localhost:8080/ -H "X-Client-ID:user123"
```

---

## JSON ошибки

```json
{ "code": 429, "message": "Rate limit exceeded" }
{ "code": 502, "message": "Backend error" }
{ "code": 503, "message": "No healthy backend" }
```