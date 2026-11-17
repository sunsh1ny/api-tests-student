# API-тесты TestBoard

Нужны Go 1.25.1+, golangci-lint 2.x и два аккаунта TestBoard. Зарегистрировать их можно в [Swagger](https://testboard.avito.com/swagger/index.html) через `POST /auth/register` (`email`, `first_name`, `password`). Первый аккаунт создаёт объявления, второй добавляет их в избранное.

Создайте `.env.override` в корне проекта:

```dotenv
TEST_LOGIN=owner@example.com
TEST_PASSWORD=owner_password
TEST_SECOND_LOGIN=reader@example.com
TEST_SECOND_PASSWORD=reader_password
```

Запуск из корня:

```sh
go test -p 1 -count=1 -v ./...
golangci-lint run
```

Объявления заранее создавать не нужно. Для трёх прогонов замените `-count=1` на `-count=3`.
