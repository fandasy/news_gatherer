Приветствие
---

Привет, читатель. Это мой первый проект, связанный с телеграм-ботом, предполагаться что этот проект будет использоваться как шаблон для чего то более полезного : )


Функционал
---

- Проводит валидацию всех кинутых в чат ссылок, и сохраняет их если они её прошли
- Вывод всех сохранённых ссылок на новостные медиа
- Удаление ссылок
- Вывод 10 (искусственное ограничение) новостных блоков с каждой новостной ленты
- Вывод с определённой страницы или платформы (VK, RSS)
- Пересказ новостей с RSS ленты
- Ограниченитель запросов

В данной версии поддерживаются только VK и RSS новостные ленты

Запуск
---

Для запуска tg bot можно использовать docker-compose файл, с заранее настроенным config.yaml

Требования: Docker, tg token, vk token, yaGpt token

- tg token: [@BotFather](https://telegram.me/BotFather)

- vk token: [dev VK](https://dev.vk.com/ru), создаёте приложение и получаете Сервисный ключ доступа

- yaGpt token [300.ya.ru](https://300.ya.ru/?nr=1#), в нижней части нажимаете на API, после кнопку Получить токен

Перейдите в директорию проекта и в Dockerfile впишите в поля свои токены:

Example
```
TG_TOKEN=0123456789:AAA0AA0AAA0Aa0AaAaaa00aA0aaaAAaaaA0 
VK_TOKEN=0a00000a0a00000a0a00000a000a00a00a00a000a00000a000a0000aa00000000000aa0  
YA_GPT_TOKEN=y0_AAAAAAA_AAaAAAaA0aAAAAAAAAaaAAAaaAAaAa00aaAAaAaAaaaaAAa0aA
```

После пропишите в консоли

```
docker-compose up --build
```

YAML
---

```
env: "local"

clients:
 - tgBotHost:    "api.telegram.org"
 - vkApiHost:    "api.vk.com"
 - vkApiVersion: "5.131"
 - yaGptHost:    "300.ya.ru"

PSQLConnection:  "user=username dbname=dbname password=password host=localhost port=5432 sslmode=disable"

batchSize:     100  |  лимит обновлений Telegram бота, от 1 до 100, по умолчанию 100

updateTimeout: 50ms |  timeout перед полчением новых данных с tg

reqLimit:
 maxNumberReq: 5    |  кол-во запросов в timeSlice * time.Second
 timeSlice:    2s   |  Промежуток времени (в секундах)
 banTime:      60s  |  Время бана (в секундах)
```

В зависимости от env запускаются типы логирования:

```
- local - text уровень Debug вывод в консоль
- dev   - json уровень Debug вывод в файл
- prod  - json уровень Info  вывод в файл
```
