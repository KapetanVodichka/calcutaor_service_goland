# Калькулятор на Go

Проект представляет собой веб-сервис, который умеет вычислять арифметические выражения, переданные по HTTP (метод POST).

## Функциональность

1. Принимает JSON с полем `"expression"`.
2. Вычисляет результат (поддерживаются операции `+`, `-`, `*`, `/` и скобки `(` и `)`).
3. Возвращает JSON с полем `"result"` и статусом 200 при успешном вычислении.
4. В случае ошибки в самом выражении (недопустимые символы, деление на ноль, несбалансированные скобки) возвращает JSON с полем `"error"` и статусом 422.
5. При любой другой непредвиденной ошибке возвращается JSON с полем `"error"` и статусом 500.

Сервис принимает POST-запрос на эндпоинт `/api/v1/calculate` с телом вида:
```json
{
  "expression": "2+2*2"
}
```
и возвращает результат вычисления в формате JSON.

Формат успешного ответа (код 200 OK):
```json
{
  "result": "6"
}
```

## Запуск сервиса
Клонируйте репозиторий:
```bash
git clone https://github.com/KapetanVodichka/calc_service_goland.git
```
Перейдите в папку с проектом:
```bash
cd calcutaor_service_goland
```

Скачайте необходимые зависимости:
```bash
go mod tidy
```

Запустите сервис:
```bash
go run ./cmd/calc_service/...
```
По умолчанию сервер стартует на порту 8080. Эндпоинт доступен по адресу:
```bash
http://localhost:8080/api/v1/calculate
```

## Возможные коды ошибок:
422 (Unprocessable Entity)
Если в выражении содержатся недопустимые символы, есть деление на ноль или скобки не сбалансированы.
Пример:

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+(3*2"
}'
```
Ответ:

```json
{
  "error": "Expression is not valid"
}
```

500 (Internal Server Error)
Любая иная непредвиденная ошибка внутри приложения.
Пример ответа:

```json
{
  "error": "Internal server error"
}
```

## Примеры запросов
Успешный запрос (200 OK):
```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
Пример ответа:

```json
{
  "result": "6"
}
```
Ошибка 422 (Unprocessable Entity):
```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+(3*2"
}'
```
Пример ответа:

```json
{
  "error": "Expression is not valid"
}
```
Ошибка 500 (Internal Server Error).
Данный код возвращается в случае непредвиденных сбоев внутри приложения:

```json
{
  "error": "Internal server error"
}
