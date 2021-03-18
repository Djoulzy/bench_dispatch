package clog

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel : Niveau de log
var LogLevel int

// StartLogging : Debut du log
var StartLogging bool

//ServiceCallback usages
var ServiceCallback func(string)

var logToFile = false
var fileDesc *os.File

var fgColors = map[string]string{
	"black":        "0;30",
	"dark_gray":    "1;30",
	"blue":         "0;34",
	"light_blue":   "1;34",
	"green":        "0;32",
	"light_green":  "1;32",
	"cyan":         "0;36",
	"light_cyan":   "1;36",
	"red":          "0;31",
	"light_red":    "1;31",
	"purple":       "0;35",
	"light_purple": "1;35",
	"brown":        "0;33",
	"yellow":       "1;33",
	"light_gray":   "0;37",
	"white":        "1;37",
}

var bgColors = map[string]string{
	"black":        "40",
	"red":          "41",
	"green":        "42",
	"brown":        "43",
	"blue":         "44",
	"magenta":      "45",
	"cyan":         "46",
	"light_gray":   "47",
	"dark_gray":    "100",
	"light_red":    "101",
	"light_green":  "102",
	"yellow":       "103",
	"light_blue":   "104",
	"light_purple": "105",
	"light_cyan":   "106",
	"white":        "107",
}

type errorsColors struct {
	name  string
	level int
	fg    string
	bg    string
}

var _Test = errorsColors{
	name:  "TEST",
	level: 0,
	fg:    "green",
	bg:    "black",
}
var _Debug = errorsColors{
	name:  "DBUG",
	level: 5,
	fg:    "dark_gray",
	bg:    "black",
}
var _Info = errorsColors{
	name:  "INFO",
	level: 4,
	fg:    "light_gray",
	bg:    "black",
}
var _Trace = errorsColors{
	name:  "TRAC",
	level: 4,
	fg:    "white",
	bg:    "black",
}
var _Warn = errorsColors{
	name:  "WARN",
	level: 2,
	fg:    "yellow",
	bg:    "black",
}
var _Error = errorsColors{
	name:  "ERRR",
	level: 1,
	fg:    "light_red",
	bg:    "black",
}
var _Fatal = errorsColors{
	name:  "FATAL",
	level: 1,
	fg:    "white",
	bg:    "red",
}

var _Servc = errorsColors{
	name:  "SERVC",
	level: 1,
	fg:    "light_blue",
	bg:    "black",
}

//GetColoredString add color info
func GetColoredString(str string, fgcolor string, bgcolor string) string {
	coloredString := ""

	if len(fgcolor) != 0 {
		if len(fgColors[fgcolor]) != 0 {
			coloredString = fmt.Sprintf("%s%c[%sm", coloredString, 27, fgColors[fgcolor])
		}
	}

	if len(bgcolor) != 0 {
		if len(bgColors[bgcolor]) != 0 {
			coloredString = fmt.Sprintf("%s%c[%sm", coloredString, 27, bgColors[bgcolor])
		}
	}

	coloredString = fmt.Sprintf("%s%s%c[0m", coloredString, str, 27)
	// return $coloredString;
	return coloredString
}

//CPrintln is a Println colored string
func CPrintln(fgcolor string, bgcolor string, str string) {
	tmp := GetColoredString(str, fgcolor, bgcolor)
	fmt.Println(tmp)
}

//CPrintf is a Printf colored string
func CPrintf(fgcolor string, bgcolor string, format string, vars ...interface{}) {
	tmp := GetColoredString(format, fgcolor, bgcolor)
	fmt.Printf(tmp, vars...)
}

//Output log mechanism
func Output(str string, vars ...interface{}) {
	before := fmt.Sprintf("%s", str)
	tmp := GetColoredString(before, _Info.fg, _Info.bg)
	log.Printf(tmp, vars...)
}

func logOutput(etype errorsColors, pack string, function string, str string, vars ...interface{}) {
	if LogLevel < etype.level || StartLogging == false {
		return
	}
	before := fmt.Sprintf("%s|%s|%s| %s", etype.name, pack, function, str)
	tmp := GetColoredString(before, etype.fg, etype.bg)
	log.Printf(tmp, vars...)
}

// File logger
func File(pack string, function string, str string, vars ...interface{}) {
	if StartLogging == false || logToFile == false {
		return
	}
	before := fmt.Sprintf("%s|%s|%s| %s\n", time.Now().Format("15:04:05"), pack, function, str)
	fileDesc.Write([]byte(fmt.Sprintf(before, vars...)))
}

//Warn message
func Warn(pack string, function string, str string, vars ...interface{}) {
	logOutput(_Warn, pack, function, str, vars...)
}

//Info message
func Info(pack string, function string, str string, vars ...interface{}) {
	logOutput(_Info, pack, function, str, vars...)
}

//Debug message
func Debug(pack string, function string, str string, vars ...interface{}) {
	logOutput(_Debug, pack, function, str, vars...)
}

//Test message
func Test(pack string, function string, str string, vars ...interface{}) {
	logOutput(_Test, pack, function, str, vars...)
}

//Error message
func Error(pack string, function string, str string, vars ...interface{}) {
	logOutput(_Error, pack, function, str, vars...)
}

//Fatal message
func Fatal(pack string, function string, err error) {
	logOutput(_Fatal, pack, function, "%s", err)
	log.Fatal()
}

//Trace output message
func Trace(pack string, function string, str string, vars ...interface{}) {
	logOutput(_Trace, pack, function, str, vars...)
}

//Service callback functions
func Service(pack string, function string, str string, vars ...interface{}) {
	logOutput(_Servc, pack, function, str, vars...)
	if ServiceCallback != nil {
		ServiceCallback(fmt.Sprintf(str, vars...))
	}
}

// EnableFileLog : Log also to file
func EnableFileLog(file string) {
	tmp, err := os.Create(file)
	if err != nil {
		panic(err)
	}

	logToFile = true
	fileDesc = tmp
}
