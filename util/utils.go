package util

import (
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func SanitizeString(str string) string {
	l := len(str)

	if l == 0 {
		return str
	}

	// Remove last \n if ending with one
	if str[l-1] == '\n' {
		str = str[:l-1]
	}

	return str
}

func StrToInt(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func StrToIntOrDefault(logger log.Logger, str string, defaultValue int64) int64 {
	result, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		level.Error(logger).Log("err", err, "str", str)
		return defaultValue
	}

	return result
}

func FirstElement(line string) string {
	return strings.Split(line, " ")[0]
}

func ExtractLeadingInteger(line string, logger log.Logger) int64 {
	el := FirstElement(line)
	v, err := StrToInt(el)
	if err != nil {
		level.Error(logger).Log("err", err, "line", line)
		return -1
	}

	return v
}

func ExtractTrailingValueAfterColon(line string, logger log.Logger) int64 {
	array := strings.Split(line, ":")

	lastValue := array[len(array)-1]
	lastValue = strings.TrimSpace(lastValue)

	v, err := StrToInt(lastValue)
	if err != nil {
		level.Error(logger).Log("err", err, "line", line)
		return -1
	}

	return v
}

func ExtractLastLine(text string) string {
	if text == "" {
		return ""
	}

	lines := strings.Split(text, "\n")

	length := len(lines)
	lastLine := lines[length-1]

	if lastLine == "" {
		lastLine = lines[length-2]
	}

	return lastLine
}

func CountLines(text string) int {
	if text == "" {
		return 0
	}

	count := strings.Count(text, "\n") + 1

	if text[len(text)-1] == '\n' {
		count--
	}

	return count
}

func ToBoolean(c byte) bool {
	switch c {
	case 'y':
		return true
	case 'Y':
		return true
	case 't':
		return true
	case 'T':
		return true
	case '1':
		return true

	case 'n':
		return false
	case 'N':
		return false
	case 'f':
		return false
	case 'F':
		return false
	case '0':
		return false
	default:
		return false
	}
}

func BoolToFloat(b bool) float64 {
	if b {
		return 1
	}

	return 0
}
