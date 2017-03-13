package main

import (
	"convertfs/admutils"
	"fmt"
)

func main() {
	fmt.Println(admutils.Hashfile("test.txt"))
}
