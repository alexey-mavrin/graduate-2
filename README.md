# graduate-2

## Общее описание

1. Протокол взаимодействия - REST (или что-то на него похожее)
1. Хранение данных на стороне сервера: сейчас sqlite, в принципе, несложно
   перейти на другой вариант SQL БД.
1. Аутентификация - basic + https (https пока не реализован, TODO: добавить https)
1. сохраняемые данные - на настоящий момент только пароли (`Accounts`)
   TODO: добавить типы
1. Разрешение конфликтов: TODO, планирую добавить поле `version`, увеличивать
   его при обновлении записей, проверять при обновлении записей и кэша.
1. Шифрование данных в базе: TODO. Будет сделано на стороне клиента с помощью
   мастер-ключа

### Разные TODO:
1. добавить проверки:
   * режим доступа конфигурационного файла клиента (должен быть 0600).
   * проверки параметров командной строки на валидность и лимиты.
1. Улучшить форматирование данных клиентом при выводе пользователю.
1. В режиме обновления записи дать возможность указывать только обновляемые поля,
   а те, что не указаны - не изменять.
1. Добавить возможность смены пароля пользователя.

## Организация кода
1. Внутренние модули:
   * `internal/store/`: код, работающий с БД
   * `internal/server/`: код, работающий в http-сервере
   * `internal/client/`: код, работающий в http-клиенте
   * `cmd/server/`: код для запуска сервера
   * `cmd/client/`: код для запуска клиента

## Кэширование на стороне клиента
1. Запись и обновление: при успешной попытки записи или обновления на сервере,
   обновляется запись в локальном кэше.
1. Чтение: делается попытка чтения с сервера. Если не удаётся связаться с сервером -
   делается попытка чтения из кэша. При удачной попытке чтения с сервера - кэш
   обновляется.
1. Удаление: при успешном удалении с сервера запись из локального кэша удаляется.

TODO:
* добавить установку прав доступа 0600 к файлу кэша.
* добавить ключ "clear cache".

## Использование

1. Сделать конфигурационный файл сервера `server.cfg` (имя можно задать
   в переменной окружения `SERVER_CFG`):
   ```
   {
     "store_file": "server_storage.db",
     "listen_port": 8088
   }
   ```
1. Сделать конфигурационный файл клиента `gosecret.cfg` (имя можно задать
   в переменной окружения `GOSECRET_CFG`):
   ```
   {
     "user_name": "user1",
     "password": "pass",
     "full_name": "Full Name",
     "server_address": "http://localhost:8088"
   }
   ```
1. Запустить сервер
   ```
   go run cmd/server/main.go
   ```
1. Зарегистрироватьс на сервере:
   ```
   $ go run cmd/client/main.go user -a register
   2022/04/30 09:17:26 user is registered with id 1
   ```
1. Проверить регистарцию:
   ```
   $ go run cmd/client/main.go user -a verify
   2022/04/30 09:18:11 user is verified
   ```
1. Сохранить данные аккаунтов:
   ```
   $ go run cmd/client/main.go acc -a store \
       -n "My account" \
       -u user123 \
       -p pass987 \
       -l http://example.com \
       -m "test account, some info"
   accout record stored with id 1

   $ go run cmd/client/main.go acc -a store \
       -n "Another account" \
       -u user22 \
       -p passW0RD \
       -l http://example.org \
       -m "second account, metainfo"
   accout record stored with id 2
   ```
1. Получить список сохранённых аккаунтов:
   ```
   $ go run cmd/client/main.go acc -a list
   map[1:{My account http://example.com user123  test account, some info} 2:{Another account http://example.org user22  second account, metainfo}]
   ```
1. Обновить аккаунт
   ```
   $ go run cmd/client/main.go acc -a update \
       -i 2 \
       -n "Another account" \
       -u user22 \
       -p passW0RD \
       -l https://example.org:4433 \
       -m "second account, Meta Info"
   account updated
   ```
1. Получить запись аккаунта:
   ```
   $ go run cmd/client/main.go acc -a get -i 2
   {Another account https://example.org:4433 user22 passW0RD second account, Meta Info}
   ```
1. Удалить аккаунт:
   ```
   $ go run cmd/client/main.go acc -a delete -i 1
   Account 1 deleted
   ```
