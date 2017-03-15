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
	genoOutPath := "out.txt"

	// TODO: get indNum
	indNum := 435

	genoFile, err := os.Open(genoPath)
	if err != nil {
		log.Fatal(err)
	}
	defer genoFile.Close()

	genoOutFile, err := os.Create(genoOutPath)
	if err != nil {
		panic(err)
	}
	defer genoOutFile.Close()

	chunkSize := int(math.Ceil(float64(indNum) / 4))
	rchunk := make([]byte, chunkSize)
	wchunk := make([]byte, indNum+1)
	wchunk[indNum] = 10

	reader := bufio.NewReaderSize(genoFile, chunkSize)
	writer := bufio.NewWriterSize(genoOutFile, indNum+1)

	// packed eigenstrat to unpacked eigenstrat core
	// for ideas see https://github.com/stasundr/co-huge-converter/blob/master/modules/utils.js

	// TODO: skip header!
	for {
		_, err := reader.Read(rchunk)
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		for i, b := range rchunk {
			for j := uint(0); j < 4; j++ {
				snp := (b >> (6 - j*2)) & 3
				if k := i*4 + int(j); k < indNum {
					switch snp {
					case 0:
						wchunk[k] = byte('0')
					case 1:
						wchunk[k] = byte('1')
					case 2:
						wchunk[k] = byte('2')
					case 3:
						wchunk[k] = byte('9')
					}

				}
			}
		}

		writer.Write(wchunk)
	}

	hashOk := mcio.Calcishash(genoPath, indPath, snpPath)
	fmt.Println(hashOk, chunkSize)
}
