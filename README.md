# graduate-2

## Общее описание

1. Протокол взаимодействия - REST (или что-то на него похожее)
1. Хранение данных на стороне сервера: сейчас sqlite, в принципе, несложно
   перейти на другой вариант SQL БД.
1. Аутентификация - basic + https
1. сохраняемые данные - на настоящий момент пароли (`Accounts`),
   текстовые записи (`Notes`) и данные платёжных карт (`Cards`)
   TODO: добавить бинарные данные
1. Шифрование данных: на стороне клиента с помощью мастер-ключа
   алгоритмом AES. Данные шифруются перед отправкой на сервер и кэшированием.
   Шифруются поля `Password` и `Meta` для `Account`, и `Text` и `Meta` для `Note`,
   `Number`, `CVC` и `Meta` - для `Card`.

## Организация кода
1. Внутренние модули:
   * `internal/store/`: код, работающий с БД
   * `internal/server/`: код, работающий в http-сервере
   * `internal/client/`: код, работающий в http-клиенте
   * `internal/client/tmpl/`: шаблоны для автогенерации кода клиента
     Часть файлов генерируются при помощи `go generate`
   * `cmd/server/`: код для запуска сервера
   * `cmd/client/`: код для запуска клиента
   * `cmd/server/internal/`, cmd/client/internal` - внутренние модули команд
     сервера и клиента.  Часть файлов генерируются при помощи `go generate`.

## Кэширование на стороне клиента
1. Сохранение и обновление: при успешной попытке сохранения или обновления
   на сервере, обновляется также запись в локальном кэше.
1. Чтение: делается попытка чтения с сервера. Если не удаётся связаться с сервером -
   делается попытка чтения из кэша. При удачной попытке чтения с сервера - кэш
   обновляется. Операция `list` локального кэша не обновляет, но при недоступности
   сервера выводит содержимое локального кэша.
1. Удаление: при успешном удалении с сервера запись из локального кэша удаляется.

## Ключи командрной строки клиента
Общая схема:
```
go run cmd/client/main.go MODE -a ACTION flags
```
гдеs
* `MODE` - один из `user`, `acc`, `note`, `card`
* `ACTION`
  * для режима `user` один из `register` или `verify`
  * для режимов `acc`, `note` или `card` - один из
    `list`, `store`, `get`, `update` или `delete`
* `flags`:
  * `-h` - получить справку по флагам
  * для режима `acc`:
    ```
    -i int
    	account ID
    -l string
    	account URL
    -m string
    	account metainfo
    -n string
    	account name
    -p string
    	account password
    -u string
    	account user name
    ```
  * для режима `note`:
    ```
    -i int
    	note ID
    -m string
    	note metainfo
    -n string
    	note name
    -t string
    	note text
    ```
  * для режима `card`:
    ```
    -c string
    	card CVC code
    -ch string
    	card holder
    -em int
    	card expiry month
    -ey int
    	card expiry year
    -i int
    	card ID
    -m string
    	card metainfo
    -n string
    	card name
    -num string
    	card number
    ```


## Использование

1. Сделать конфигурационный файл сервера `server.cfg` (имя можно задать
   в переменной окружения `SERVER_CFG`):
   ```
   {
     "store_file": "server_storage.db",
     "listen_port": 8443,
     "server_key": "keys/server.key",
     "server_crt": "keys/server.crt"
   }
   ```
1. Скопировать ключ и сертификат сервера в соответствующие файлы.
   Можно изготовить самоподписанный сертификат через команду `make key`.
1. Записать секретную фразу в файл `secret_phrase.txt` или другой,
   назначить ему режим доступа `0400` или `0600`.
1. Сделать конфигурационный файл клиента `gosecret.cfg` (имя можно задать
   в переменной окружения `GOSECRET_CFG`). Если для сервера применяется
   самоподписанный сертификат, нужно выставить `https_insecure` в `true`.
   Режим доступа конфигурационного файла клиента должен быть `0400` или `0600`.
   ```
   {
     "user_name": "user1",
     "password": "pass",
     "full_name": "Full Name",
     "server_address": "https://localhost:8443"
     "cache_file": "cache_store.db",
     "https_insecure": true,
     "key_phrase_file": "secret_phrase.txt"
   }
   ```
1. Запустить сервер
   ```
   go run cmd/server/main.go
   ```
1. Зарегистрироваться на сервере:
   ```
   $ go run cmd/client/main.go user -a register
   2022/04/30 09:17:26 user is registered with id 1
   ```
1. Проверить регистрацию:
   ```
   $ go run cmd/client/main.go user -a verify
   2022/04/30 09:18:11 user is verified
   ```

Далее показаны действия по работе с аккаунтами (логин-пароль-url-мета).
Для других типов данных (платёжные карты и текстовые записи)
действия аналогичны.

Для бинарных данных действия отличаются (TODO: добавить)

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

## Разные TODO:
* добавить проверки параметров командной строки на валидность и лимиты
* В режиме обновления записи дать возможность указывать только обновляемые поля,
   а те, что не указаны - не изменять
* Добавить возможность смены пароля пользователя
* Вынести настройку тайм-аута клиента http в конфигурационный файл
* Добавить режим выбора записей по паре "тип-имя" в дополнение
   к режиму выбора записей по ID
* добавить установку прав доступа 0600 к файлу кэша; однако, секретные
  записи в кэше зашифрованы, так что, возможно, это не требуется
* добавить ключ "clear cache"
* Добавить кеширование всех пользовательских данных по запросу единовременно
