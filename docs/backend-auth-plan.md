# Аутентификация с GORM

## Summary

Реализовать backend-аутентификацию на Go с PostgreSQL и ORM `GORM v2`.

`GORM` выбран за зрелость, поддержку PostgreSQL/UUID/транзакций и удобную интеграцию с `pgx`. Миграции останутся отдельными SQL-миграциями через `golang-migrate`, чтобы схема БД не зависела от автогенерации ORM.

## Основные изменения

- Добавить зависимости:
  - `gorm.io/gorm`
  - `gorm.io/driver/postgres`
  - PostgreSQL-драйвер на базе `pgx`
  - `github.com/go-chi/chi/v5`
  - JWT, bcrypt и конфигурационные зависимости по необходимости.
- Реализовать GORM-модели:
  - `User`: UUID, email, bcrypt password hash, timestamps;
  - `RefreshSession`: UUID, user ID, SHA-256 token hash, expiry, revoked timestamp, timestamps.
- Настроить GORM через отдельный infrastructure/database-пакет:
  - подключение к PostgreSQL;
  - connection pooling;
  - проверка подключения;
  - отключить `AutoMigrate` в production;
  - использовать явные SQL-миграции `golang-migrate`.
- Реализовать repository-слой поверх GORM с контекстом во всех запросах.
- Для refresh rotation использовать транзакцию GORM:
  - проверить текущую сессию;
  - пометить ее отозванной;
  - создать новую сессию;
  - корректно обработать конкурентное повторное использование токена.
- Добавить маршруты:
  - `POST /auth/register`
  - `POST /auth/login`
  - `POST /auth/refresh`
  - `POST /auth/logout`
- Refresh token хранить только в `HttpOnly` cookie; access JWT возвращать в JSON.
- Использовать `bcrypt` для паролей и `HS256` для JWT.
- Секрет JWT, DSN и TTL получать из ENV.
- Добавить SQL-миграции и `docker-compose.yml` с PostgreSQL.
- Обновить `backend/README.md` с настройкой ENV, запуском PostgreSQL, миграций и сервера.
- Обновить `backend/AGENTS.md`:
  - зафиксировать использование GORM v2;
  - запретить `AutoMigrate` как замену production-миграциям;
  - описать разделение моделей БД и доменных типов;
  - указать правила repository-слоя, транзакций и context-aware запросов;
  - добавить команды миграций и тестирования;
  - сохранить требование тестов рядом с реализуемым поведением.

## Тесты

Добавить тесты на:

- регистрацию и уникальность email;
- валидацию email и непустой пароль;
- хранение только bcrypt-хэша;
- login с корректными и некорректными credentials;
- выдачу и проверку JWT claims;
- refresh по cookie;
- отказ для истекшего, отозванного и неизвестного refresh-токена;
- rotation и повторное использование старого токена;
- атомарность refresh-транзакции;
- logout и очистку cookie;
- HTTP-контракт через `httptest`;
- repository-операции GORM с PostgreSQL test database.

Проверки: `gofmt`, `go test ./...`, `go build ./...`.

## Предположения

- Используется `GORM v2` с PostgreSQL-драйвером на базе `pgx`.
- Миграции выполняются только через `golang-migrate`; GORM `AutoMigrate` не используется для production.
- Refresh-сессии stateful и ротируются при каждом обновлении.
- Минимальная политика пароля в первой версии: непустое значение.
- Восстановление пароля, подтверждение email, rate limiting и logout-all не входят в текущий scope.
