
# queenx - Утилита для сборки проектов под QNX4 на target системе.

## В чём идея
Я гоняю QNX4 в виртуальной машине, т.к. программировать прямо под ней - боль и страдание. Раньше я выходил из положения так: подключал файловую систему по sshfs, а собирал из ssh сессии. Казалось бы, всё хорошо, но я использую git. Файловая система QNX4 имеет ограничение на длину имён файлов, в которое не вписываются некоторые файлы из каталога .git, да и держать исходники в QNX4 - такое себе, иногда это не VM а реальная железка, и неудобно таскать всё это туда-сюда. Потому была наспех написана данная тулза (ранее был написан bash скрипт, но честно, меня тошнит от шелл скриптинга). Что она делает:

- Создаёт дерево каталогов, необходимое для проекта
- Копирует каталоги и файлы из текущего каталога в соответствии с описанием проекта в файле конфигурации
- Выполняет на хосте с QNX4 сборку проекта

// Можно было бы написать хитровывернутый Makefile, но я из другой категории извращенцев.

// Заранее прошу прощения за мой кривой английский, можете мне в этом помочь (всем как обычно пофиг).

## Требования
localhost:
* Клиент OpenSSH
* rsync

qnx4:
* OpenSSH

## Как пользоваться
Положить в каталог с  проектом файл конфигурации queenx.yml, пример:
```yml
local:
  # Имя проекта
  project_name: "PtSoko"
  # Каталоги, которые будут скопированы
  project_dirs:
    - src
    - inc
    - obj
    - bin
  # Файлы, которые будут скопированы
  project_files:
    - Makefile
remote:
  # Адрес машины
  host: "qnx4vm"
  # Путь, по которому будет создан каталог с именем проекта
  projects_path: "/root/projects"
build:
  # Выполняется перед сборкой
  cmd_pre: ""
  # Сборка
  cmd_build: "make"
  # Выполняется после сборки
  cmd_post: ""
  # Очистка
  cmd_clean: "make clean"  
```
Адрес хоста можно передать в параметрах при помощи ключа -h

Далее проинициализируем проект (на target системе будет создано дерево каталогов, в соответствии с конфигурацией).

```
~/q/p/PtSoko > queenx init
 -- Checking the directory structure on remote host...
 -- [/root/projects/PtSoko]: Creating... OK
 -- [/root/projects/PtSoko/src]: Creating... OK
 -- [/root/projects/PtSoko/inc]: Creating... OK
 -- [/root/projects/PtSoko/obj]: Creating... OK
 -- [/root/projects/PtSoko/bin]: Creating... OK
```

Ну и соберём проект
```
~/q/p/PtSoko > queenx build
 -- Checking the directory structure on remote host...
 -- [/root/projects/PtSoko]: OK
 -- [/root/projects/PtSoko/src]: OK
 -- [/root/projects/PtSoko/inc]: OK
 -- [/root/projects/PtSoko/obj]: OK
 -- [/root/projects/PtSoko/bin]: OK
 -- Transferring files to remote host...
 -- [./src --> /root/projects/PtSoko/src]: 
sending incremental file list
src/box.cpp
src/box_place.cpp
src/brick.cpp
src/game.cpp
src/help.cpp
src/main.cpp
src/object.cpp
src/player.cpp

sent 31,278 bytes  received 169 bytes  62,894.00 bytes/sec
total size is 30,695  speedup is 0.98
 -- [./inc --> /root/projects/PtSoko/inc]: 
sending incremental file list
inc/box.h
inc/box_place.h
inc/brick.h
inc/game.h
inc/help.h
inc/object.h
inc/player.h

sent 10,723 bytes  received 150 bytes  21,746.00 bytes/sec
total size is 10,228  speedup is 0.94
 -- [./obj --> /root/projects/PtSoko/obj]: 
sending incremental file list
obj/.placeholder

sent 118 bytes  received 36 bytes  102.67 bytes/sec
total size is 0  speedup is 0.00
 -- [./bin --> /root/projects/PtSoko/bin]: 
sending incremental file list
bin/.placeholder
sent 118 bytes  received 36 bytes  308.00 bytes/sec
total size is 0  speedup is 0.00
 -- [./Makefile --> /root/projects/PtSoko/Makefile]: 
sending incremental file list
Makefile

sent 790 bytes  received 35 bytes  1,650.00 bytes/sec
total size is 699  speedup is 0.85
 -- Prebuild...
 -- Nothing to do
 -- Build...
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/help.o src/help.cpp
/usr/watcom/10.6/bin/wpp386 -zq -oentx -w1 -i=./inc -ms -fo=obj/help.o -xss -5r -i=/usr/watcom/10.6/usr/include -i=/usr/include src/help.cpp 
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/box.o src/box.cpp
/usr/watcom/10.6/bin/wpp386 -zq -oentx -w1 -i=./inc -ms -fo=obj/box.o -xss -5r -i=/usr/watcom/10.6/usr/include -i=/usr/include src/box.cpp 
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/box_place.o src/box_place.cpp
/usr/watcom/10.6/bin/wpp386 -zq -oentx -w1 -i=./inc -ms -fo=obj/box_place.o -xss -5r -i=/usr/watcom/10.6/usr/include -i=/usr/include src/box_place.cpp 
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/object.o src/object.cpp
/usr/watcom/10.6/bin/wpp386 -zq -oentx -w1 -i=./inc -ms -fo=obj/object.o -xss -5r -i=/usr/watcom/10.6/usr/include -i=/usr/include src/object.cpp 
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/brick.o src/brick.cpp
/usr/watcom/10.6/bin/wpp386 -zq -oentx -w1 -i=./inc -ms -fo=obj/brick.o -xss -5r -i=/usr/watcom/10.6/usr/include -i=/usr/include src/brick.cpp 
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/player.o src/player.cpp
/usr/watcom/10.6/bin/wpp386 -zq -oentx -w1 -i=./inc -ms -fo=obj/player.o -xss -5r -i=/usr/watcom/10.6/usr/include -i=/usr/include src/player.cpp 
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/game.o src/game.cpp
/usr/watcom/10.6/bin/wpp386 -zq -oentx -w1 -i=./inc -ms -fo=obj/game.o -xss -5r -i=/usr/watcom/10.6/usr/include -i=/usr/include src/game.cpp 
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/main.o src/main.cpp
/usr/watcom/10.6/bin/wpp386 -zq -oentx -w1 -i=./inc -ms -fo=obj/main.o -xss -5r -i=/usr/watcom/10.6/usr/include -i=/usr/include src/main.cpp 
cc -M -N 64k -lphoton -l/qnx4/phtk/lib/phrender_s.lib -l/qnx4/phtk/lib/phexlib3r.lib -o bin/PtSoko obj/help.o obj/box.o obj/box_place.o obj/object.o obj/brick.o obj/player.o obj/game.o obj/main.o
/usr/watcom/10.6/bin/wlink op quiet form qnx flat na bin/PtSoko op static op map=bin/PtSoko.map op priv=3 op c libp /usr/watcom/10.6/usr/lib:/usr/lib:. l /usr/lib/photon3r.lib l /qnx4/phtk/lib/phrender_s.lib l /qnx4/phtk/lib/phexlib3r.lib f obj/help.o f obj/box.o f obj/box_place.o f obj/object.o f obj/brick.o f obj/player.o f obj/game.o f obj/main.o op offset=388k op st=64k  
Warning(1027): file obj/box_place.o(/root/projects/PtSoko/src/box_place.cpp): redefinition of _PxImageFunc ignored
Warning(1027): file obj/object.o(/root/projects/PtSoko/src/object.cpp): redefinition of _PxImageFunc ignored
Warning(1027): file obj/brick.o(/root/projects/PtSoko/src/brick.cpp): redefinition of _PxImageFunc ignored
Warning(1027): file obj/player.o(/root/projects/PtSoko/src/player.cpp): redefinition of _PxImageFunc ignored
Warning(1027): file obj/game.o(/root/projects/PtSoko/src/game.cpp): redefinition of _PxImageFunc ignored
Warning(1027): file obj/main.o(/root/projects/PtSoko/src/main.cpp): redefinition of _PxImageFunc ignored
 -- Postbuild...
 -- Nothing to do
```

Можно запустить собранный проект через ssh сессию выполнив
```
 > queenx run
```
Все аргументы cli после run будут переданы запускаемому приложению. Бинарник должен лежать в bin/ и иметь имя, соответствующее названию проекта (возможно вынесу это в конфигурацию проекта)

Чтобы запустить проект на определённой node, номер можно передать через параметр -n


## Список команд
#### Инициализация
Создание структуры каталогов проекта на удалённом хосте с qnx. 
Проверяет их наличие, и при необходимости создаёт. Также выполняется вместе с build
```
> queenx init 
```
Ключи:

-r - Переинициализировать структуру каталогов (удаляет при наличии, и создаёт заново)

-h - Адрес хоста

---

#### Сборка
Выполнение последовательно трёх команд, указаных в конфигурационном файле проекта (cmd_pre, cmd_build, cmd_post).
```
> queenx build
```
Ключи:

-r - Переинициализировать структуру каталогов (удаляет при наличии, и создаёт заново)

-h - Адрес хоста

---

#### Очистка
Выполнение команды очистки (cmd_clean)
```
> queenx clean
```
Ключи:

-h - Адрес хоста

---

#### Запуск
Запуск проекта на удалённом хосте с qnx
```
> queenx run
```
Ключи:

-n - Запуск приложения на определённой node

-h - Адрес хоста
