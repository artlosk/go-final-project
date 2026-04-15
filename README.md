# Планировщик задач (TODO-лист)

Итоговый проект: веб-сервер TODO-листа

## Что реализовано

- создание/инициализация SQLite БД при первом запуске;
- API:
  - `GET /api/nextdate`,
  - `POST /api/task`,
  - `GET /api/task`,
  - `PUT /api/task`,
  - `DELETE /api/task`,
  - `GET /api/tasks`,
  - `POST /api/task/done`;
- поддержка авторизации:
  - `POST /api/signin`,
  - проверка JWT в middleware для защищённых API;
- поддержка правила повторения: `d`, `y`, `w`, `m`

## Задания со звездочкой

- Реализованы правила повторения (`w`, `m`)
- Реализована авторизация через пароль (`TODO_PASSWORD`) и JWT (если пароль не установлен, то авторизация не требуется)
- Добавлен `Dockerfile`

## Переменные окружения

- `TODO_PORT` — порт сервера (по умолчанию `7540`);
- `TODO_DBFILE` — путь к файлу SQLite;
- `TODO_PASSWORD` — пароль для входа:
  - пустой: авторизация отключена,
  - непустой: требуется вход через `/login.html`

## Локальный запуск

1. Подготовить `.env`:

```bash
cp .env.example .env
```

Пример:

```env
TODO_PORT=7540
TODO_DBFILE=./data/scheduler.db
TODO_PASSWORD=
```

2. Запустить сервер:

```bash
go run ./cmd
```

3. Открыть в браузере:

`http://localhost:7540/`
если `TODO_PASSWORD` пароль не пустой, будет редирект `http://localhost:7540/login.html`

## Запуск тестов

Обычный запуск:

```bash
go test -v ./tests
```
С очисткой кэша
```bash
go test -count=1 -v ./tests
```

Очистка кэша тестов Go:

```bash
go clean -testcache
```

## Docker

Через Docker Compose + docker exec:

```bash
docker compose up -d --build
docker exec -it app /usr/local/go/bin/go test -count=1 -v ./tests
docker compose down
```

Параметры в `tests/settings.go`:

- `Port` — порт тестируемого сервера
- `DBFile` — путь к тестовой БД
- `FullNextDate` — расширенные тесты `w/m`
- `Search` — тесты поиска
- `Token` — JWT токен при включенной авторизации

Каталог `data` уже есть в репозитории; БД `data/scheduler.db` создаётся автоматически.

Docker реализован не одним файлом, а через контейнер, чтобы можно было прогонять тесты через тот же контейнер

## Структура проекта

- `cmd`
  - `main.go` — точка входа приложения
- `internal`
  - `api` — регистрация роутов
  - `auth` — работа с JWT
  - `config` — чтение конфигурации из окружения
  - `db` — инициализация БД и SQL запросы
  - `handlers` — обработчики API
  - `helpers` — общие вспомогательные функции
  - `middleware` — middleware (для авторизации)
  - `schedule` — логика расчёта следующей даты
  - `server` — настройка и запуск сервера
- `tests` — интеграционные тесты проекта
- `web` — фронтенд (HTML/CSS/JS)
- `data` — файл SQLite БД для локального запуска/контейнера
- `.env.example` — пример переменных окружения
- `docker-compose.yml` — запуск проекта в контейнере
- `Dockerfile` — сборка docker-образа
- `README.md` — документация по проекту