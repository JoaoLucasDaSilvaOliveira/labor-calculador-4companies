package valueobject

import "regexp"

func onlyDigits(value string) string {
	//searches for any character which is not a digit
	regex := regexp.MustCompile(`\D`)
	return regex.ReplaceAllString(value, "")
}
