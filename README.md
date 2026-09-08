# Avito.Kitchen MVP

Система заказа еды с публичным REST/OpenAPI контрактом и приватным gRPC-контрактом между сервисами.

- `kitchen-api` хранит реестр ресторанов и публичные snapshots заказов.
- `restaurant-demo` владеет меню, доступностью, остатками и ресторанными заказами.
- Kitchen не хранит `menu_items` и не синхронизирует меню.

Основные потоки:

```text
Оформление заказа
Web -> REST Kitchen -> gRPC RestaurantDemo -> PostgreSQL ресторана
                         |
                         +-> при успехе Kitchen сохраняет snapshot заказа
                             в своей PostgreSQL

Исполнение заказа
Заведение -> gRPC RestaurantDemo -> PostgreSQL ресторана
                         |
                         +-> синхронный gRPC в Kitchen обновляет
                             публичный статус и его историю
```

## Быстрый запуск

Требования: Docker Compose v2. Для генерации и тестов также нужны Go 1.25+ и [Task](https://taskfile.dev/).

```bash
cp .env.example .env
task up
```

После запуска:

- Kitchen REST: `http://localhost:8080`;
- Kitchen партнёрский gRPC: `localhost:9090`;
- Restaurant gRPC: `localhost:9091`;
- Kitchen PostgreSQL: `localhost:15432`;
- Restaurant PostgreSQL: `localhost:15433`;
- Kitchen health: `GET /healthz`, `GET /readyz`.

Основные команды:

```bash
task gen
task db:up
task migrate
task build
task up
task logs
task test
task test:integration
task test:e2e
task lint
task down
```

## Бизнес-правила

- Kitchen получает меню live-вызовом `RestaurantService.GetMenu`.
- Restaurant возвращает активные блюда, включая временно недоступные; клиент видит `is_available`.
- UUID и `version` блюда принадлежат restaurant-demo. Отдельный `external_id` между сервисами не используется.
- Заказ относится к одному ресторану, содержит 1–100 уникальных блюд и 1–99 единиц каждого блюда.
- Web возвращает `observed_version`, увиденную при чтении меню.
- `CreateRestaurantOrder` в одной транзакции блокирует блюда через `FOR UPDATE`, проверяет существование, версию, доступность, валюту и остаток, уменьшает остаток и записывает ресторанный заказ.
- Если хотя бы одна проверка не прошла, вся транзакция откатывается, а Kitchen возвращает `409 MENU_CHANGED`.
- При успехе restaurant-demo возвращает авторитетные snapshots названий и цен. Kitchen сохраняет их в своих `orders/order_items`; дальнейшее изменение меню не изменяет заказ.
- Успешно созданный заказ сразу имеет статус `ACCEPTED`: отдельного polling-подтверждения больше нет.
- Денежные суммы хранятся целыми числами в минимальных единицах валюты.

Внутри каждого сервиса соблюдается направление `api/client -> service -> repository`. Transport DTO, domain model и PostgreSQL record разделены converter-пакетами. Транзакции открывает только service через Avito Transaction Manager; repository выбирает транзакцию из context посредством `trmpgx.DefaultCtxGetter`.

## Контракты

- REST-исходники: `shared/api/kitchen/v1/`;
- ogen-код: `shared/pkg/openapi/kitchen/v1/`;
- Protobuf-исходник: `shared/proto/partner/v1/partner.proto`;
- protobuf/gRPC-код: `shared/pkg/proto/partner/v1/`;
- генерация: `task gen`.

### Public REST

| Метод | Путь | Назначение |
|---|---|---|
| `GET` | `/api/v1/restaurants` | Получить активные рестораны |
| `GET` | `/api/v1/restaurants/{restaurantId}/menu` | Получить live-меню через restaurant-demo |
| `POST` | `/api/v1/orders` | Атомарно проверить корзину и создать заказ |
| `GET` | `/api/v1/orders/{orderId}` | Получить сохранённый snapshot заказа |
| `GET` | `/healthz` | Проверить, что процесс работает |
| `GET` | `/readyz` | Проверить готовность Kitchen и его PostgreSQL |

Ошибки возвращаются как `application/problem+json`: validation — `400`, not found — `404`, конфликт меню — `409`, недоступный restaurant-demo — `503`, неизвестная ошибка — `500`. SQL-ошибки наружу не передаются.

### Private gRPC

Оба направления интеграции описаны в одном версионированном Protobuf-контракте:

| Сервис и RPC | Кто вызывает | Назначение |
|---|---|---|
| `RestaurantService.GetMenu` | Kitchen | Прочитать актуальное меню из PostgreSQL ресторана |
| `RestaurantService.CreateRestaurantOrder` | Kitchen | Атомарно проверить, зарезервировать блюда и создать заказ |
| `RestaurantService.UpdateRestaurantOrderStatus` | система ресторана / demo-вызов | Изменить локальный статус и синхронно уведомить Kitchen |
| `KitchenPartnerService.UpdateOrderStatus` | restaurant-demo | Сохранить новый статус и history в Kitchen |

Kitchen использует сгенерированный `RestaurantServiceClient`, а restaurant-demo —
`KitchenPartnerServiceClient`.

`KitchenPartnerService` закрыт API key-аутентификацией. Restaurant-demo передаёт ключ
в gRPC metadata `x-api-key`; unary interceptor Kitchen сопоставляет ключ с UUID
заведения. UUID из тела запроса дополнительно сверяется с аутентифицированным UUID,
поэтому подставить идентификатор другого ресторана недостаточно. Для локального Compose
ключ задаётся переменной `RESTAURANT_API_KEY`, а Kitchen получает отображение
`restaurant_id=key` через `RESTAURANT_API_KEYS`.

## Архитектура

Диаграммы описаны на Mermaid и рендерятся GitHub без внешних инструментов, поэтому
исходник диаграммы и её отображение всегда совпадают.

### C4 Level 1 — контекст

```mermaid
flowchart TB
    Client["Клиент Авито.Кухня<br/><i>[Person]</i><br/>Выбирает заведение,<br/>собирает корзину, следит за заказом"]
    Partner["Заведение<br/><i>[Person / внешняя система]</i><br/>Держит меню и остатки,<br/>исполняет заказ"]

    System["<b>Авито.Кухня</b><br/><i>[Software System]</i><br/>Каталог заведений, оформление заказов<br/>и отслеживание их состояния"]

    Client -->|"Смотрит меню,<br/>оформляет заказ"| System
    System -->|"Передаёт заказ,<br/>получает подтверждение"| Partner
    Partner -->|"Сообщает новый статус заказа"| System
```

### C4 Level 2 — контейнеры

```mermaid
flowchart TB
    Web["Web-клиент<br/><i>[SPA, вне репозитория]</i>"]

    subgraph platform["Авито.Кухня — граница системы"]
        KitchenAPI["<b>kitchen-api</b><br/><i>[Container: Go, ogen, gRPC]</i><br/>Каталог заведений, оформление заказов,<br/>публичные snapshots, партнёрский API"]
        KitchenDB[("<b>Kitchen PostgreSQL</b><br/><i>[Container: PostgreSQL 17]</i><br/>restaurants, orders,<br/>order_items, order_status_history")]
    end

    subgraph venue["Заведение — отдельный деплой"]
        Demo["<b>restaurant-demo</b><br/><i>[Container: Go, gRPC]</i><br/>Меню, остатки, резервирование,<br/>исполнение заказа"]
        DemoDB[("<b>Restaurant PostgreSQL</b><br/><i>[Container: PostgreSQL 17]</i><br/>menu_items, orders,<br/>order_items, order_status_history")]
    end

    Web -->|"HTTP/JSON :8080<br/>OpenAPI kitchen/v1"| KitchenAPI
    KitchenAPI -->|"pgx :5432"| KitchenDB
    KitchenAPI -->|"gRPC :9091<br/>GetMenu, CreateRestaurantOrder"| Demo
    Demo -->|"gRPC :9090<br/>UpdateOrderStatus"| KitchenAPI
    Demo -->|"pgx :5432"| DemoDB
```

Оба направления приватные и синхронные. Kitchen одновременно REST-сервер для web,
gRPC-клиент заведения и gRPC-сервер партнёрского API; restaurant-demo — gRPC-сервер
для Kitchen и gRPC-клиент его партнёрского API.

### C4 Level 3 — компоненты kitchen-api

```mermaid
flowchart TB
    subgraph transport["api — транспортный слой"]
        Rest["<b>api/kitchen/v1</b><br/><i>[ogen Handler]</i><br/>Публичный REST для web"]
        PartnerSrv["<b>api/partner/v1</b><br/><i>[gRPC Server]</i><br/>Партнёрский API заведений"]
        ApiConv["<b>api/converter</b><br/>DTO ⇄ model,<br/>problem+json"]
    end

    subgraph domain["service — бизнес-логика"]
        Catalog["<b>service/catalog</b><br/>Рестораны, live-меню"]
        Orders["<b>service/order</b><br/>Валидация, allowedTransitions,<br/>границы транзакций"]
    end

    subgraph infra["repository / client — инфраструктура"]
        RepoConv["<b>repository/convert</b><br/>model ⇄ record"]
        Repo["<b>repository/kitchen</b><br/><i>[pgx + transaction manager]</i>"]
        Client["<b>client/grpc/partner</b><br/>Resolver по restaurant_id"]
    end

    KitchenDB[("Kitchen PostgreSQL")]
    Demo["restaurant-demo"]

    Rest --> ApiConv
    PartnerSrv --> ApiConv
    ApiConv --> Catalog
    ApiConv --> Orders
    Catalog --> Repo
    Catalog --> Client
    Orders --> Repo
    Orders --> Client
    Repo --> RepoConv
    Repo --> KitchenDB
    Client --> Demo
```

Направление зависимостей всегда `api/client -> service -> repository`. Transport DTO,
domain model и PostgreSQL record изолированы converter-пакетами: `api/converter`
переводит DTO ⇄ model, `repository/convert` — model ⇄ record. Ни один слой не видит
типы соседа через голову конвертера.

### Состояния заказа

```mermaid
stateDiagram-v2
    [*] --> ACCEPTED: CreateRestaurantOrder успешен
    ACCEPTED --> PREPARING: заведение начало готовить
    PREPARING --> READY: заказ собран
    READY --> DELIVERING: передан в доставку
    DELIVERING --> DELIVERED: вручён клиенту
    DELIVERED --> [*]
```

Переходы проверяются в service-слоях обеих систем. Если ресторан не может создать
заказ из-за отсутствия блюда, количества или изменения версии меню, создание
атомарно отклоняется и public REST возвращает `409`; отдельный заказ `REJECTED`
в текущем синхронном MVP не сохраняется.

### Схема данных

Kitchen PostgreSQL:

- `restaurants` — реестр доступных заведений и адрес интеграции на уровне конфигурации MVP;
- `orders` — публичное состояние заказа;
- `order_items` — snapshots UUID, названия, цены и количества;
- `order_status_history` — история статусов.

Restaurant PostgreSQL:

- `menu_items` — меню, доступность, версия и `available_quantity`;
- `orders` — заказ, принятый рестораном;
- `order_items` — ресторанные snapshots;
- `order_status_history` — история исполнения.

Между `Kitchen.order_items.menu_item_id` и Restaurant PostgreSQL физического внешнего ключа нет: это межсервисная ссылка.

```mermaid
erDiagram
    KITCHEN_RESTAURANTS {
        uuid id PK
        varchar name
        boolean is_active
        boolean accepting_orders "требует is_active"
        varchar currency
        timestamptz created_at
        timestamptz updated_at
    }
    KITCHEN_ORDERS {
        uuid id PK "общий с restaurant-demo"
        uuid restaurant_id FK
        order_status status
        varchar customer_name
        varchar customer_phone
        varchar delivery_address
        bigint total_minor "минимальные единицы валюты"
        varchar currency
        varchar rejection_reason "NOT NULL только для REJECTED"
        timestamptz created_at
        timestamptz updated_at
    }
    KITCHEN_ORDER_ITEMS {
        uuid order_id PK,FK
        uuid menu_item_id PK "логическая ссылка в другую БД"
        varchar name_snapshot
        bigint unit_price_minor
        integer quantity "CHECK 1..99"
        bigint line_total_minor "GENERATED ALWAYS"
    }
    KITCHEN_STATUS_HISTORY {
        bigserial id PK
        uuid order_id FK
        order_status status
        varchar reason
        timestamptz created_at
    }

    RESTAURANT_MENU_ITEMS {
        uuid id PK
        varchar name
        varchar description
        bigint price_minor
        varchar currency
        boolean is_active
        boolean is_available
        integer available_quantity
        bigint version "оптимистичная блокировка"
    }
    RESTAURANT_ORDERS {
        uuid id PK "приходит из Kitchen"
        restaurant_order_status status
        varchar customer_name
        varchar customer_phone
        varchar delivery_address
        bigint total_minor
        varchar currency
        timestamptz created_at
        timestamptz updated_at
    }
    RESTAURANT_ORDER_ITEMS {
        uuid order_id PK,FK
        uuid menu_item_id PK,FK
        varchar name_snapshot
        bigint unit_price_minor
        integer quantity
        bigint line_total_minor "GENERATED ALWAYS"
    }
    RESTAURANT_STATUS_HISTORY {
        bigserial id PK
        uuid order_id FK
        restaurant_order_status from_status
        restaurant_order_status to_status
        timestamptz created_at
    }

    KITCHEN_RESTAURANTS ||--o{ KITCHEN_ORDERS : receives
    KITCHEN_ORDERS ||--|{ KITCHEN_ORDER_ITEMS : contains
    KITCHEN_ORDERS ||--|{ KITCHEN_STATUS_HISTORY : records

    RESTAURANT_ORDERS ||--|{ RESTAURANT_ORDER_ITEMS : contains
    RESTAURANT_MENU_ITEMS ||--o{ RESTAURANT_ORDER_ITEMS : snapshots
    RESTAURANT_ORDERS ||--|{ RESTAURANT_STATUS_HISTORY : records

    RESTAURANT_MENU_ITEMS }o..o{ KITCHEN_ORDER_ITEMS : "logical menu_item_id"
    KITCHEN_ORDERS }o..o{ RESTAURANT_ORDERS : "общий order id"
```

Ключевые решения по схеме:

- `KITCHEN_ORDERS.id` и `RESTAURANT_ORDERS.id` — один и тот же UUID, сгенерированный
  Kitchen. В будущем его можно использовать как ключ идемпотентности; полноценная
  обработка повторной отправки в MVP не реализована.
- `order_items` в обеих базах — snapshots названия и цены на момент оформления.
  Последующее изменение меню не переписывает историю заказа.
- `line_total_minor` в обеих базах — генерируемая колонка, рассинхрон с
  `unit_price_minor * quantity` невозможен по построению.
- Первичный ключ `order_items` — составной `(order_id, menu_item_id)`: суррогатный
  `id` не нёс смысла, а уникальность пары и так требовалась.
- Деньги везде целые в минимальных единицах валюты, валюта хранится рядом с ценой.

Один экземпляр `restaurant-demo` представляет один ресторан, поэтому локальная
таблица `restaurants` и повторяющиеся `restaurant_id` в таблицах ресторана не
нужны. Валюта хранится вместе с ценами позиций меню и snapshot заказа, а не в
конфигурации сервиса.

## Транзакции и согласованность

Restaurant выполняет в одной локальной транзакции:

```text
SELECT menu_items FOR UPDATE
-> validate version/availability/quantity
-> decrement available_quantity
-> INSERT orders
-> INSERT order_items
-> INSERT order_status_history(ACCEPTED)
```

Kitchen после успешного gRPC-ответа отдельной локальной транзакцией записывает свой
`order + items + history`. При смене статуса restaurant-demo сначала атомарно меняет
свой `order + history`, затем синхронно вызывает Kitchen. Общей ACID-транзакции между
двумя PostgreSQL нет. Поэтому возможны orphan после создания или временное расхождение
статусов при сбое между локальным commit и внешним вызовом. Для MVP это известный
trade-off; production-улучшения — idempotency key/order ID, retry, recovery job и
transactional outbox.

## CJM

Два пути: клиента, который заказывает еду, и заведения, которое этот заказ исполняет.

### CJM клиента

```mermaid
journey
    title Путь клиента: от выбора заведения до полученного заказа
    section Выбор заведения
      Открывает список заведений: 4: Клиент
      Видит, кто принимает заказы прямо сейчас: 5: Клиент
    section Сбор корзины
      Открывает меню заведения: 5: Клиент
      Видит цену, состав и доступность блюда: 4: Клиент
      Собирает корзину: 5: Клиент
    section Оформление
      Отправляет заказ: 3: Клиент
      Узнаёт, что блюдо разобрали, пока он выбирал: 1: Клиент
      Обновляет меню и повторяет заказ: 2: Клиент
      Получает подтверждение и сумму: 5: Клиент
    section Ожидание
      Следит за статусом заказа: 4: Клиент
      Получает заказ: 5: Клиент
```

Болевая точка — шаг «блюдо разобрали, пока он выбирал». Она решена не блокировкой
корзины, а `observed_version`: клиент присылает версию меню, которую видел, и Kitchen
отвечает `409 MENU_CHANGED` с перечнем конфликтных позиций, их актуальными ценами и
доступностью. Клиент видит именно то, что изменилось, а не абстрактную ошибку.

```mermaid
sequenceDiagram
    autonumber
    actor Web as Web-клиент
    participant K as kitchen-api
    participant KDB as Kitchen PostgreSQL
    participant R as restaurant-demo
    participant RDB as Restaurant PostgreSQL

    Web->>K: GET /api/v1/restaurants
    K->>KDB: SELECT активные заведения
    KDB-->>K: список
    K-->>Web: 200 RestaurantList

    Web->>K: GET /api/v1/restaurants/{id}/menu
    K->>R: gRPC GetMenu
    R->>RDB: SELECT menu_items
    RDB-->>R: позиции с version
    R-->>K: MenuItem[]
    K-->>Web: 200 Menu (is_available, version)

    Note over Web: клиент собирает корзину<br/>и запоминает observed_version

    Web->>K: POST /api/v1/orders
    K->>KDB: SELECT заведение
    Note over K: проверка is_active<br/>и accepting_orders
    K->>R: gRPC CreateRestaurantOrder

    alt меню изменилось или не хватает остатка
        R->>RDB: SELECT menu_items FOR UPDATE
        RDB-->>R: version не совпала
        Note over R: rollback транзакции,<br/>остаток не тронут
        R-->>K: FailedPrecondition
        K-->>Web: 409 MENU_CHANGED + конфликтные позиции
    else всё доступно
        R->>RDB: FOR UPDATE, списание остатка,<br/>INSERT orders/items/history
        RDB-->>R: commit
        R-->>K: RestaurantOrder + snapshots цен
        K->>KDB: INSERT orders/items/history в своей транзакции
        KDB-->>K: commit
        K-->>Web: 201 Order (ACCEPTED)
    end

    Web->>K: GET /api/v1/orders/{orderId}
    K->>KDB: SELECT order + items
    KDB-->>K: snapshot
    K-->>Web: 200 Order с текущим статусом
```

### CJM заведения

```mermaid
journey
    title Путь заведения: от подключения до закрытого заказа
    section Подключение
      Получает UUID заведения и адрес приватного gRPC: 3: Заведение
      Публикует меню, цены и остатки: 4: Заведение
    section Приём заказа
      Отдаёт актуальное меню по запросу Kitchen: 5: Заведение
      Получает заказ с observed_version: 4: Заведение
      Отклоняет заказ по устаревшей версии: 2: Заведение
      Резервирует остаток и принимает заказ: 5: Заведение
    section Исполнение
      Начинает готовить и сообщает PREPARING: 4: Заведение
      Собирает заказ и сообщает READY: 4: Заведение
      Передаёт в доставку: 4: Заведение
      Закрывает заказ как DELIVERED: 5: Заведение
```

Заведение никогда не пишет в базу Kitchen напрямую: оно меняет статус у себя и
сообщает о нём партнёрским вызовом. Kitchen остаётся единственным владельцем своего
состояния и сам проверяет допустимость перехода.

```mermaid
sequenceDiagram
    autonumber
    actor Venue as Персонал заведения
    participant R as restaurant-demo
    participant RDB as Restaurant PostgreSQL
    participant K as kitchen-api
    participant KDB as Kitchen PostgreSQL
    actor Web as Web-клиент

    Note over R,K: заказ уже создан и находится в ACCEPTED

    Venue->>R: перевести заказ в PREPARING
    R->>RDB: SELECT order FOR UPDATE
    RDB-->>R: текущий статус
    Note over R: проверка перехода<br/>у себя
    R->>RDB: UPDATE orders + INSERT history
    RDB-->>R: commit

    R->>K: gRPC UpdateOrderStatus<br/>(restaurant_id, order_id,<br/>expected_status, new_status)
    Note over K: проверка restaurant_id,<br/>принадлежности заказа<br/>и допустимости перехода

    alt статус в Kitchen не совпал с expected_status
        K-->>R: FailedPrecondition STATUS_CONFLICT
        Note over R: заведение перечитывает<br/>заказ и повторяет
    else переход разрешён
        K->>KDB: UPDATE orders + INSERT history<br/>в одной транзакции
        KDB-->>K: commit
        K-->>R: обновлённый заказ
    end

    Web->>K: GET /api/v1/orders/{orderId}
    K-->>Web: 200 Order (PREPARING)

    Note over Venue,Web: далее READY -> DELIVERING -> DELIVERED
```

## Тестирование

- Unit: валидация и отсутствие резервирования при конфликте версии.
- Integration: реальные PostgreSQL и транзакции обоих сервисов.
- E2E: live-меню, успешное создание заказа через gRPC и `409` при устаревшей версии.

```bash
task up
task test
task test:integration
task test:e2e
```

## Осознанные упрощения

Решения, принятые ради объёма MVP. Каждое видно в коде или схеме, поэтому описано явно.

**Отказ ресторана не сохраняется как состояние заказа.** Если позиции разобрали или
изменилась версия меню, `CreateRestaurantOrder` откатывает транзакцию, заказ не
создаётся ни в одной из баз, а клиент получает `409 MENU_CHANGED`. Отдельного статуса
`REJECTED` у сохранённого заказа нет.

В Kitchen зарезервированы ещё два состояния — `PENDING_CONFIRMATION` для будущего
асинхронного подтверждения через outbox и `REJECTED` для отказа уже принятого заказа.
Они пока остаются в доменной модели, OpenAPI, enum базы и заготовке машины переходов,
но текущий runtime-поток их не устанавливает: без outbox первое недостижимо, а второе
требует компенсирующего возврата `available_quantity`, иначе отклонённый заказ навсегда
списал бы продукты со склада.

Для этого будущего расширения также сохранены:

- значения `PENDING_CONFIRMATION` и `REJECTED` в enum `order_status` — PostgreSQL не
  умеет `ALTER TYPE ... DROP VALUE`, а пересоздавать тип через `ALTER TABLE ... USING`
  ради двух неиспользуемых меток дороже, чем оставить их;
- колонка `orders.rejection_reason`, констрейнт `orders_rejection_reason_consistent` и
  колонка `order_status_history.reason` — они не заполняются и ждут появления отмены
  заказа с возвратом резерва.

**Согласованность между двумя базами.** Описана выше в разделе «Транзакции и
согласованность»: общей ACID-транзакции нет, окна рассогласования known trade-off.

**Синхронная доставка статусов.** Заведение вызывает Kitchen напрямую и получает
ошибку, если вызов не прошёл. Повтор того же перехода идемпотентен, но автоматического
retry нет.

## Вне текущего MVP

Ротация API keys и mTLS, access token заказа, полноценная идемпотентность, retries/backoff,
автоматический recovery, transactional outbox как надёжная доставка статусов вместо
текущего синхронного вызова, отмена заказа с возвратом зарезервированного количества.
