package log

import (
	"fmt"
	"strconv"
)

var DebugMode = false

func Debug(a ...interface{}) {
	if DebugMode {
		fmt.Println(Grey(fmt.Sprint(a...)))
	}
}

const esc = "\033"
const reset = esc + "[0m"
const grey = 2
const yellow = 33
const blue = 34

func printWithCode(s string, c int) string {
	return fmt.Sprintf(esc+"["+strconv.Itoa(c)+"m%s"+reset, s)
}

// Grey formats colored terminal text.
func Grey(s string) string {
	return printWithCode(s, grey)
}
