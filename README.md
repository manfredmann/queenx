# queenx - Утилита для сборки проекта под QNX4 на target системе.

#### В чём идея:
Я гоняю QNX4 в виртуальной машине, т.к. программировать прямо под ней - боль и страдание. Раньше я выходил из положения так: подключал файловую систему по sshfs, а собирал из ssh сессии. Казалось бы, всё хорошо, но я использую git. Файловая система QNX4 имеет ограничение на длину имён файлов, в которое не вписываются некоторые файлы из каталога .git, да и держать исходники в QNX4 - такое себе, иногда это не VM а реальная железка, и неудобно таскать всё это туда-сюда. Потому была наспех написана данная тулза (ранее был написан bash скрипт, но честно, меня тошнит от шелл скриптинга). Что она делает:

- Создаёт дерево каталогов, необходимое для проекта
- Копирует каталоги и файлы из текущего каталога в соответствии с описанием проекта в файле конфигурации
- Выполняет на хосте с QNX4 make clean && make (позже добавлю команду выполнения сборки в файл конфигурации проекта) 

Можно было бы написать хитровывернутый Makefile, но я из другой категории извращенцев.

// Заранее прошу прощения за мой кривой английский, можете мне в этом помочь (всем как обычно пофиг).

#### Требования
Установленный и запущенный OpenSSH на хосте с qnx4
Любой GNU/Linux дистрибутив, где есть ssh клиент

#### Как пользоваться
Положить в каталог с  проектом файл конфигурации queenx.yml, пример:
```
local:
  project_name: "PtSoko" //Имя проекта
  project_dirs:  // Каталоги, которые буду скопированы
    - src
    - obj
    - bin
  project_files: // Файлы, которые будут скопированы
    - Makefile
remote:
  host: "qnx4vm" // Адрес машины
  projects_path: "/root/projects" //Путь, по которому будет создан каталог с именем проекта
build:
  cmd_pre: "make clean" // Выполняется перед сборкой
  cmd_build: "make" // Сборка
  cmd_post: "" // Выполняется после сборки
```

Далее проинициализируем проект (На target системе будет создано дерево каталогов, в соответствии с конфигурацией).

```
~/q/p/PtSoko > queenx init
 -- Checking project dirs on remote host...
 -- [/root/projects/PtSoko]: Creating... OK
 -- [/root/projects/PtSoko/src]: Creating... OK
 -- [/root/projects/PtSoko/inc]: Creating... OK
 -- [/root/projects/PtSoko/obj]: Creating... OK
 -- [/root/projects/PtSoko/bin]: Creating... OK
```

Ну и соберём проект
```
~/q/p/PtSoko > queenx build
 -- Transferring files to remote host...
 -- [./src --> /root/projects/PtSoko/src]: 
box.cpp                                                                                                                                                                                                   100% 2170   670.0KB/s   00:00    
main.cpp                                                                                                                                                                                                  100% 1029   942.7KB/s   00:00    
player.cpp                                                                                                                                                                                                100% 1903     1.3MB/s   00:00    
help.cpp                                                                                                                                                                                                  100%  257   183.2KB/s   00:00    
game.cpp                                                                                                                                                                                                  100%   21KB   5.4MB/s   00:00    
brick.cpp                                                                                                                                                                                                 100% 1804     1.6MB/s   00:00    
object.cpp                                                                                                                                                                                                100%  809     1.0MB/s   00:00    
box_place.cpp                                                                                                                                                                                             100% 1725     2.8MB/s   00:00    
 -- [./inc --> /root/projects/PtSoko/inc]: 
game.h                                                                                                                                                                                                    100% 3582   771.7KB/s   00:00    
object.h                                                                                                                                                                                                  100% 1550     1.2MB/s   00:00    
brick.h                                                                                                                                                                                                   100% 1222     1.6MB/s   00:00    
player.h                                                                                                                                                                                                  100% 1248     1.8MB/s   00:00    
help.h                                                                                                                                                                                                    100%  164   249.1KB/s   00:00    
box.h                                                                                                                                                                                                     100% 1281     2.2MB/s   00:00    
box_place.h                                                                                                                                                                                               100% 1181   744.4KB/s   00:00    
 -- [./obj --> /root/projects/PtSoko/obj]: 
.placeholder                                                                                                                                                                                              100%    0     0.0KB/s   00:00    
 -- [./bin --> /root/projects/PtSoko/bin]: 
.placeholder                                                                                                                                                                                              100%    0     0.0KB/s   00:00    
 -- [./Makefile --> /root/projects/PtSoko/Makefile]: 
Makefile                                                                                                                                                                                                  100%  699   133.6KB/s   00:00    
 -- Prebuild...
rm -f ./obj/*.o ./bin/PtSoko ./bin/*.map *.err 
 -- Build...
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/help.o src/help.cpp
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/box.o src/box.cpp
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/box_place.o src/box_place.cpp
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/object.o src/object.cpp
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/brick.o src/brick.cpp
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/player.o src/player.cpp
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/game.o src/game.cpp
cc -Oentx -ms -s -w1 -5r, -WC,-xss -I./inc -c -o obj/main.o src/main.cpp
cc -M -N 64k -lphoton -l/qnx4/phtk/lib/phrender_s.lib -l/qnx4/phtk/lib/phexlib3r.lib -o bin/PtSoko obj/help.o obj/box.o obj/box_place.o obj/object.o obj/brick.o obj/player.o obj/game.o obj/main.o
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
