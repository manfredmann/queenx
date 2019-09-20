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
	Color_RED     = "\033[1;31m"
	Color_REDL    = "\033[1;91m"
	Color_WHITE   = "\033[1;97m"
	Color_YELLOW  = "\033[1;33m"
	Color_YELLOWL = "\033[1;93m"
)

func Colorf(color string, format string, a ...interface{}) {
	var str = fmt.Sprintf("%s%s\033[0m", color, format)

	fmt.Printf(str, a...)
}

func Colorln(color string, str string) {
	fmt.Printf("%s%s\033[0m\n", color, str)
}

func Printf(format string, a ...interface{}) {
	Colorf(Color_WHITE, format, a...)
}

func Println(str string) {
	Colorln(Color_WHITE, str)
}

func Errorf(format string, a ...interface{}) {
	Colorf(Color_REDL, format, a...)
}

func Errorln(str string) {
	Colorln(Color_REDL, str)
}

func Warningf(format string, a ...interface{}) {
	Colorf(Color_YELLOWL, format, a...)
}

func Warningln(str string) {
	Colorln(Color_YELLOWL, str)
}
