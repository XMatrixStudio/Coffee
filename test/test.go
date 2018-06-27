package main

import (
	"fmt"
	"strings"
)

func main() {
	name := "file23"
	index := strings.Split(name, "file")
	fmt.Println(index[1])
}