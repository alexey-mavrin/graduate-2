# graduate-2

## Общее описание

1. Протокол взаимодействия - REST (или что-то на него похожее)
1. Хранение данных на стороне сервера: сейчас sqlite, в принципе, несложно
   перейти на другой вариант SQL БД.
1. Аутентификация - basic + https
1. сохраняемые данные - пароли (`Account`), текстовые записи (`Note`),
   данные платёжных карт (`Card`) и бинарные данные (`Binary`).
1. Шифрование данных: на стороне клиента с помощью мастер-ключа
   алгоритмом AES. Данные шифруются перед отправкой на сервер и кэшированием.
   Шифруются поля `Opaque` и `Meta` структуры `Record`. В поле `Opaque`
   сохраняется содержимое `Account`, `Card`, `Note` или `Binary`.

## Организация кода
1. Внутренние модули:
   * `internal/store/`: код, работающий с БД
   * `internal/server/`: код, работающий в http-сервере
   * `internal/client/`: код, работающий в http-клиенте
   * `cmd/server/`: код для запуска сервера
   * `cmd/client/`: код для запуска клиента
   * `cmd/server/internal/`, cmd/client/internal` - внутренние модули команд
     сервера и клиента.

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
* `MODE` - один из `user`, `cache`, `acc`, `note`, `card` или `bin`
* `ACTION`
  * для режима `user` один из `register` или `verify`
  * для режима `cache` один из `clean` или `sync`
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
  * для режима `bin`:
    ```
    -a string
    	action: list|store|get|update|delete (default "list")
    -f string
    	file name
    -i int
    	binary record ID
    -n string
    	binary record name

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
действия аналогичны. Для бинарных данных примеры приведены ниже -
они отличаются заданием имени файла.

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

   Id: 2
     Type: account
     Name: Another account
   Id: 1
     Type: account
     Name: My account
   ```
1. Обновить аккаунт по ID
   ```
   $ go run cmd/client/main.go acc -a update \
       -i 2 \
       -n "Another account" \
       -u user22 \
       -p passW0RD \
       -l https://example.org:4433 \
       -m "second account, Meta Info"
   record updated
   ```
1. Обновить аккаунт по имени
   ```
   $ go run cmd/client/main.go acc -a update \
       -n "Another account" \
       -u user22 \
       -p passW0RD \
       -l https://example.org:4433 \
       -m "second account, Meta Info"
   record updated
   ```
1. Получить запись аккаунта по ID:
   ```
   $ go run cmd/client/main.go acc -a get -i 2

    Type: account
    Name: Another account
    Meta info: second account, Meta Info
    Data: {"url":"https://example.org:4433","user_name":"user22","password":"passW0RD"}
   ```
1. Получить запись аккаунта по имени:
   ```
   $ go run cmd/client/main.go acc -a get -n "Another account"

    Type: account
    Name: Another account
    Meta info: second account, Meta Info
    Data: {"url":"https://example.org:4433","user_name":"user22","password":"passW0RD"}
   ```
1. Удалить аккаунт по ID:
   ```
   $ go run cmd/client/main.go acc -a delete -i 1
   Record 1 deleted
   ```
1. Удалить аккаунт по имени:
   ```
   $ go run cmd/client/main.go acc -a delete -n "Another account"
   Record 1 deleted
   ```

1. Для бинарных данных использование предполагает задание имени файла
   для исходных данных или для их сохранения на локальной машине.
   * для сохранения содержимого файла `FILE_NAME` на сервере:
     ```
     go run cmd/client/main.go bin -a store -n REC_NAME -f FILE_NAME
     record stored with id 3
     ```
   * для получения с сервера и сохранения в файле `FILE_NAME`:
     ```
     go run cmd/client/main.go bin -a get -i 3 -f FILE_NAME

       Type: binary
       Name: REC_NAME
       File FILE_NAME is written
     ```
1. Очистка локального кэша
   ```
   go run cmd/client/main.go cache -a clean
   2022/05/15 17:39:59 cache is cleaned
   ```
1. Синхронизация кэша с данными сервера
   ```
   go run cmd/client/main.go cache -a sync
   2022/05/15 17:36:57 cache is synchronized
   ```
1. Работа от другого пользователя: создать другой файл конфигурации,
   именить его параметры, указать через переменную окружения:
   ```
   GOSECRET_CFG=gosecret1.cfg go run cmd/client/main.go user -a register
   2022/05/13 22:52:40 user is registered with id 2

   GOSECRET_CFG=gosecret1.cfg go run cmd/client/main.go user -a verify
   2022/05/13 22:58:04 user is verified
   ```
1. В случае утери секретной фразы, зная пароль учётной записи,
   можно получить список записей, но нельзя получить их секретную часть:
   ```
   echo "wrongphrase" > keys/secret_key_phrase_2.txt

   GOSECRET_CFG=gosecret1.cfg go run cmd/client/main.go user -a verify
   2022/05/13 23:01:04 user is verified

   GOSECRET_CFG=gosecret1.cfg go run cmd/client/main.go acc -a list

   Id: 4
     Type: account
     Name: rec1

   GOSECRET_CFG=gosecret1.cfg go run cmd/client/main.go acc -a get -i 4
   2022/05/13 23:01:46 cipher: message authentication failed
   ```

## Разные TODO:
* добавить проверки параметров командной строки на валидность и лимиты
* В режиме обновления записи дать возможность указывать только обновляемые поля,
   а те, что не указаны - не изменять
* Добавить возможность смены пароля пользователя
* Вынести настройку тайм-аута клиента http в конфигурационный файл
* добавить установку прав доступа 0600 к файлу кэша; однако, секретные
  записи в кэше зашифрованы, так что, возможно, это не требуется
* Сделать сообщения об ошибках более user friendly. Например, при попытке создания
  записи с повторяющимся именем и типом, клиент сейчас пишет
  ```
  2022/05/13 22:48:15 storing account: http status 500: Cannot Store Record: UNIQUE constraint failed: records.user_id, records.name, records.type
  ```
  А при попытке получить несуществующую запись
  ```
  2022/05/13 22:54:58 get account: http status 404
  ```
