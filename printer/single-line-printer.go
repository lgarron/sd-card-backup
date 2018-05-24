package printer

import (
	"bytes"
	"fmt"
)

func Clear() {
	var buffer bytes.Buffer
	// Writing out `\r` *should* work, but seems to mix lines that were meant to be
	// separate. This works better, and doesn't seem to have a significant impact on
	// overall program performance.
	for i := 0; i < 200; i++ {
		buffer.WriteString("\b \b")
	}
	fmt.Printf(buffer.String())
}

func Printf(format string, args ...interface{}) {
	Clear()
	fmt.Printf(fmt.Sprintf(format, args...))
}
