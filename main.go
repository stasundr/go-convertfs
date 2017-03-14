package main

import (
	"flag"
	"fmt"
)

func main() {
	// fmt.Println(fmt.Sprintf("%x", mcio.HashFileFirstColumn("full230.ind")))

	var parFile string
	flag.StringVar(&parFile, "p", "", "par file")
	flag.Parse()
	fmt.Println(parFile)
}
