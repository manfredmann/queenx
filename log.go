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
