## Локальное окружение

### Зависимости

- docker [https://docs.docker.com/engine/install/](https://docs.docker.com/engine/install/)
- easyjson [https://github.com/mailru/easyjson](https://github.com/mailru/easyjson)

### Запуск

```shell
make up
```

### Остановка

```shell
make down
```

### Подключение к БД  
[docker-compose](./docker/docker-compose.yaml)

## Приложение

### Перед первым запуском
Собрать конфиг:
```shell
make dev-config
```

### Запуск
```shell
 go run . --config build/config.yaml
```

### Локальные ресурсы

- grafana UI [http://localhost:3000/](http://localhost:3000/)
- prometheus UI [http://localhost:9090/](http://localhost:9090/)

## Makefile

Запустить линтер:

```shell
make lint
```

Сгенерировать //go:generate:

```shell
make gen-go
```

Создать файлы миграции БД:

```shell
make migration NAME=example
```

Накатить новые миграции в запущенное локальное окружение:

```shell
make db-migrate-up
```

Откатить последнюю миграцию в запущенном локальном окружении:

```shell
make db-migrate-down
```

Пересобрать локальный вспомогательный докер образ:

```shell
make build FLAGS="--no-cache"
```
