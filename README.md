# go-cloud-camp

### GoCloudCamp test assignment

Сервер и клиентская библиотека для динамического управления конфигурацией приложений.

Сервер реализован на языке GoLang и использует в качестве хранилища базу данных MongoDB. Конфигурации хранятся в

### Варианты запроса GET:

Получить последнюю версию конфига

```
http://host:port/config?service=name
```

Получить определенный номер версии конфига

```
http://host:port/config?service=name&version=number
```

Варианты кодов ответа сервера:

- 200 –Ок. Запрос выполнен успешно

В теле ответа будет содержаться JSON в формате

```json
{
   "key1": "value1",
   "key2": "value2",
   . . .
}
```

- 400 – Неправильный формат запроса
- 404 – Конфиг не найден
- 500 – Внутренняя ошибка сервера

```json
{
	"error": "текст сообщения об ошибке"
}
```

Дополнительные библиотеки, использованные в проекте:

- [github.com/ilyakaznacheev/cleanenv](github.com/ilyakaznacheev/cleanenv)
- [github.com/julienschmidt/httprouter](github.com/julienschmidt/httprouter)
- [go.mongodb.org/mongo-driver](go.mongodb.org/mongo-driver)
- [go.uber.org/zap](go.uber.org/zap)
