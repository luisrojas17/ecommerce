package helper

import (
	"fmt"
	"strings"
	"time"
)

const EMPTY_STRING = ""

func GetDate() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func RemoveQuotationMarks(str string) string {

	desc := strings.ReplaceAll(str, "'", EMPTY_STRING)

	desc = strings.ReplaceAll(desc, "\"", EMPTY_STRING)

	return desc
}
