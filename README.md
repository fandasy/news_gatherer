### Приветствие
Привет, читатель. Это мой первый проект, связанный с телеграм-ботом, предполагаться что этот проект будет использоваться как шаблон для чего то более полезного : )

### Функционал:
- Проводит валидацию всех кинутых в чат ссылок, и сохраняет их если они её прошли
- Вывод всех сохранённых ссылок на новостные медиа
- Удаление ссылок
- Вывод 10 (искусственное ограничение) новостных блоков с каждой новостной ленты
- Вывод с определённой страницы или платформы (VK, RSS)
- Пересказ новостей с RSS ленты

В данной версии поддерживаются только VK и RSS новостные ленты

---
### Как запустить мой проект: файл _exe содержит всё необходимое для запуска

- data.json — данные для подключения к базе данных, получения API VK, API Telegram и указания лимита обновлений для Telegram бота
- app.exe - приложение
- run.bat - для удобства (не забудьте поменять <token>)

---
### JSON

- "env": "local" 
- "tgBotHost": "api.telegram.org",
- "vkApiHost": "api.vk.com",
- "vkApiVersion": "5.131",
- "yaGptHost": "300.ya.ru",
- "PSQLconnection": "user=username dbname=db password=pass host=localhost port=5432 sslmode=disable",
- "batchSize": 100
- "maxNumberReq": 5, // кол-во запросов в timeSlice * time.Second
- "timeSlice": 1,    // Промежуток времени (в секундах)
- "banTime": 60      // Время бана (в секундах)

batchSize — лимит обновлений Telegram бота, от 1 до 100, по умолчанию 100

Варианты env
- local - text, уровень Debug
- dev   - json, уровень Debug
- prod  - json, уровень Info

---
### Для запуска бота существует два варианта:
- Запустить в консоли с флагами: -config-path, -tg-bot-token, -vk-bot-token, -ya-gpt-token
- Использовать переменные среды: CONFIG_PATH, TG_TOKEN, VK_TOKEN, YA_GPT_TOKEN


```
start name.exe -config-path data.json -tg-bot-token <token> -vk-bot-token <token> -ya-gpt-token <token>

// Example

start v2.0.1.exe -config-path data.json -tg-bot-token 0123456789:AAA0AA0AAA0Aa0AaAaaa00aA0aaaAAaaaA0 -vk-bot-token 0a00000a0a00000a0a00000a000a00a00a00a000a00000a000a0000aa00000000000aa0  -ya-gpt-token y0_AAAAAAA_AAaAAAaA0aAAAAAAAAaaAAAaaAAaAa00aaAAaAaAaaaaAAa0aA
```
