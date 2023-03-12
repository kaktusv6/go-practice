# LOMS

## Локальное развертывание

1. Скопировать пример конфига
```shell
cp config.example.yml config.yml
```

2. Заполнить его недостающими данными

3. Выполнить комманды

```shell
make install-go-deps
make get-go-deps
make vendor-proto
```

## Сборка

```shell
make build
```

## Генерация кода из proto файлов

Если запускаете на Unix системе то выполняйте

```shell
make generate-unix
```

Если запускаете на системе Windows то выполняйте

> make на Windows устанавливать через choco
> 
> `choco install make`

```shell
make generate-win
```
