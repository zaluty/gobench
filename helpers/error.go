package helpers

import "regexp"

func MustCompile(str string) *regexp.Regexp {
	regexp, err := regexp.Compile(str)
	if err != nil {
		panic(`regexp.Compile(` + str + `): ` + err.Error())
	}
	return regexp
}
