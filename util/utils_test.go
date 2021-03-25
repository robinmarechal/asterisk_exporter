package util

import (
	"strings"
	"testing"

	"github.com/prometheus/common/promlog"
)

var (
	logCfg = &promlog.Config{}
	logger = promlog.New(logCfg)
)

func TestSanitizeString(t *testing.T) {
	var expected string
	var result string

	expected = ""
	if result = SanitizeString(""); result != expected {
		t.Errorf("Invalid SanitizeString result.\nExpected: '%s'\nActual: '%s'", escape(expected), escape(result))
	}

	expected = "foo"
	if result = SanitizeString("foo"); result != expected {
		t.Errorf("Invalid SanitizeString result.\nExpected: '%s'\nActual: '%s'", escape(expected), escape(result))
	}

	expected = "multi\nline\nstring"
	if result = SanitizeString("multi\nline\nstring"); result != expected {
		t.Errorf("Invalid SanitizeString result.\nExpected: '%s'\nActual: '%s'", escape(expected), escape(result))
	}

	expected = "multi\nline\nstring\nwith tailing new line"
	if result = SanitizeString("multi\nline\nstring\nwith tailing new line\n"); result != expected {
		t.Errorf("Invalid SanitizeString result.\nExpected: '%s'\nActual: '%s'", escape(expected), escape(result))
	}

	expected = "multi\nline\nstring\nwith tailing new line\n\n"
	if result = SanitizeString("multi\nline\nstring\nwith tailing new line\n\n\n"); result != expected {
		t.Errorf("Invalid SanitizeString result.\nExpected: '%s'\nActual: '%s'", escape(expected), escape(result))
	}
}

func escape(str string) string {
	return strings.ReplaceAll(str, "\n", "\\\n")
}

func TestStrToInt_ValidString(t *testing.T) {
	var expected int64
	var result int64
	var err error

	expected = 27
	if result, err = StrToInt("27"); result != expected || err != nil {
		t.Errorf("Invalid StrToInt result.\nErr: '%s'\nExpected: '%d'\nActual: '%d'", err.Error(), expected, result)
	}

	expected = 27
	if result, err = StrToInt("00027"); result != expected || err != nil {
		t.Errorf("Invalid StrToInt result.\nErr: '%s'\nExpected: '%d'\nActual: '%d'", err.Error(), expected, result)
	}
}

func TestStrToInt_InvalidString(t *testing.T) {
	var err error

	if _, err = StrToInt(""); err == nil {
		t.Errorf("Invalid StrToInt result. Call with invalid int should return an error.\nInput: '%s'", "")
	}

	if _, err = StrToInt("aze"); err == nil {
		t.Errorf("Invalid StrToInt result. Call with invalid int should return an error.\nInput: '%s'", "aze")
	}

	if _, err = StrToInt("aze 50"); err == nil {
		t.Errorf("Invalid StrToInt result. Call with invalid int should return an error.\nInput: '%s'", "aze 50")
	}

	if _, err = StrToInt("50 aze"); err == nil {
		t.Errorf("Invalid StrToInt result. Call with invalid int should return an error.\nInput: '%s'", "50 aze")
	}

	if _, err = StrToInt("00027.00"); err == nil {
		t.Errorf("Invalid StrToInt result. Call with invalid int should return an error.\nInput: '%s'", "00027.00")
	}
}

func TestFirstElement(t *testing.T) {
	var expected string
	var result string

	expected = ""
	if result = FirstElement(""); result != expected {
		t.Errorf("Invalid FirstElement result.\nExpected: '%s'\nActual: '%s'", expected, result)
	}

	expected = "aze"
	if result = FirstElement("aze"); result != expected {
		t.Errorf("Invalid FirstElement result.\nExpected: '%s'\nActual: '%s'", expected, result)
	}

	expected = "abc"
	if result = FirstElement("abc def ghi"); result != expected {
		t.Errorf("Invalid FirstElement result.\nExpected: '%s'\nActual: '%s'", expected, result)
	}

	expected = "0"
	if result = FirstElement("0 peers"); result != expected {
		t.Errorf("Invalid FirstElement result.\nExpected: '%s'\nActual: '%s'", expected, result)
	}
}

func TestExtractLeadingInteger(t *testing.T) {
	var expected int64
	var result int64

	expected = -1
	if result = ExtractLeadingInteger("", &logger); result != expected {
		t.Errorf("Invalid ExtractLeadingInteger result. It should return default value -1. Param: ''")
	}

	expected = -1
	if result = ExtractLeadingInteger("abc", &logger); result != expected {
		t.Errorf("Invalid ExtractLeadingInteger result. It should return default value -1. Param: 'abc'")
	}

	expected = -1
	if result = ExtractLeadingInteger("abc def", &logger); result != expected {
		t.Errorf("Invalid ExtractLeadingInteger result. It should return default value -1. Param: 'abc def'")
	}

	expected = -1
	if result = ExtractLeadingInteger("abc 5 def ijk", &logger); result != expected {
		t.Errorf("Invalid ExtractLeadingInteger result. It should return default value -1. Param: '5 def'")
	}

	expected = 5
	if result = ExtractLeadingInteger("5 def ijk", &logger); result != expected {
		t.Errorf("Invalid ExtractLeadingInteger result. It should return 5. Param: '5 def'")
	}
}

func TestExtractTrailingValueAfterColon(t *testing.T) {
	var expected int64
	var result int64

	expected = -1
	if result = ExtractTrailingValueAfterColon("", &logger); result != expected {
		t.Errorf("Invalid ExtractTrailingValueAfterColon result. It should return default value -1. Param: ''")
	}

	expected = -1
	if result = ExtractTrailingValueAfterColon("abc", &logger); result != expected {
		t.Errorf("Invalid ExtractTrailingValueAfterColon result. It should return default value -1. Param: 'abc'")
	}

	expected = -1
	if result = ExtractTrailingValueAfterColon("abc: def", &logger); result != expected {
		t.Errorf("Invalid ExtractTrailingValueAfterColon result. It should return default value -1. Param: 'abc: def'")
	}

	expected = -1
	if result = ExtractTrailingValueAfterColon("abc def: 5 ijk", &logger); result != expected {
		t.Errorf("Invalid ExtractTrailingValueAfterColon result. It should return default value -1. Param: 'abc def: 5 ijk'")
	}

	expected = 5
	if result = ExtractTrailingValueAfterColon("def ijk :5", &logger); result != expected {
		t.Errorf("Invalid ExtractTrailingValueAfterColon result. It should return 5. Param: 'def ijk :5'")
	}

	expected = 5
	if result = ExtractTrailingValueAfterColon("def ijk:    5", &logger); result != expected {
		t.Errorf("Invalid ExtractTrailingValueAfterColon result. It should return 5. Param: 'def ijk:    5'")
	}
}

func TestExtractLastLine(t *testing.T) {
	samples := map[string]string{
		"":                "",
		"abc":             "abc",
		"abc\ndef\nijk":   "ijk",
		"abc\ndef\nijk\n": "ijk",
	}

	for param, expected := range samples {
		if result := ExtractLastLine(param); result != expected {
			t.Errorf("Invalid ExtractLastLine result. Param; '%s', Expected: '%s', Actual: '%s'", param, expected, result)
		}
	}
}

func TestCountLines(t *testing.T) {
	samples := map[string]int{
		"":                0,
		"abc":             1,
		"abc\ndef\nijk":   3,
		"abc\ndef\nijk\n": 3,
	}

	for param, expected := range samples {
		if result := CountLines(param); result != expected {
			t.Errorf("Invalid CountLines result. Param: '%s', Expected: %d, Actual: %d", param, expected, result)
		}
	}
}

func TestBoolToFloat(t *testing.T) {
	if BoolToFloat(true) != 1 {
		t.Errorf("BoolToFloat should convert 'true' to '1'.")
	}

	if BoolToFloat(false) != 0 {
		t.Errorf("BoolToFloat should convert 'false' to '0'.")
	}
}
