

# queenx - Утилита для сборки проектов под QNX4 на target системе.
- [В чём идея](#в-чём-идея)
- [Требования](#требования)
- [Как пользоваться](#как-пользоваться)
- [Шаблоны](#шаблоны)
- [Список команд](#список-команд)

## В чём идея
Я гоняю QNX4 в виртуальной машине, т.к. программировать прямо под ней - боль и страдание. Раньше я выходил из положения так: подключал файловую систему по sshfs, а собирал из ssh сессии. Казалось бы, всё хорошо, но я использую git. Файловая система QNX4 имеет ограничение на длину имён файлов, в которое не вписываются некоторые файлы из каталога .git, да и держать исходники в QNX4 - такое себе, иногда это не VM а реальная железка, и неудобно таскать всё это туда-сюда. Потому была наспех написана данная тулза (ранее был написан bash скрипт, но честно, меня тошнит от шелл скриптинга). Что она делает:

- Создаёт дерево каталогов, необходимое для проекта
- Копирует каталоги и файлы из текущего каталога в соответствии с описанием проекта в файле конфигурации
- Выполняет на хосте с QNX4 сборку проекта
- Позволяет создавать структуру каталогов для нового проекта из шаблона

// Можно было бы написать хитровывернутый Makefile, но я из другой категории извращенцев.

// Заранее прошу прощения за мой кривой английский, можете мне в этом помочь (всем как обычно пофиг).

## Требования
localhost:
* GNU/Linux
* Клиент OpenSSH
* rsync

qnx4:
* Сервер OpenSSH (также нужен prngd, т.к. 4ка не имеет генератора псевдослучайных числе)
* rsync

Взять можно тут: https://stuff.pentium02.org/qnx4/ (prngd и rsync были добыты тут: http://forum.kpda.ru/)

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
# Параметры для команды run
run:
  # Каталог с бинарником
  bin_path: "bin"
  # Имя бинарника
  bin_name: "PtSoko"
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
          2,170 100%    0.00kB/s    0:00:00 (xfr#1, to-chk=7/9)
src/box_place.cpp
          1,725 100%    1.65MB/s    0:00:00 (xfr#2, to-chk=6/9)
src/brick.cpp
          1,804 100%    1.72MB/s    0:00:00 (xfr#3, to-chk=5/9)
src/game.cpp
         20,998 100%   20.03MB/s    0:00:00 (xfr#4, to-chk=4/9)
src/help.cpp
            257 100%  250.98kB/s    0:00:00 (xfr#5, to-chk=3/9)
src/main.cpp
          1,029 100% 1004.88kB/s    0:00:00 (xfr#6, to-chk=2/9)
src/object.cpp
            809 100%  790.04kB/s    0:00:00 (xfr#7, to-chk=1/9)
src/player.cpp
          1,903 100%    1.81MB/s    0:00:00 (xfr#8, to-chk=0/9)
 -- [./inc --> /root/projects/PtSoko/inc]: 
sending incremental file list
inc/box.h
          1,281 100%    0.00kB/s    0:00:00 (xfr#1, to-chk=6/8)
inc/box_place.h
          1,181 100%    1.13MB/s    0:00:00 (xfr#2, to-chk=5/8)
inc/brick.h
          1,222 100%    1.17MB/s    0:00:00 (xfr#3, to-chk=4/8)
inc/game.h
          3,582 100%    3.42MB/s    0:00:00 (xfr#4, to-chk=3/8)
inc/help.h
            164 100%  160.16kB/s    0:00:00 (xfr#5, to-chk=2/8)
inc/object.h
          1,550 100%    1.48MB/s    0:00:00 (xfr#6, to-chk=1/8)
inc/player.h
          1,248 100%    1.19MB/s    0:00:00 (xfr#7, to-chk=0/8)
 -- [./obj --> /root/projects/PtSoko/obj]: 
sending incremental file list
obj/.placeholder
              0 100%    0.00kB/s    0:00:00 (xfr#1, to-chk=0/2)
 -- [./bin --> /root/projects/PtSoko/bin]: 
sending incremental file list
bin/.placeholder
              0 100%    0.00kB/s    0:00:00 (xfr#1, to-chk=0/2)
 -- [./Makefile --> /root/projects/PtSoko/Makefile]: 
sending incremental file list
Makefile
            699 100%    0.00kB/s    0:00:00 (xfr#1, to-chk=0/1)
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
Все аргументы cli после run будут переданы запускаемому приложению. Каталог с бинарником относительно корня проекта, и имя бинарника берутся из файла конфигураци проекта. Если не заданы, каталог по умолчанию - bin, имя бинарника = имя проекта

Чтобы запустить проект на определённой node, номер можно передать через параметр -n

## Шаблоны
Утилита позволяет создавать новые проекты из шаблонов, при помощи команды new.

Шаблон представляет собой обычный tar архив, который также может быть упакован в gzip. Путь, по которому располагаются шаблоны: ~/.config/queenx/templates.

Чтобы создать новый проект из шаблона, нужно выполнить:
```
> queenx new <имя шаблона> <имя проекта>
```
Имя шаблона представляет собой имя tar архива без расширения

Пример:
```
> queenx new main test
 -- Trying to open "/home/manfredmann/.config/queenx/templates/main.tar"
 -- [test/bin]: OK
 -- [test/bin/.placeholder]: OK
 -- [test/inc]: OK
 -- [test/inc/.placeholder]: OK
 -- [test/Makefile]: OK
 -- [test/obj]: OK
 -- [test/obj/.placeholder]: OK
 -- [test/queenx.yml]: OK
 -- [test/src]: OK
 -- [test/src/.placeholder]: OK
 -- OK
> ls test
итого 32K
drwxr-xr-x 6 manfredmann manfredmann 4,0K сен 19 02:03 .
drwxr-xr-x 3 manfredmann manfredmann 4,0K сен 19 02:03 ..
drwxr-xr-x 2 manfredmann manfredmann 4,0K сен 19 02:03 bin
drwxr-xr-x 2 manfredmann manfredmann 4,0K сен 19 02:03 inc
drwxr-xr-x 2 manfredmann manfredmann 4,0K сен 19 02:03 obj
drwxr-xr-x 2 manfredmann manfredmann 4,0K сен 19 02:03 src
-rw-r--r-- 1 manfredmann manfredmann  552 сен 19 02:03 Makefile
-rw-r--r-- 1 manfredmann manfredmann  261 сен 19 02:03 queenx.yml
```

## Список команд
Ключи следует указывать перед именем команды
#### Инициализация
Создание структуры каталогов проекта на удалённом хосте с qnx. 
Проверяет их наличие, и при необходимости создаёт. Также выполняется вместе с build
```
> queenx init 
```
Ключи:

* -r - Переинициализировать структуру каталогов (удаляет при наличии, и создаёт заново)
* -h - Адрес хоста

---

#### Сборка
Выполнение последовательно трёх команд, указаных в конфигурационном файле проекта (cmd_pre, cmd_build, cmd_post).
```
> queenx build
```
Ключи:

* -r - Переинициализировать структуру каталогов (удаляет при наличии, и создаёт заново)
* -h - Адрес хоста

---

#### Очистка
Выполнение команды очистки (cmd_clean)
```
> queenx clean
```
Ключи:
* -h - Адрес хоста

---

#### Запуск
Запуск проекта на удалённом хосте с qnx. Все аргументы, указанные после run будут переданы запускаемому приложению
```
> queenx run <аргументы>
```
Ключи:
* -n - Запуск приложения на определённой node
* -h - Адрес хоста

---

#### Новый проект
Создание нового проекта из шаблона
```
> queenx new <имя шаблона> <имя проекта>
```
Каталог с шаблонами: ~/.config/queenx/templates
