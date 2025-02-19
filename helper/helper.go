package helper

import (
	"fmt"
	"strconv"
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

func BuildStatement(statement string, fieldName string, typeField string, valueN int, valueF float64, valueS string) string {

	if (typeField == "S" && len(valueS) == 0) ||
		(typeField == "F" && valueF == 0) ||
		(typeField == "N" && valueN == 0) {

		return statement
	}

	if !strings.HasSuffix(statement, "SET ") {
		statement += ", "
	}

	switch typeField {
	case "S":
		statement += fieldName + " = '" + RemoveQuotationMarks(valueS) + "'"
	case "N":
		statement += fieldName + " = " + strconv.Itoa(valueN)
	case "F":
		statement += fieldName + " = " + strconv.FormatFloat(valueF, 'e', -1, 64)

	}

	fmt.Println("Statement built: " + statement)

	return statement
}
