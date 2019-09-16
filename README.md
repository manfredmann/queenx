# queenx - Утилита для сборки проекта под QNX4 на target системе.

#### В чём идея
Я гоняю QNX4 в виртуальной машине, т.к. программировать прямо под ней - боль и страдание. Раньше я выходил из положения так: подключал файловую систему по sshfs, а собирал из ssh сессии. Казалось бы, всё хорошо, но я использую git. Файловая система QNX4 имеет ограничение на длину имён файлов, в которое не вписываются некоторые файлы из каталога .git, да и держать исходники в QNX4 - такое себе, иногда это не VM а реальная железка, и неудобно таскать всё это туда-сюда. Потому была наспех написана данная тулза (ранее был написан bash скрипт, но честно, меня тошнит от шелл скриптинга). Что она делает:

- Создаёт дерево каталогов, необходимое для проекта
- Копирует каталоги и файлы из текущего каталога в соответствии с описанием проекта в файле конфигурации
- Выполняет на хосте с QNX4 сборку проекта

// Можно было бы написать хитровывернутый Makefile, но я из другой категории извращенцев.

// Заранее прошу прощения за мой кривой английский, можете мне в этом помочь (всем как обычно пофиг).

#### Требования
- Любой GNU/Linux дистрибутив, где есть ssh клиент, и возможность собирать Go приложения
- Установленный и запущенный OpenSSH на хосте с qnx4

#### Как пользоваться
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
  cmd_pre: "make clean"
  # Сборка
  cmd_build: "make"
  # Выполняется после сборки
  cmd_post: ""
```
Адрес хоста можно передать в параметрах при помощи ключа -h

Далее проинициализируем проект (на target системе будет создано дерево каталогов, в соответствии с конфигурацией).

```
~/q/p/PtSoko > queenx init
 -- Checking the project dirs on remote host...
 -- [/root/projects/PtSoko]: Creating... OK
 -- [/root/projects/PtSoko/src]: Creating... OK
 -- [/root/projects/PtSoko/inc]: Creating... OK
 -- [/root/projects/PtSoko/obj]: Creating... OK
 -- [/root/projects/PtSoko/bin]: Creating... OK
```

Ну и соберём проект
```
 -- Transferring files to remote host...
 -- [./src --> /root/projects/PtSoko/src]: 
box.cpp                                                                                                                                                                                                   100% 2170   579.4KB/s   00:00    
main.cpp                                                                                                                                                                                                  100% 1029     1.0MB/s   00:00    
player.cpp                                                                                                                                                                                                100% 1903     2.4MB/s   00:00    
help.cpp                                                                                                                                                                                                  100%  257   413.7KB/s   00:00    
game.cpp                                                                                                                                                                                                  100%   21KB   9.6MB/s   00:00    
brick.cpp                                                                                                                                                                                                 100% 1804     2.1MB/s   00:00    
object.cpp                                                                                                                                                                                                100%  809     1.1MB/s   00:00    
box_place.cpp                                                                                                                                                                                             100% 1725     1.8MB/s   00:00    
 -- [./inc --> /root/projects/PtSoko/inc]: 
game.h                                                                                                                                                                                                    100% 3582   575.0KB/s   00:00    
object.h                                                                                                                                                                                                  100% 1550     1.2MB/s   00:00    
brick.h                                                                                                                                                                                                   100% 1222     1.5MB/s   00:00    
player.h                                                                                                                                                                                                  100% 1248     1.8MB/s   00:00    
help.h                                                                                                                                                                                                    100%  164   224.4KB/s   00:00    
box.h                                                                                                                                                                                                     100% 1281     1.6MB/s   00:00    
box_place.h                                                                                                                                                                                               100% 1181     1.3MB/s   00:00    
 -- [./obj --> /root/projects/PtSoko/obj]: 
.placeholder                                                                                                                                                                                              100%    0     0.0KB/s   00:00    
 -- [./bin --> /root/projects/PtSoko/bin]: 
.placeholder                                                                                                                                                                                              100%    0     0.0KB/s   00:00    
 -- [./Makefile --> /root/projects/PtSoko/Makefile]: 
Makefile                                                                                                                                                                                                  100%  699   252.2KB/s   00:00    
 -- Prebuild...
rm -f ./obj/*.o ./bin/PtSoko ./bin/*.map *.err 
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
```

Можно запустить собранный проект через ssh сессиию выполнив
```
 > queenx run
```
Все аргументы cli после run будут переданы запускаемому приложению. Бинарник должен лежать в bin/ и иметь имя, соответствующее названию проекта (возможно вынесу это в конфигурацию проекта)
