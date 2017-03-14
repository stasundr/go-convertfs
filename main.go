package main

import (
	"bufio"
	"convertfs/mcio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

func main() {
	// TODO: par file parse
	var parFile string
	flag.StringVar(&parFile, "p", "", "par file")
	flag.Parse()
	fmt.Println(parFile)

	genoPath := "/Users/me/Desktop/v19/v19.0_HO.pruned.geno"
	indPath := "/Users/me/Desktop/v19/v19.0_HO.pruned.ind"
	snpPath := "/Users/me/Desktop/v19/v19.0_HO.pruned.snp"

	// TODO: get indNum
	indNum := 435
	chunkSize := int(math.Ceil(float64(indNum) / 4))
	chunk := make([]byte, chunkSize)

	file, err := os.Open(genoPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, chunkSize)
	// packed eigenstrat to unpacked eigenstrat core
	// for ideas see https://github.com/stasundr/co-huge-converter/blob/master/modules/utils.js
	for {
		_, err := reader.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		for i, b := range chunk {
			for j := uint(0); j < 4; j++ {
				snp := (b >> (3 - j)) & 3
				fmt.Println(i, b, snp)
			}
		}
	}

	hashOk := mcio.Calcishash(genoPath, indPath, snpPath)
	fmt.Println(hashOk, chunkSize)
}
