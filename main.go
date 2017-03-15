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

	// genoPath := "/Users/me/Desktop/v19/v19.0_HO.pruned.geno"
	// indPath := "/Users/me/Desktop/v19/v19.0_HO.pruned.ind"
	// snpPath := "/Users/me/Desktop/v19/v19.0_HO.pruned.snp"
	// genoOutPath := "out.txt"

	// packedAncestryMapToEigenstrat(genoPath, indPath, snpPath, genoOutPath)
	packedAncestryMapToBed()
}

func packedAncestryMapToEigenstrat(genoPath, indPath, snpPath, genoOutPath string) {
	hashOk := mcio.Calcishash(genoPath, indPath, snpPath)
	if !hashOk {
		log.Fatal("Hash check is failed")
	}

	// TODO: copy *.ind and *.snp to indOutPath and snpOutPath
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

	reader.Read(rchunk)
	for {
		_, err := reader.Read(rchunk)
		if (err == nil) || (err == io.EOF) {
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

			if err == io.EOF {
				break
			}
		} else {
			log.Fatal(err)
		}
	}
}

func packedAncestryMapToBed() {
	// .bed
	// see https://www.cog-genomics.org/plink2/formats#bed
	// bed: magicNumbers + V blocks of math.Ceil(N/4) bytes each
	// V - snp number
	// N - ind number
	// The first block corresponds to the first marker in the .bim file, etc.

	// .bim
	// one line per variant with the following six fields:
	// Chromosome code (either an integer, or 'X'/'Y'/'XY'/'MT'; '0' indicates unknown) or name
	// Variant identifier
	// Position in morgans or centimorgans (safe to use dummy value of '0')
	// Base-pair coordinate (normally 1-based, but 0 ok; limited to 231-2)
	// Allele 1 (corresponding to clear bits in .bed; usually minor)
	// Allele 2 (corresponding to set bits in .bed; usually major)
	// Allele codes can contain more than one character. Variants with negative bp coordinates are ignored by PLINK.

	magicNumbers := []byte{0x6c, 0x1b, 0x01}
	fmt.Println(magicNumbers)
}
