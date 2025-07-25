# Zloy Billing Service

Микросервис для управления биллингом и пользователями. Обеспечивает аутентификацию с CAPTCHA, управление балансом пользователей и работу с отчетами.

## Архитектура

### Clean Architecture
Проект построен по принципам Clean Architecture:
- **Domain** - бизнес-логика и интерфейсы
- **UseCase** - бизнес-правила приложения  
- **Repository** - доступ к данным
- **Delivery** - HTTP handlers и middleware
- **Config** - централизованная конфигурация

### Базы данных

**PostgreSQL** - пользователи и баланс:
- Таблица `users`: id, login, password_hash, balance, created_at
- Баланс в центах для точности (1000 центов = $10)
- Пароли хешируются bcrypt

**MongoDB** - отчеты:
- Коллекция `reports`: _id, report_id, user_id, client_generated_id, is_purchased, created_at
- Поддержка анонимных отчетов (user_id = null)
- Уникальные report_id

### Конфигурация через переменные окружения

Все настройки вынесены в `.env`:
- `REPORT_COST_CENTS=1000` - стоимость отчета
- `JWT_EXPIRATION_HOURS=24` - время жизни токена
- `CAPTCHA_LIFETIME_MINUTES=5` - время жизни CAPTCHA
- `DEFAULT_PAGE_LIMIT=10` - лимит пагинации

## Запуск проекта

### Требования
- Docker и Docker Compose

### Не забудьте установить jq

### Запуск
```bash
# Клонирование и запуск
git clone <repository-url>
cd Zloy
docker-compose up -d --build
migrate-up
```

Сервис доступен: **http://localhost:8080**

## API Endpoints

### 🔐 Аутентификация


#### Получение CAPTCHA(в json файле возьмите id сессии с png прочитайте код)
```bash
curl -X GET "http://localhost:8080/api/captcha/generate" > fresh_captcha.json && cat fresh_captcha.json | jq -r '.image' | base64 -d > fresh_captcha.png
```
**Ответ:**
```json
{
  "image": "base64_png_image_data",
  "session_id": "captcha_1753453172755014131"
}
```

#### Регистрация пользователя
```bash
curl -X POST "http://localhost:8080/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "login": "newuser2024",
    "password": "testpass123",
    "captcha_id": "captcha_1753453172755014131",
    "captcha_code": "3625"
  }' | jq .
```
**Ответ:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 4,
    "login": "newuser2024"
  }
}
```

#### Вход в систему
```bash
curl -X POST "http://localhost:8080/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "login": "newuser2024", 
    "password": "testpass123"
  }' | jq .
```

### 📊 Отчеты (требуют аутентификации)

#### Создание mock отчета (для тестирования)
```bash
curl -X POST "http://localhost:8080/api/mock/create-report" \
  -H "Content-Type: application/json" \
  -d '{
    "client_generated_id": "user-2024-test-001"
  }' | jq .
```
**Ответ:**
```json
{
  "message": "Mock report created",
  "report_id": "d0379cab-41f0-4cc1-828a-9d13b818cb6d"
}
```

#### Привязка анонимных отчетов к пользователю
```bash
TOKEN="your_jwt_token_here"

curl -X POST "http://localhost:8080/api/user/link-anonymous" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "client_generated_id": "user-2024-test-001"
  }' | jq .
```
**Ответ:**
```json
{
  "message": "Reports linked successfully",
  "updated": 1
}
```

#### Получение отчетов пользователя
```bash
curl -X GET "http://localhost:8080/api/user/reports?limit=10&offset=0" \
  -H "Authorization: Bearer $TOKEN" | jq .
```
**Ответ:**
```json
{
  "reports": [
    {
      "_id": "688392ae2c85076810a41b1e",
      "report_id": "d0379cab-41f0-4cc1-828a-9d13b818cb6d",
      "user_id": 4,
      "client_generated_id": "user-2024-test-001",
      "is_purchased": false,
      "created_at": "2025-07-25T14:20:30.148Z"
    }
  ],
  "total": 1
}

#пополняем баланс
curl -X POST "http://localhost:8080/api/user/topup?amount=2500" \
  -H "Authorization: Bearer $TOKEN" | jq .

#чекаем баланс
curl -X GET "http://localhost:8080/api/user/balance" \
  -H "Authorization: Bearer $TOKEN" | jq .

```
#### Покупка отчета
```bash
# Успешная покупка (при достаточном балансе)(подставляем report_id)
curl -X POST "http://localhost:8080/api/reports/d0379cab-41f0-4cc1-828a-9d13b818cb6d/purchase" \
  -H "Authorization: Bearer $TOKEN" | jq .
```
**Ответ при успехе:**
```json
{
  "message": "Report purchased successfully"
}
```

**Возможные ошибки:**
- `"Insufficient balance"` - недостаточно средств
- `"Report not found"` - отчет не найден  
- `"Report already purchased"` - отчет уже куплен


## 🔒 Безопасность

- **CAPTCHA валидация** при регистрации
- **JWT токены** с настраиваемым временем жизни
- **bcrypt хеширование** паролей
- **Middleware аутентификации** для защищенных endpoints
- **Проверка баланса** перед покупками
- **Защита от повторных покупок**

## ⚙️ Технологии

- **Go 1.23** + Clean Architecture
- **Chi Router** - HTTP маршрутизация
- **JWT** - аутентификация
- **PostgreSQL** - пользователи и транзакции
- **MongoDB** - отчеты и документы
- **Docker** - контейнеризация
- **bcrypt** - безопасность паролей

## 📁 Структура проекта

```
Zloy/
├── cmd/main.go                    # Точка входа
├── internal/
│   ├── config/config.go           # Конфигурация
│   ├── domain/                    # Бизнес-логика
│   │   ├── user.go
│   │   ├── report.go
│   │   ├── auth.go
│   │   └── errors.go
│   ├── usecase/                   # Бизнес-правила
│   │   ├── user_usecase.go
│   │   ├── auth_usecase.go
│   │   ├── report_usecase.go
│   │   └── captcha_usecase.go
│   ├── repository/                # Доступ к данным
│   │   ├── postgres_user_repository.go
│   │   └── mongo_report_repository.go
│   └── delivery/                  # HTTP слой
│       ├── http/
│       │   ├── auth_handler.go
│       │   ├── report_handler.go
│       │   ├── captcha_handler.go
│       │   └── dto.go
│       └── middleware/
│           └── auth_middleware.go
├── migrations/postgres/           # SQL миграции
├── docker-compose.yml            # Docker окружение
├── Dockerfile                    # Сборка образа
└── Makefile                      # Команды разработки
```

## 🛠️ Команды разработки

```bash
#быстрый старт
make start
migrate-up
##логи
logs-app: ## Логи только приложения
logs: ## Показать логи Docker контейнеров
logs-db: ## Логи баз данных
```

## 🌍 Переменные окружения

Скопируйте `env.example` в `.env` и настройте:

```bash
# Порт сервиса
PORT=8080

# MongoDB
MONGO_URI=mongodb://localhost:27017

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=zl0y_billing

# JWT Secret (в продакшене должен быть сложным)
JWT_SECRET=your-secret-key-change-in-production 
```

## ✅ Протестированные функции

- ✅ Генерация и валидация CAPTCHA
- ✅ Регистрация с проверкой дубликатов
- ✅ Аутентификация и JWT токены  
- ✅ Создание и привязка анонимных отчетов
- ✅ Пагинация списка отчетов
- ✅ Покупка отчетов с проверкой баланса
- ✅ Защита от повторных покупок
- ✅ Обработка ошибок и валидация
- ✅ Middleware аутентификации
- ✅ Конфигурация через env переменные 
