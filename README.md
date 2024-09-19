# Телеграм бот для создания и обработки тикетов клиентов службы поддержки.
## Задачи определены в файле `Taskfile.yml` и могут быть выполнены с помощью команды `task`.
### Предварительные требования
1. Docker установлен на вашем компьютере
2. Установленный инструмент командной строки task (вы можете установить его отсюда [ссылка](https://taskfile.dev/#/installation)
##  Запуск бота с помощью docker compose
1. Сборка образов:
```shell
docker-compose build
```
2. Запуск сервисов:
```shell
docker-compose up -d
```
3. Остановка сервисов:
```shell
docker-compose down
```
## Команды Taskfile
1. Собрать бота:
```shell
task build-bot
```
2. Запустить бота:
```shell
task run-bot
```
#### Очистка
1. Удалить Docker-образы и сеть:
```shell 
task cleanup
```

