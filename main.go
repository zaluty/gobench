package main

import (
	"fmt"

	"github.com/zaluty/buildbeforecommmit/helpers"
)

func main() {
	file := helpers.MustCompile(`^.*\d.txt$`)
	fmt.Println(file)
}
