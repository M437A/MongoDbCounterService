## Сервис подсчета документов MongoDB

Этот репозиторий содержит простой сервис для подсчета документов в коллекциях MongoDB.

### Установка

1. Клонируйте репозиторий:
git clone https://github.com/M437A/MongoDbCounterService.git


2. Запустите контейнер MongoDB. Для этого можете воспользоваться скриптом `run.sh`.

3. Запустите приложение.

### Эндпоинты

#### 1. Получение данных из всех коллекций во всех базах данных

- **Метод:** GET
- **URL:** `http://localhost:8080/`
- **Описание:** Получает данные из всех коллекций во всех базах данных.
- **Ответ:**
```json
[
   {
       "CollectionName": "startup_log",
       "DocumentCount": 3
   },
   {
       "CollectionName": "system.version",
       "DocumentCount": 2
   },
   {
       "CollectionName": "system.users",
       "DocumentCount": 1
   },
   {
       "CollectionName": "system.sessions",
       "DocumentCount": 0
   }
]
```


#### 2) Получение данных о коллекциях по имени базы данных 

- **Метод:** GET
- **URL:** `http://localhost:8080/database`
- **Описание:** Получает данные из всех коллекций конкретной базы данных.
  
- **Пример запроса:**
```json
{
    "database": "admin"
}
```

- **Ответ:**
```json
[
    {
        "CollectionName": "system.version",
        "DocumentCount": 2
    },
    {
        "CollectionName": "system.users",
        "DocumentCount": 1
    }
]
```
