/*
* queenx - CLI tool for building projects for the QNX4 on target machine
* Copyright (C) 2019  Roman Serov <roman@serov.co>
*
* This file is part of queenx.
*
* queenx is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* queenx is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with queenx. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
)

const (
	Color_GREEN   = "\033[1;32m"
	Color_GREENL  = "\033[1;92m"
	Color_RED     = "\033[1;31m"
	Color_REDL    = "\033[1;91m"
	Color_WHITE   = "\033[1;97m"
	Color_YELLOW  = "\033[1;33m"
	Color_YELLOWL = "\033[1;93m"
)

type Logger struct {
	new_line     bool
	prefix       string
	prefix_color string
}

func LoggerInit(prefix string, prefix_color string) *Logger {
	var lg Logger

	lg.prefix = prefix
	lg.new_line = true
	lg.prefix_color = prefix_color

	return &lg
}

func (l *Logger) Colorf(color string, format string, a ...interface{}) {
	if l.new_line {
		fmt.Printf("%s%s\033[0m", l.prefix_color, l.prefix)
	}

	if format[len(format)-1] == '\n' {
		l.new_line = true
	} else {
		l.new_line = false
	}

	var str = fmt.Sprintf("%s%s\033[0m", color, format)

	fmt.Printf(str, a...)
}

func (l *Logger) Colorln(color string, str string) {
	if l.new_line {
		fmt.Printf("%s%s\033[0m", l.prefix_color, l.prefix)
	}

	fmt.Printf("%s%s\033[0m\n", color, str)

	l.new_line = true

}

func (l *Logger) Printf(format string, a ...interface{}) {
	l.Colorf(Color_WHITE, format, a...)
}

func (l *Logger) Println(str string) {
	l.Colorln(Color_WHITE, str)
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Colorf(Color_REDL, format, a...)
}

func (l *Logger) Errorln(str string) {
	l.Colorln(Color_REDL, str)
}

func (l *Logger) Warningf(format string, a ...interface{}) {
	l.Colorf(Color_YELLOWL, format, a...)
}

func (l *Logger) Warningln(str string) {
	l.Colorln(Color_YELLOWL, str)
}
