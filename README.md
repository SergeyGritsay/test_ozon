## Тестовое задание
### Задача

Необходимо создать сервис для хранения и подачи объявлений. Объявления должны храниться в базе данных. Сервис должен предоставлять API, работающее поверх HTTP в формате JSON.

### Настройка и запуск

Приложение и БД развернуты в контейнерах. Для удобства запуска используются команды make.

* Сборка контейнеров

```bash
make app-setup-and-up
```

После сборки контейнеров необходимо запустить миграцию для создания таблиц.
В данной реализации используется утилита goose.

* Применить миграцию

```bash
make db-migrate-up
```

* Запуск приложения

```bash
make app-up
```

### Описание методов

* Создать объявление
```http request
http://localhost:4000/api/advertisements [POST]
```
```json
{
	"title": "Вентилятор",
	"description": "Немного спасает от жары",
	"price": 1300,
	"photos": [
		"https://www.stadlerform.ru/upload/iblock/559/SF2_Otto_FULL_600_11.jpg",
		"https://cs.petrovich.ru/images/1980247/original-925x925-fit.jpg"
	]
}
```

Ответ:
```json
{
  "id": 3,
  "status": 200
}
```

* Получить объявление по id
```http request
http://localhost:4000/api/advertisements/:id [GET]
```

JSON Body (опционально):
```json
{
  "fields": [
    "description",
    "photos"
  ]
}
```

Ответ:
```json
{
  "description": "Немного спасает от жары",
  "photos": [
    "https://www.stadlerform.ru/upload/iblock/559/SF2_Otto_FULL_600_11.jpg",
    "https://cs.petrovich.ru/images/1980247/original-925x925-fit.jpg"
  ],
  "price": 1300,
  "title": "Вентилятор"
}
```

* Получить все объявления
```http request
http://localhost:4000/api/advertisements/ [GET]
```

JSON Body (опционально)
По умолчанию в качестве ответа возвращается 0-ая страница

Параметры sort:
* +createdAt
* -createdAt
* +price
* -price

```json
{
  "sort": "+createdAt",
  "page": 1
}
```

Ответ
```json
[
    {
      "photos": "https://www.stadlerform.ru/upload/iblock/559/SF2_Otto_FULL_600_11.jpg",
      "price": 1300,
      "title": "Вентилятор"
    },
    {
    "photos": "https://www.stadlerform.ru/upload/iblock/559/SF2_Otto_FULL_600_11.jpg",
    "price": 1200,
    "title": "Вентилятор 2"
    }
]
```

### Архитектура 

Приложение разделено на 4 слоя в соответствии с "Чистой архитектурой Дяди Боба":
* delivery - инфраструктурный слой
* repository - интерфейс работы с БД
* services (usecases) - бизнес-логика приложения
* models - сущности

Для реализации "Правила зависимостей" использовалась техника "Dependency Injection".

Стек:

* Golang
* Gin
* Goose
* PostgreSQL

### Дополнительно

Для удобства создания/запуска/отката миграций были добавлены соответствующие инструкции make:

* Создать миграцию
```bash
make db-migration-create name="NAME" lang="LANG"
```

* Применить миграции
```bash
make db-migrate-up
```

* Откат миграции
```bash
make db-migrate-down
```



