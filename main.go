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
	"strconv"
	"strings"
)

func main() {
	// TODO: par file parse
	var parFile string
	flag.StringVar(&parFile, "p", "", "par file")
	flag.Parse()
	fmt.Println(parFile)

	// genoPath := "/Users/me/Desktop/v19/v19.0_HO.pruned.geno"
	// indPath := "/Users/me/Desktop/v19/v19.0_HO.pruned.ind"
	snpPath := "/Users/me/Desktop/v19/v19.0_HO.pruned.snp"
	// genoOutPath := "out.geno"
	// indOutPath := "out.ind"
	// snpOutPath := "out.snp"
	bimOutPath := "out.bim"

	// packedAncestryMapToEigenstrat(genoPath, indPath, snpPath, genoOutPath, indOutPath, snpOutPath)
	packedAncestryMapToBed(snpPath, bimOutPath)

	snps := readSnps(snpPath)

	fmt.Println(snps[100])
}

func packedAncestryMapToEigenstrat(genoPath, indPath, snpPath, genoOutPath, indOutPath, snpOutPath string) {
	hashOk := mcio.Calcishash(genoPath, indPath, snpPath)
	if !hashOk {
		log.Fatal("Hash check is failed")
	}

	err := mcio.CopyFile(indPath, indOutPath)
	if err != nil {
		log.Fatal(err)
	}
	err = mcio.CopyFile(snpPath, snpOutPath)
	if err != nil {
		log.Fatal(err)
	}

	indNum := getRowsNumber(indPath)

	genoFile, err := os.Open(genoPath)
	if err != nil {
		log.Fatal(err)
	}
	defer genoFile.Close()

	genoOutFile, err := os.Create(genoOutPath)
	if err != nil {
		log.Fatal(err)
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

func packedAncestryMapToBed(snpPath, bimOutPath string) {
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

	bimOutFile, err := os.Create(bimOutPath)
	if err != nil {
		log.Fatal(err)
	}
	defer bimOutFile.Close()

	i := 0
	snps := readSnps(snpPath)
	bimWriter := bufio.NewWriter(bimOutFile)
	for _, snp := range snps {
		i++
		// TODO: ignore variants with negative bp coordinates
		str := strconv.Itoa(i) + " " + snp.chromosome + " " + snp.id + " " + snp.position + " " + snp.coordinate + " " + snp.allele1 + " " + snp.allele2 + string(10)
		bimWriter.WriteString(str)
		// some kind of a bug here!
	}

	// magicNumbers := []byte{0x6c, 0x1b, 0x01}
	fmt.Println(i, len(snps))
}

func getRowsNumber(path string) int {
	var num int
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		num++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return num
}

// TODO: correct name for function? read eigenstrat snp?
func readSnps(path string) []Snp {
	size := getRowsNumber(path)
	snps := make([]Snp, size)

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(strings.Trim(scanner.Text(), " "))
		// coordinate, err := strconv.Atoi(fields[3])
		// if err != nil {
		// 	log.Fatal(err)
		// }

		snps[i] = Snp{
			id:         fields[0],
			chromosome: fields[1],
			position:   fields[2],
			coordinate: fields[3], //coordinate,
			allele1:    fields[4], //[0],
			allele2:    fields[5]} //[0]}
		i++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return snps
}

// Snp is a type for storing information about single nucleotide polymorphysm
type Snp struct {
	id, chromosome, position, allele1, allele2, coordinate string
	// allele1, allele2         byte
	// coordinate               int
}
