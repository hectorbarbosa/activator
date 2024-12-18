### About
Sample REST API Service with email activation 

### Установка
1. Клонировать репозиторий
2. Создайте базу данных `activator`.
3. Проверьте настройки PosgreSQL в файле `.env`. Сервер запускается на порту 8080 по умолчанию.
4. Создайте таблицы при помощи скрипта.
5. На локальном хосте дожен быть установлен smtp-клиент (без аутентификации), например Postfix.

В терминале линукс, тестировалось на Ubuntu 22:
```shell
git clone https://github.com/hectorbarbosa/activator.git
make createdb
# build
make
# Start server (port 8080 by default)
make run
# Create tables
make migrateup
```