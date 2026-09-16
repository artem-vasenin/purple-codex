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

## Задачи

### Инфраструктура и конфигурация

- [x] Определить структуру backend-пакетов для domain, application, transport и infrastructure-слоев.
- [x] Добавить зависимости GORM v2, PostgreSQL-драйвера, `chi`, JWT, bcrypt и миграций.
- [x] Реализовать загрузку конфигурации из ENV: DSN, адрес сервера, JWT secret, access TTL, refresh TTL и cookie-настройки.
- [x] Добавить проверку обязательных production-настроек и безопасные dev-default значения TTL.
- [x] Реализовать подключение GORM к PostgreSQL с настройкой connection pool.
- [x] Добавить graceful shutdown HTTP-сервера и закрытие database connection.
- [x] Добавить `docker-compose.yml` для локального PostgreSQL.
- [x] Добавить команды или инструкции запуска и применения миграций.

### Схема данных и repositories

- [x] Спроектировать таблицы `users` и `refresh_sessions`.
- [x] Создать SQL-миграцию `up` для таблиц, индексов, ограничений и внешнего ключа.
- [x] Создать SQL-миграцию `down`.
- [x] Реализовать GORM-модели `User` и `RefreshSession`.
- [ ] Разделить persistence-модели и доменные типы, если они различаются по ответственности.
- [x] Реализовать user repository с поиском по нормализованному email и созданием пользователя.
- [x] Реализовать refresh-session repository: создание, поиск, отзыв и удаление/очистку просроченных сессий.
- [x] Добавить context-aware запросы во все repository-методы.
- [x] Реализовать транзакционный метод rotation refresh-сессии с защитой от повторного использования старого токена.
- [x] Убедиться, что `AutoMigrate` не используется вместо production-миграций.

### Auth domain и token service

- [x] Реализовать нормализацию и валидацию email.
- [x] Реализовать проверку непустого пароля.
- [x] Реализовать bcrypt-хэширование и проверку пароля.
- [x] Реализовать генерацию криптографически случайного refresh-токена.
- [x] Реализовать SHA-256 хэширование refresh-токена перед сохранением.
- [x] Реализовать выпуск access JWT с claims `sub`, `iat`, `exp` и `jti`.
- [x] Реализовать проверку алгоритма `HS256`, подписи и срока действия JWT.
- [x] Реализовать use cases регистрации, входа, обновления токенов и logout.
- [x] Скрыть различия между неизвестным email и неверным паролем в публичных auth-ошибках.

### HTTP API

- [x] Поднять HTTP router на `chi`.
- [x] Реализовать `POST /auth/register`.
- [x] Реализовать `POST /auth/login`.
- [x] Реализовать `POST /auth/refresh`.
- [x] Реализовать `POST /auth/logout`.
- [x] Реализовать единый JSON-формат успешных ответов и ошибок.
- [x] Установить refresh cookie как `HttpOnly` с корректными `Secure`, `SameSite`, `Path` и `MaxAge`.
- [x] Возвращать access JWT в JSON без сохранения его на сервере.
- [x] Обработать malformed JSON, отсутствующие поля, истекшие токены и внутренние ошибки без утечки деталей.
- [x] Добавить базовый health endpoint для проверки доступности сервиса и БД.

### Тестирование

- [x] Добавить unit-тесты нормализации и валидации email и пароля.
- [x] Добавить unit-тесты bcrypt-хэширования и проверки пароля.
- [ ] Добавить unit-тесты генерации и валидации JWT.
- [ ] Добавить unit-тесты хэширования и генерации refresh-токенов.
- [x] Добавить тест регистрации нового пользователя.
- [x] Добавить тест отказа при повторном email.
- [x] Добавить тест login с корректными credentials.
- [x] Добавить тест одинаковой ошибки для неизвестного email и неверного пароля.
- [ ] Добавить тест refresh по валидной cookie.
- [ ] Добавить тест отказа для истекшего, отозванного и неизвестного refresh-токена.
- [ ] Добавить тест rotation и повторного использования старого токена.
- [ ] Добавить тест атомарности refresh-транзакции и конкурентного запроса.
- [ ] Добавить тест logout и очистки cookie.
- [ ] Добавить HTTP-тесты маршрутов через `httptest`.
- [ ] Добавить repository-тесты GORM с PostgreSQL test database.
- [x] Запустить `gofmt`, `go test ./...` и `go build ./...`.

### Документация и правила разработки

- [x] Обновить `backend/README.md` с ENV, локальным PostgreSQL, миграциями и запуском сервера.
- [x] Обновить `backend/AGENTS.md` правилами GORM v2, repository-слоя и транзакций.
- [ ] Документировать API-контракт auth-маршрутов и формат cookie.
- [ ] Документировать production-требования: HTTPS, `Secure` cookie, JWT secret и CORS/CSRF policy.
- [ ] Проверить, что секреты и локальные environment-файлы не добавлены в репозиторий.

## Предположения

- Используется `GORM v2` с PostgreSQL-драйвером на базе `pgx`.
- Миграции выполняются только через `golang-migrate`; GORM `AutoMigrate` не используется для production.
- Refresh-сессии stateful и ротируются при каждом обновлении.
- Минимальная политика пароля в первой версии: непустое значение.
- Восстановление пароля, подтверждение email, rate limiting и logout-all не входят в текущий scope.
