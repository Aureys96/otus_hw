package main

import (
	"fmt"

	strUtil "golang.org/x/example/stringutil"
)

func main() {
	originalString := "Hello, OTUS!"
	reversedString := strUtil.Reverse(originalString)
	fmt.Println(reversedString)
}
