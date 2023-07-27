//go:build !solution

package varfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	sprintCache := make([]string, 0, len(args))

	capacity := cacheArgsAndReturnLength(&sprintCache, args)

	var buffer strings.Builder

	stringLength := len(format)
	counter := 0

	buffer.Grow(len(format) + capacity)

	counter = 0

	for i := 0; i < stringLength; i++ {
		if format[i] != '{' {
			buffer.WriteByte(format[i])
			continue
		}

		if format[i] == '{' && i < stringLength-1 && format[i+1] == '}' {
			buffer.WriteString(sprintCache[counter])
			counter++
			i += 1
		} else {
			j := i + 1
			for ; format[j] != '}'; j++ {
			}

			num, err := strconv.Atoi(format[i+1 : j])

			if err != nil {
				panic("Error with string to int conversion")
			}

			buffer.WriteString(sprintCache[num])
			i = j
			counter++
		}
	}

	return buffer.String()
}

func cacheArgsAndReturnLength(sprintCache *[]string, args []interface{}) int {
	capacity := 0
	for _, arg := range args {
		str := fmt.Sprint(arg)
		*sprintCache = append(*sprintCache, str)
		capacity += len(str)
	}

	return capacity
}
