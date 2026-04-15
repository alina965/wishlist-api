# Wishlist API

REST API сервис для управления вишлистами

## Запуск

```bash
docker-compose up --build
```

API доступно: `http://localhost:8080`  
Swagger документация: `http://localhost:8080/swagger/index.html`

## API Endpoints

### Публичные
| Method | Endpoint | Описание |
|--------|----------|----------|
| POST | `/api/register` | Регистрация |
| POST | `/api/login` | Вход |
| GET | `/api/wishlists/public?token={token}` | Просмотр вишлиста |
| GET | `/api/gifts?wishlist_id={id}` | Подарки вишлиста |
| POST | `/api/gifts/reserve` | Забронировать подарок |

### Защищенные (требуют JWT)
| Method | Endpoint | Описание |
|--------|----------|----------|
| POST | `/api/wishlists` | Создать вишлист |
| GET | `/api/wishlists` | Все вишлисты |
| GET | `/api/wishlists/get` | Вишлист по ID |
| PUT | `/api/wishlists/update` | Обновить вишлист |
| DELETE | `/api/wishlists/delete` | Удалить вишлист |
| POST | `/api/gifts` | Добавить подарок |
| DELETE | `/api/gifts` | Удалить подарок |

## Структура

```
├── cmd/main.go           # Точка входа
├── internal/
│   ├── domain/           # Модели
│   ├── application/      # Сервисы (бизнес-логика)
│   └── infrastructure/   # Репозитории, хендлеры, БД
├── migrations/           # SQL миграции
└── docker-compose.yml
```

## Unit-тесты для бизнес-логики

Написаны тесты для всех сервисов:

- **AuthService** - регистрация, логин, проверка паролей, дубликаты email
- **WishlistService** - создание, удаление, обновление, получение вишлистов, генерация токенов
- **GiftService** - создание, удаление, бронирование/разбронирование подарков, проверка приоритетов
