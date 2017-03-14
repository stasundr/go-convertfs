package main

import (
	"convertfs/mcio"
	"flag"
	"fmt"
)

func main() {
	// fmt.Println(fmt.Sprintf("%x", mcio.HashFileFirstColumn("full230.ind")))

	var parFile string
	flag.StringVar(&parFile, "p", "", "par file")
	flag.Parse()
	fmt.Println(parFile)

	mcio.Calcishash("/Users/me/Desktop/v19/v19.0_HO.pruned.geno", "/Users/me/Desktop/v19/v19.0_HO.pruned.ind", "/Users/me/Desktop/v19/v19.0_HO.pruned.snp")
}
