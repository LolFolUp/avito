# Avito.Kitchen MVP

Минимальная система заказа еды с публичным REST API для web-клиента и единственным приватным gRPC Partner API для ресторана. `kitchen-api` хранит каталог и заказы в одной PostgreSQL; `restaurant-demo` не имеет собственной БД, синхронизирует меню и polling'ом забирает заказы на подтверждение.

## Быстрый запуск

Требования: Docker с Compose v2. Для локальной генерации и тестов также нужны Go 1.25+ и [Task](https://taskfile.dev/).

```bash
cp .env.example .env
task up
```

После запуска:

- REST/OpenAPI: `http://localhost:8080`;
- приватный gRPC: `localhost:9090`;
- PostgreSQL: `localhost:5432` только для разработки;
- liveness/readiness: `GET /healthz`, `GET /readyz`.

Если локальный порт PostgreSQL занят, достаточно запустить `POSTGRES_PORT=55432 task up`; контейнеры всё равно обращаются к `postgres:5432`.

Основные команды:

```bash
task gen                 # ogen, protobuf/gRPC и GoMock
task db:up               # только PostgreSQL
task migrate             # применить Goose-миграции
task build               # собрать Docker-образы
task up                  # поднять полный стек
task logs                # смотреть логи
task test                # unit-тесты
task test:integration    # нужен TEST_DATABASE_URL и применённые миграции
task test:e2e            # нужен запущенный полный стек
task lint
task down
```

## Бизнес-правила MVP

- Web-клиент видит только активные рестораны и активные позиции меню. Флаг `is_available` показывается отдельно.
- Источник доступности — `restaurant-demo`: при `SyncMenu` он передаёт `is_active` и `is_available` в Kitchen. Отдельного Inventory Service нет.
- Синхронизация — полный снимок меню. Новые позиции создаются; изменённые обновляются с увеличением `version`; пропавшие из снимка деактивируются. `external_id` уникален внутри ресторана.
- Заказ относится к одному ресторану и содержит от 1 до 100 уникальных позиций, количество каждой — от 1 до 99.
- Web-клиент при создании заказа передаёт увиденную `observed_version`. Kitchen под блокировкой проверяет принадлежность позиции ресторану, активность, доступность, версию и валюту.
- Если меню изменилось или позиция стала недоступна, заказ не создаётся, REST возвращает `409 MENU_CHANGED` с актуальными сведениями о конфликтных позициях.
- Цена, название и внешний идентификатор позиции копируются в `order_items`. Последующие изменения меню не меняют уже созданный заказ.
- Суммы хранятся целыми числами в минимальных денежных единицах (`price_minor`), валюта заказа совпадает с валютой ресторана.
- Новый заказ получает `pending_confirmation`. `restaurant-demo` polling'ом получает такие заказы и принимает их, только если все внешние позиции доступны в его текущем in-memory меню; иначе отклоняет с причиной.
- Обновление статуса использует `expected_status`, поэтому конкурирующее или устаревшее обновление завершается `FailedPrecondition`.
- Допустимые переходы:

```text
pending_confirmation -> accepted -> preparing -> ready -> delivering -> delivered
                    \-> rejected
```

`rejected` и `delivered` терминальны. Причина обязательна только для `rejected`.

Границы транзакций задаёт service-слой через Avito Transaction Manager:

- полный `SyncMenu`;
- `CreateOrder` вместе с позициями и первой записью истории;
- смена статуса вместе с записью истории.

Простые GET-операции выполняются без транзакции. Репозитории используют `trmpgx.DefaultCtxGetter`: берут транзакцию из `context.Context`, а вне неё — `pgxpool`.

## Контракты

Исходники и generated-код версионированы и лежат рядом по назначению:

- REST: `shared/api/kitchen/v1/openapi.yaml` → `shared/pkg/openapi/kitchen/v1/` (ogen);
- gRPC: `shared/proto/partner/v1/partner.proto` → `shared/pkg/proto/partner/v1/` (Buf/protoc plugins);
- конфигурация: `shared/api/kitchen/v1/ogen.yaml`, `buf.yaml`, `buf.gen.yaml`.

Generated-файлы коммитятся, чтобы проект собирался без установленных генераторов и `protoc`. `task gen` остаётся единственным способом их обновления.

### Public REST API

| Метод | Путь | Назначение |
|---|---|---|
| `GET` | `/api/v1/restaurants` | Активные рестораны |
| `GET` | `/api/v1/restaurants/{restaurantId}/menu` | Меню ресторана с версиями и доступностью |
| `POST` | `/api/v1/orders` | Создать заказ, ответ `201` |
| `GET` | `/api/v1/orders/{orderId}` | Получить снимок и текущий статус заказа |
| `GET` | `/healthz` | Liveness процесса |
| `GET` | `/readyz` | Readiness с проверкой PostgreSQL |

Ошибки возвращаются как `application/problem+json`. Основной mapping: validation → `400`, not found → `404`, конфликт меню/состояния → `409`, неизвестная или инфраструктурная ошибка → `500` без SQL-деталей.

Пример создания заказа:

```bash
curl http://localhost:8080/api/v1/restaurants/00000000-0000-0000-0000-000000000001/menu

curl -X POST http://localhost:8080/api/v1/orders \
  -H 'Content-Type: application/json' \
  -d '{
    "restaurant_id":"00000000-0000-0000-0000-000000000001",
    "customer_name":"Иван",
    "customer_phone":"+79990000000",
    "delivery_address":"Москва",
    "items":[{"menu_item_id":"<id из меню>","quantity":1,"observed_version":1}]
  }'
```

### Private gRPC Partner API

| RPC | Назначение |
|---|---|
| `SyncMenu` | Передать полный актуальный снимок меню ресторана |
| `ListOrders` | Получить страницу заказов ресторана по статусу; demo запрашивает `PENDING_CONFIRMATION` |
| `UpdateOrderStatus` | Accept/reject и последующие статусы с optimistic guard `expected_status` |

Доменные ошибки единообразно отображаются в `InvalidArgument`, `NotFound`, `FailedPrecondition`, `AlreadyExists`, `Internal`. PostgreSQL-ошибки наружу не возвращаются.

## Архитектура

REST и gRPC — два входных адаптера одного `kitchen-api`, а не параллельные Partner API. В каждом сервисе соблюдается направление `api/client → service → repository`; модели transport, domain и PostgreSQL разделены converter-пакетами. `pkg/app/app.go` содержит только composition root и lifecycle серверов.

```mermaid
C4Container
    title Avito.Kitchen MVP — containers
    Person(customer, "Клиент", "Web-клиент")
    Person(restaurant, "Заведение", "Доступность и исполнение заказа")
    Container(kitchen, "kitchen-api", "Go, REST + gRPC", "Каталог, заказы, правила и транзакции")
    Container(demo, "restaurant-demo", "Go, gRPC client", "Меню в памяти, polling и решение accept/reject")
    ContainerDb(db, "PostgreSQL", "PostgreSQL 17", "Рестораны, меню, заказы, снимки и история")
    Rel(customer, kitchen, "Каталог, создание/получение заказа", "REST/OpenAPI :8080")
    Rel(demo, kitchen, "SyncMenu, ListOrders, UpdateOrderStatus", "private gRPC :9090")
    Rel(kitchen, db, "SQL", "pgxpool")
    Rel(restaurant, demo, "Управляет локальным меню")
```

Гибрид оправдан для MVP: web получает простой и обозримый HTTP-контракт, а интеграция сервис-сервис — строгий protobuf-контракт и сгенерированный клиент. Цена решения — два transport-стека, две генерации, gRPC health tooling и отсутствие удобного ручного вызова Partner API через обычный `curl`. Polling проще push/webhook и не требует доступности ресторана для входящих вызовов, но создаёт задержку до polling-интервала и лишние запросы без заказов.

### Структура

```text
kitchen/
  cmd/main.go
  internal/
    api/{converter,kitchen/v1,partner/v1}
    model/
    errs/
    service/{catalog,order}/
    repository/{kitchen,record,convert}/
  pkg/app/app.go
restaurant-demo/
  cmd/main.go
  internal/
    client/grpc/partner/v1/
    model/
    errs/
    service/restaurant/
    repository/memory/
  pkg/app/app.go
shared/
  api/kitchen/v1/
  proto/partner/v1/
  pkg/{openapi/kitchen/v1,proto/partner/v1}/
db/{migrations,schema.sql}
deploy/docker/
tests/e2e/
```

Интерфейсы зависимостей определяются потребителями в `deps.go`. Mocks находятся внутри соответствующих API/service-пакетов. Конкретные PostgreSQL repository и gRPC client остаются адаптерами.

## Данные и ER

```mermaid
erDiagram
    RESTAURANTS ||--o{ MENU_ITEMS : owns
    RESTAURANTS ||--o{ ORDERS : receives
    ORDERS ||--|{ ORDER_ITEMS : contains
    MENU_ITEMS ||--o{ ORDER_ITEMS : snapshotted_as
    ORDERS ||--|{ ORDER_STATUS_HISTORY : records

    RESTAURANTS {
        uuid id PK
        varchar name
        boolean is_active
        boolean accepting_orders
        varchar currency
        timestamptz created_at
        timestamptz updated_at
    }
    MENU_ITEMS {
        uuid id PK
        uuid restaurant_id FK
        varchar external_id UK
        varchar name
        bigint price_minor
        boolean is_active
        boolean is_available
        bigint version
    }
    ORDERS {
        uuid id PK
        uuid restaurant_id FK
        order_status status
        varchar customer_name
        varchar customer_phone
        varchar delivery_address
        bigint total_minor
        varchar currency
        varchar rejection_reason
    }
    ORDER_ITEMS {
        uuid id PK
        uuid order_id FK
        uuid menu_item_id FK
        varchar external_item_id
        varchar name_snapshot
        bigint unit_price_minor
        int quantity
        bigint line_total_minor
    }
    ORDER_STATUS_HISTORY {
        bigint id PK
        uuid order_id FK
        order_status from_status
        order_status to_status
        varchar reason
        timestamptz created_at
    }
```

Ограничения целостности находятся и в service, и в PostgreSQL. Индексы покрывают активное меню и основной polling-запрос `(restaurant_id, status, created_at, id)`. Последняя Goose-миграция добавляет только demo-ресторан; меню загружает `restaurant-demo` через Partner API.

## CJM

```mermaid
journey
    title Клиент
    section Выбор
      Открывает список ресторанов: 4: Клиент, Kitchen
      Смотрит актуальное меню: 4: Клиент, Kitchen
    section Заказ
      Отправляет корзину с версиями: 4: Клиент, Kitchen
      При конфликте обновляет меню: 2: Клиент, Kitchen
      Получает pending_confirmation: 4: Клиент, Kitchen
    section Исполнение
      Проверяет текущий статус заказа: 4: Клиент, Kitchen
      Получает delivered либо rejected: 3: Клиент, Kitchen
```

```mermaid
journey
    title Заведение
    section Подготовка
      Ведёт доступность в своей системе: 4: Заведение, Restaurant demo
      Синхронизирует полный снимок меню: 4: Restaurant demo, Kitchen
    section Новые заказы
      Polling pending_confirmation: 3: Restaurant demo, Kitchen
      Проверяет локальную доступность: 4: Restaurant demo
      Принимает или отклоняет заказ: 4: Restaurant demo, Kitchen
    section Исполнение
      Меняет статусы строго по графу: 4: Заведение, Restaurant demo, Kitchen
```

## Тестирование

- Unit: бизнес-валидация, переходы статусов, конфликт версий, error mapping и in-memory availability.
- Integration (`//go:build integration`): реальная PostgreSQL и проверка, что repository использует транзакцию из context и откатывает order/items/history вместе.
- E2E (`//go:build e2e`): REST-каталог → создание заказа → gRPC polling restaurant-demo → статус `accepted`.

Пример полного прогона:

```bash
task up
task test
TEST_DATABASE_URL='postgres://kitchen:kitchen@localhost:5432/kitchen?sslmode=disable' task test:integration
task test:e2e
task lint
```

## Осознанно вне MVP

API keys и transport security Partner API, access token публичного заказа, идемпотентность создания заказа, retries/backoff, lease/claim для polling, recovery после сбоев, observability и отдельный inventory-контур. Они нужны перед production, но не смешиваются с минимальной проверяемой бизнес-цепочкой этого задания.
