# Avito.Kitchen MVP

Сервис доставки еды: клиент выбирает заведение, собирает корзину из его живого меню и
оформляет заказ, а заведение исполняет этот заказ у себя и сообщает статус обратно.

В репозитории два сервиса, каждый со своей БД(выбрал PostgreSQL):

- **`kitchen-api`** - сама площадка. Держит реестр заведений, оформляет заказы,
  хранит их публичные snapshots(для истории заказов) и историю статусов. Наружу отдаёт REST по OpenAPI
  для web-клиента, внутрь - приватный gRPC для заведений.
- **`restaurant-demo`** — пример подключённого заведения. Владеет меню, ценами,
  доступностью и остатками, принимает окончательное решение по заказу и ведёт его
  исполнение.

Оба поднимаются одной командой через Docker Compose(ниже сама команда).

## Быстрый запуск

```bash
cp .env.example .env
task up
```

После запуска:

- Kitchen REST: `http://localhost:8080`;
- Kitchen партнёрский gRPC: `http://localhost:9090`;
- Restaurant gRPC: `http://localhost:9091`;
- Kitchen PostgreSQL: `http://localhost:15432`;
- Restaurant PostgreSQL: `http://localhost:15433`;
- health-проверки Kitchen:
  `GET /health`, - обычный http запрос к kitchen, возвращает 200 если все работает
  `GET /ready`.  - проверка работоспособности kitchen. Делается пинг в БД

Основные команды:

```bash
task up          # поднять всё через Compose
task down        # остановить
task logs        # логи сервисов
task gen         # перегенерировать REST, gRPC и моки из контрактов
task migrate     # накатить миграции обеих баз
task lint        # запустить линтер
```

## CJM

Два пути:
1) Клиент заказывает еду
2) Заведение исполняет заказ

### CJM клиента

```mermaid
journey
    title Путь клиента: от выбора заведения до полученного заказа
    section Выбор заведения
      Открывает список заведений: 5: Клиент
      Видит, кто принимает заказы прямо сейчас: 5: Клиент
    section Сбор корзины
      Открывает меню заведения: 5: Клиент
      Видит цену, состав и доступность блюда: 5: Клиент
      Собирает корзину: 5: Клиент
    section Оформление
      Отправляет заказ: 5: Клиент
      Узнаёт, что блюдо разобрали, пока он выбирал: 1: Клиент
      Обновляет меню и повторяет заказ: 1: Клиент
      Получает подтверждение и сумму: 5: Клиент
    section Ожидание
      Следит за статусом заказа: 5: Клиент
      Получает заказ: 5: Клиент
```

Как я решил проблему устаревшей корзины - изменения цены или доступности блюда после его добавления.
Она решена не блокировкой корзины, а с помощью `observed_version`: при чтении меню клиент получает `version` каждой позиции
и возвращает её при оформлении. Если меню успело измениться, заказ не создаётся частично - вся проверка откатывается, и Kitchen отвечает `409 MENU_CHANGED`.

Разбивки по конкретным позициям в `409` нет: ответ сообщает только факт конфликта.
Что именно изменилось, клиент узнаёт из повторного `GET .../menu`
(клиент сам поймет, где именно цена изменилась, или позиция уже закончилась, когда будет просматривать меню)
там уже видны актуальные цены, `is_available` и новые `version`. Это осознанное упрощение MVP,
см. «Осознанные упрощения»(не успел).

Не знаю стоило ли показывать все возможные ошибки, которые пользователь может получить в CJM(например, при POST /api/v1/orders, когда у нас is_active/acceptings_orders == false, возвращает клиенту 409)
так что не стал показывать, кроме транзакций, на схеме покажу все явно.

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
    K->>KDB: SELECT заведение по id
    KDB-->>K: restaurant (is_active)
    Note over K: проверка is_active
    K->>R: gRPC GetMenu
    R->>RDB: SELECT menu_items
    RDB-->>R: позиции с version
    R-->>K: MenuItem[]
    K-->>Web: 200 Menu (is_available, version)

    Note over Web: клиент собирает корзину<br/>и запоминает observed_version

    Web->>K: POST /api/v1/orders
    K->>KDB: SELECT заведение по id
    KDB-->>K: restaurant (is_active, accepting_orders, currency)
    Note over K: проверка is_active<br/>и accepting_orders
    K->>R: gRPC CreateRestaurantOrder
    R->>RDB: BEGIN транзакции Restaurant
    R->>RDB: SELECT menu_items FOR UPDATE
    RDB-->>R: актуальные позиции

    alt меню изменилось или не хватает остатка
        Note over R: version не совпала<br/>или недостаточно остатка
        R->>RDB: ROLLBACK Restaurant
        RDB-->>R: остаток не изменён
        R-->>K: FailedPrecondition
        K-->>Web: 409 MENU_CHANGED
    else всё доступно
        R->>RDB: UPDATE остатка +<br/>INSERT orders/items/history
        RDB-->>R: операции выполнены
        R->>RDB: COMMIT Restaurant
        RDB-->>R: транзакция зафиксирована
        R-->>K: RestaurantOrder + snapshots цен
        K->>KDB: BEGIN транзакции Kitchen
        K->>KDB: INSERT orders/items/history
        KDB-->>K: операции выполнены
        K->>KDB: COMMIT Kitchen
        KDB-->>K: транзакция зафиксирована
        K-->>Web: 201 Order (ACCEPTED)
    end


    Web->>K: GET /api/v1/orders/{orderId}
    K->>KDB: SELECT order + items
    KDB-->>K: snapshot
    K-->>Web: 200 Order с текущим статусом
```

### CJM заведения

```mermaid
%%{init: {"journey": {"width": 230, "height": 70, "taskFontSize": 14}}}%%
journey
    title Путь заведения: от подключения до закрытого заказа
    section Подключение
      Получает UUID заведения и адрес приватного gRPC: 3
      Публикует меню, цены и остатки: 4
    section Приём заказа
      Отдаёт актуальное меню по запросу Kitchen: 5
      Получает заказ с observed_version: 4
      Отклоняет заказ по устаревшей версии: 2
      Резервирует остаток и принимает заказ: 5
    section Исполнение
      Начинает готовить и сообщает PREPARING: 4
      Собирает заказ и сообщает READY: 4
      Передаёт в доставку: 4
      Закрывает заказ как DELIVERED: 5
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

## Архитектура

Диаграммы описаны на Mermaid и рендерятся GitHub без внешних инструментов, поэтому
исходник диаграммы и её отображение всегда совпадают.

### C4 Level 1 — контекст

```mermaid
flowchart TB
    Client["Клиент Авито.Кухня<br/><i>[Person]</i><br/>Выбирает заведение,<br/>собирает корзину,<br/>следит за заказом"]
    Partner["Заведение<br/><i>[Person / внешняя система]</i><br/>Держит меню и остатки,<br/>исполняет заказ"]

    System["<b>Авито.Кухня</b><br/><i>[Software System]</i><br/>Каталог заведений,<br/>оформление заказов<br/>и отслеживание их состояния"]

    Client -->|"Смотрит меню,<br/>оформляет заказ"| System
    System -->|"Запрашивает меню,<br/>передаёт заказ,<br/>получает подтверждение"| Partner
    Partner -->|"Сообщает новый статус заказа"| System
```

### C4 Level 2 — контейнеры

```mermaid
flowchart TB
    Web["Web-клиент<br/><i>[SPA, вне репозитория]</i>"]
    Staff["Персонал заведения<br/><i>[Person]</i>"]

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
    Staff -->|"gRPC :9091<br/>UpdateRestaurantOrderStatus"| Demo
    Demo -->|"gRPC :9090<br/>UpdateOrderStatus<br/>x-api-key"| KitchenAPI
    Demo -->|"pgx :5432"| DemoDB
```

Оба направления приватные и синхронные. Kitchen одновременно REST-сервер для web,
gRPC-клиент заведения и gRPC-сервер партнёрского API; restaurant-demo — gRPC-сервер
для Kitchen и gRPC-клиент его партнёрского API.

Доверие одностороннее: партнёрский API Kitchen требует `x-api-key` и по нему определяет
заведение, а вызовы Kitchen в заведение ключом не подписываются. Kitchen защищает
запись в своё состояние, а обращения к заведению в MVP считаются доверенными внутри
приватной сети; при выходе за её пределы это направление тоже придётся подписывать.

Меню целиком принадлежит заведению: Kitchen не хранит `menu_items` и ничего не
синхронизирует, а на каждый запрос меню ходит в restaurant-demo через
`GetMenu`. Компромисс - доступность меню зависит от заведения.

### C4 Level 3 — компоненты kitchen-api

```mermaid
flowchart TB
    subgraph transport["api — транспортный слой"]
        Rest["<b>api/kitchen/v1</b><br/><i>[ogen Handler]</i><br/>Публичный REST для web"]
        PartnerSrv["<b>api/partner/v1</b><br/><i>[gRPC Server]</i><br/>Партнёрский API заведений"]
        Auth["<b>api/partner/v1/auth</b><br/><i>[Unary Interceptor]</i><br/>x-api-key ⇒ restaurant_id"]
        ApiConv["<b>api/converter</b><br/>DTO ⇄ model,<br/>problem+json"]
    end

    subgraph domain["service — бизнес-логика"]
        Catalog["<b>service/catalog</b><br/>Рестораны, live-меню"]
        Orders["<b>service/order</b><br/>Валидация, allowedTransitions,<br/>границы транзакций<br/><i>[Avito Transaction Manager]</i>"]
    end

    subgraph infra["repository / outbound client — инфраструктура"]
        RepoConv["<b>repository/convert</b><br/>model ⇄ record"]
        Repo["<b>repository/kitchen</b><br/><i>[pgx + CtxGetter]</i>"]
        Client["<b>client/grpc/partner</b><br/><i>[исходящий gRPC-клиент]</i><br/>Resolver по restaurant_id"]
    end

    KitchenDB[("Kitchen PostgreSQL")]
    Demo["restaurant-demo"]

    Auth --> PartnerSrv
    Rest --> ApiConv
    PartnerSrv --> ApiConv
    Rest --> Catalog
    Rest --> Orders
    PartnerSrv --> Orders
    Catalog --> Repo
    Catalog --> Client
    Orders --> Repo
    Orders --> Client
    Repo --> RepoConv
    Repo --> KitchenDB
    Client --> Demo
    Rest -.->|"GET /ready:<br/>Ping пула мимо repository"| KitchenDB
```

### C4 Level 3 — компоненты restaurant-demo

```mermaid
flowchart TB
    subgraph transport["api — транспортный слой"]
        PartnerSrv["<b>api/partner/v1</b><br/><i>[gRPC Server]</i><br/>RestaurantService:<br/>GetMenu, CreateRestaurantOrder,<br/>UpdateRestaurantOrderStatus"]
        ApiConv["<b>api/converter</b><br/>DTO ⇄ model,<br/>errs ⇄ gRPC status"]
    end

    subgraph domain["service — бизнес-логика"]
        Restaurant["<b>service/restaurant</b><br/>Проверка observed_version,<br/>резерв остатков, расчёт суммы,<br/>allowedStatusTransitions,<br/>границы транзакций<br/><i>[Avito Transaction Manager]</i>"]
    end

    subgraph infra["repository / outbound client — инфраструктура"]
        RepoConv["<b>repository/convert</b><br/>model ⇄ record"]
        Repo["<b>repository/restaurant</b><br/><i>[pgx + CtxGetter]</i><br/>menu_items SELECT ... FOR UPDATE"]
        Client["<b>client/grpc/kitchen/v1</b><br/><i>[исходящий gRPC-клиент]</i><br/>x-api-key в metadata"]
    end

    DemoDB[("Restaurant PostgreSQL")]
    Kitchen["kitchen-api"]
    Staff["Персонал заведения<br/><i>[Person]</i>"]

    Kitchen -->|"gRPC :9091<br/>GetMenu,<br/>CreateRestaurantOrder"| PartnerSrv
    Staff -->|"gRPC :9091<br/>UpdateRestaurantOrderStatus"| PartnerSrv
    PartnerSrv --> ApiConv
    PartnerSrv --> Restaurant
    Restaurant --> Repo
    Restaurant --> Client
    Repo --> RepoConv
    Repo --> DemoDB
    Client -->|"gRPC :9090<br/>UpdateOrderStatus"| Kitchen
```

## Инженерные решения

Что сделано осознанно и где это видно в коде.

### Чистая трёхслойная архитектура

Основной поток вызовов в каждом сервисе: `api` -> `service` ->
`repository/client`. API использует converter для transport DTO, repository -
для record-моделей. Каждый слой не знает внутренностей другого слоя, а для взаимодействия
использует интерфейсы `deps.go`.

Интерфейсы объявлены на стороне потребителя - `service/order/deps.go`,
`api/kitchen/v1/deps.go`, `service/catalog/deps.go`. Каждый слой описывает ровно те методы, которые ему нужны. 
Из этих же файлов генерируются заглушки (`task gen:mocks`), поэтому любую зависимость можно
подменить, не трогая продакшн-код.

### У каждого слоя своя модель

Их три, и они не пересекаются:

| Слой | Тип | Откуда берётся |
|---|---|---|
| транспорт | ogen DTO, protobuf-сообщения | сгенерированы из OpenAPI и `.proto` |
| домен | `internal/model` | написаны руками, ни от чего не зависят |
| хранилище | `internal/repository/record` | повторяют колонки таблиц |

Границы держат конвертеры: `api/converter` переводит DTO <-> model,
`repository/convert` — model <-> record. Ни один слой не видит типы соседа, 
поэтому переименование колонки не протекает в JSON, а правка OpenAPI не
ломает SQL. 


### Как написан cmd файлы.

`cmd/main.go` отвечает только за процесс: конфиг, логгер, сигналы, код возврата.                                                                                                   
Зависимости и серверы собираются в `app.New`, поэтому приложение остаётся                                                                                                    
пакетом, который можно поднять из теста. Такой подход как я слышал назвается
библиотечным; он легче читается и легче тестируется

### Транзакции открывает только service

Repository не начинает и не коммитит транзакцию - он выполняет один запрос в том
соединении, которое ему дали. Границу задаёт service, потому что только он знает,
какие изменения обязаны примениться вместе.

Границ ровно две, обе многошаговые:

- **Kitchen** — `orders` + `order_items` + `order_status_history` одной транзакцией
  (`service/order/create.go`, `service/order/update_status.go`);
- **restaurant-demo** — `SELECT menu_items FOR UPDATE` -> проверка версии, доступности
  и остатка -> списание остатка -> `orders` + `order_items` + `order_status_history`
  (`service/restaurant/create_order.go`). Блокировка строк меню на весь этот путь и
  делает проверку корзины атомарной: две параллельные попытки забрать последнюю
  порцию не могут обе пройти.

### Менеджер транзакций от Авито

Используется [`github.com/avito-tech/go-transaction-manager`](https://github.com/avito-tech/go-transaction-manager) (узнал о нем, когда смотрел разборы ошибок за прошлые года, надеюсь кто-то дойдет до этой readMe части):
service вызывает `txManager.Do(ctx, func(txCtx context.Context) error { ... })`, а
repository внутри достаёт соединение через `r.getter.DefaultTrOrDB(ctx, r.pool)`.

Что это даёт по сравнению с ручной передачей `pgx.Tx` параметром:

- **Не надо создавать две отдельные версии с транзакциями и без транзакцией**.
    Да и в принципе не надо думать, используется ли здесь транзакция или не используется ли вовсе, так как
    все зависит от ctx(если содержит транзакция будет pgx.Tx, если нет, то просто r.pool);
- **`DefaultTrOrDB` возвращает либо активную транзакцию из контекста, либо пул**, если
  транзакции нет. Один и тот же `GetOrder` одинаково работает внутри `Do` и снаружи;
- **commit и rollback привязаны к возврату из функции**: вернулась ошибка - rollback,
  вернулся `nil` - commit. Закоммититься дважды или
  выйти из функции с открытой транзакцией здесь нереально.

### доменные ошибки и логирование

Внутри сервисов ходит `errs.Error` с `Kind` и стабильным `Code`. В HTTP-код и
gRPC-статус он превращается ровно на границе: `api/kitchen/v1/errors.go` и
`api/partner/v1/errors.go`. Наружу уходит `application/problem+json` со стабильным
машинным `code`, а тексты драйвера и SQL-ошибки не покидают repository - они
заворачиваются в `repositoryError` и логируются, но не отдаются клиенту.

### Использование OpenAPI и Proto контрактов

REST описан в `shared/api/kitchen/v1/`, межсервисный gRPC - в
`shared/proto/partner/v1/partner.proto`. Серверы, клиенты и DTO генерируются из них
(`ogen`, `buf`). Схема не может разойтись с кодом, потому что код
из неё и получается.



### Деньги и статическая проверка

Все суммы - целые в минимальных единицах валюты (590 ₽ хранятся как `59000`), валюта
лежит рядом с ценой. Плавающей точки в денежных путях нет вообще; переполнение при
умножении цены на количество проверяется явно.

Статический анализ настроен в `.golangci.yml` и запускается как `task lint` вместе с
`go vet` и `buf lint`.

## Контракты

- REST-исходники: `shared/api/kitchen/v1/`;
- ogen-код: `shared/pkg/openapi/kitchen/v1/`;
- Protobuf-исходник: `shared/proto/partner/v1/partner.proto`;
- protobuf/gRPC-код: `shared/pkg/proto/partner/v1/`;
- генерация: `task gen`.

### Public REST

| Метод | Путь | Назначение |
|---|---|---|
| `GET` | `/api/v1/restaurants` | Получить активные заведения |
| `GET` | `/api/v1/restaurants/{restaurantId}/menu` | Получить live-меню через restaurant-demo |
| `POST` | `/api/v1/orders` | Атомарно проверить корзину и создать заказ |
| `GET` | `/api/v1/orders/{orderId}` | Получить сохранённый snapshot заказа |
| `GET` | `/health` | Проверить, что процесс работает |
| `GET` | `/ready` | Проверить готовность Kitchen и его PostgreSQL |

Заказ относится к одному заведению и содержит от 1 до 100 уникальных позиций по
1–99 единиц каждой. Меню возвращает все активные блюда, включая временно недоступные:
клиент видит их с `is_available: false` и понимает, почему не может их заказать.

Ошибки приходят как `application/problem+json`: валидация — `400`, не найдено — `404`,
конфликт меню - `409`, restaurant-demo недоступен — `503`, всё прочее — `500`.

### Private gRPC

Оба направления интеграции описаны в одном версионированном Protobuf-контракте:

| Сервис и RPC | Кто вызывает | Назначение |
|---|---|---|
| `RestaurantService.GetMenu` | Kitchen | Прочитать актуальное меню из PostgreSQL заведения |
| `RestaurantService.CreateRestaurantOrder` | Kitchen | Атомарно проверить, зарезервировать позиции и создать заказ |
| `RestaurantService.UpdateRestaurantOrderStatus` | система заведения / demo-вызов | Изменить локальный статус и синхронно уведомить Kitchen |
| `KitchenPartnerService.UpdateOrderStatus` | restaurant-demo | Сохранить новый статус и историю в Kitchen |

`KitchenPartnerService` закрыт API key-аутентификацией. Restaurant-demo передаёт ключ в
gRPC metadata `x-api-key`; unary interceptor Kitchen сопоставляет ключ с UUID заведения.
UUID из тела запроса дополнительно сверяется с аутентифицированным, поэтому подставить
чужой идентификатор недостаточно. Для локального Compose ключ задаётся переменной
`RESTAURANT_API_KEY`, а Kitchen получает отображение `restaurant_id=key` через
`RESTAURANT_API_KEYS`.

## Схема БД

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

Переходы проверяются в service-слоях обеих систем. Успешно созданный заказ сразу
получает `ACCEPTED`. Если заведение не может
принять заказ, создание атомарно отклоняется и заказ не сохраняется нигде.

### Таблицы

Kitchen PostgreSQL:

- `restaurants` - реестр доступных заведений;
- `orders` - публичное состояние заказа;
- `order_items` - snapshots UUID, названия, цены и количества;
- `order_status_history` - история статусов.

Restaurant PostgreSQL:

- `menu_items` - меню, доступность, версия и `available_quantity`;
- `orders` - заказ, принятый заведением;
- `order_items` - snapshots заведения;
- `order_status_history` - история исполнения.

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


Миграции лежат в `db/migrations/kitchen` и `db/migrations/restaurant`, накатываются
goose и запускаются в Compose до старта сервисов.

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
свой `order + history`, затем синхронно вызывает Kitchen.

Общей ACID-транзакции между двумя PostgreSQL нет, и это главный компромисс системы.
Поэтому возможны orphan-заказы после создания (заведение списало остаток, Kitchen не
успел записать snapshot) или временное расхождение статусов при сбое между локальным
commit и внешним вызовом. См.Осознанные упрощения

## Осознанные упрощения(или то что не успел сделать)

- **Ошибки `409 Conflict`.** Используются коды `MENU_CHANGED` и
  `RESTAURANT_NOT_ACCEPTING_ORDERS`. Для `MENU_CHANGED` не возвращается список
  конфликтных позиций: клиент повторно получает меню и пересобирает корзину.

- **Статусы.** В рабочем сценарии используются переходы
  `ACCEPTED → PREPARING → READY → DELIVERING → DELIVERED`.
  Статусы `PENDING_CONFIRMATION` и `REJECTED` зарезервированы для будущего
  асинхронного подтверждения и отмены заказа, но сейчас не устанавливаются.

- **Transactional Outbox** не реализован. Kitchen и restaurant-demo используют
  синхронные gRPC-вызовы, поэтому между фиксацией изменений в двух базах остаётся
  небольшое окно рассогласования(Например, restaurant-demo может зафиксировать заказ и уменьшить остатки, после чего
  транзакция Kitchen завершится ошибкой. Тогда заказ останется в restaurant-demo,
  но клиент получит ошибку и будет считать, что заказ не создан).

- **DDD.** Использовано разделение на `model`, `service` и `repository`, однако
  полноценные агрегаты, value objects, domain events и отдельный domain-слой
  для MVP не вводились.

## Масштабирование

- Исправление всех выше упущений
- Реализация Transactional Outbox и Kafka для асинхронной и надёжной передачи
  заказов и изменений статусов. Kitchen будет сохранять заказ со статусом
  `PENDING_CONFIRMATION` и событие в outbox одной транзакцией, после чего отправлять
  его в Kafka. Это позволит повторять доставку при сбоях, сглаживать пики нагрузки
  (так как у нас не статично повышаются rps, а динамически с пиками).
- Обработчики событий должны быть идемпотентными. Конечно можно взять 
  гарантию доставки exactly once, но это дороже и сложнее.
