package logging

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

var DebugEnable bool = false

// Print the detailed content
func Debug(format string, a ...interface{}) {
	format = fmt.Sprintf("[%v] %s\n", aurora.Cyan("DBG"), format)
	if DebugEnable {
		fmt.Printf(format, a...)
	}
}

func Info(format string, a ...interface{}) {
	format = fmt.Sprintf("[%v] %s\n", aurora.Green("INF"), format)
	fmt.Printf(format, a...)
}

func Warn(format string, a ...interface{}) {
	format = fmt.Sprintf("[%v] %s\n", aurora.Yellow("WRN"), format)
	fmt.Printf(format, a...)

}

func Error(format string, a ...interface{}) {
	format = fmt.Sprintf("[%v] %s\n", aurora.Red("ERR"), format)
	fmt.Printf(format, a...)

}

func Critical(format string, a ...interface{}) {
	format = fmt.Sprintf("[%v] %s\n", aurora.BrightRed("CRITICAL"), format)
	fmt.Printf(format, a...)
	os.Exit(1)
}
