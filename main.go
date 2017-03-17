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
	indPath := "/Users/me/Desktop/v19/v19.0_HO.pruned.ind"
	snpPath := "/Users/me/Desktop/v19/v19.0_HO.pruned.snp"
	// genoOutPath := "out.geno"
	// indOutPath := "out.ind"
	// snpOutPath := "out.snp"
	famOutPath := "out.fam"
	bimOutPath := "out.bim"

	// packedAncestryMapToEigenstrat(genoPath, indPath, snpPath, genoOutPath, indOutPath, snpOutPath)
	packedAncestryMapToBed(indPath, snpPath, famOutPath, bimOutPath)
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

func packedAncestryMapToBed(indPath, snpPath, famOutPath, bimOutPath string) {
	// .bed (see https://www.cog-genomics.org/plink2/formats#bed)
	// bed: magicNumbers + V blocks of math.Ceil(N/4) bytes each
	// V - snp number
	// N - ind number
	// The first block corresponds to the first marker in the .bim file, etc.

	// .fam (see https://www.cog-genomics.org/plink2/formats#fam)
	// Family ID ('FID')
	// Within-family ID ('IID'; cannot be '0')
	// Within-family ID of father ('0' if father isn't in dataset)
	// Within-family ID of mother ('0' if mother isn't in dataset)
	// Sex code ('1' = male, '2' = female, '0' = unknown)
	// Phenotype value ('1' = control, '2' = case, '-9'/'0'/non-numeric = missing data if case/control)
	famOutFile, err := os.Create(famOutPath)
	if err != nil {
		log.Fatal(err)
	}
	defer famOutFile.Close()

	inds := readEigenstratInd(indPath)
	famWriter := bufio.NewWriter(famOutFile)
	for i, ind := range inds {
		// TODO: convert sex from "F/M/other" to 0/1/2
		str := strconv.Itoa(i) + " " + ind.id + " 0 0 " + ind.sex + " 1" + string(10)
		famWriter.WriteString(str)
	}

	if err = famWriter.Flush(); err != nil {
		log.Fatal(err)
	}

	// .bim (see https://www.cog-genomics.org/plink2/formats#bim)
	// Allele 1 (usually minor)
	// Allele 2 (usually major)
	// Allele codes can contain more than one character. Variants with negative bp coordinates are ignored by PLINK.
	bimOutFile, err := os.Create(bimOutPath)
	if err != nil {
		log.Fatal(err)
	}
	defer bimOutFile.Close()

	snps := readEigenstratSnp(snpPath)
	bimWriter := bufio.NewWriter(bimOutFile)
	for _, snp := range snps {
		// TODO: ignore variants with negative bp coordinates
		str := snp.chromosome + " " + snp.id + " " + snp.position + " " + snp.coordinate + " " + snp.allele1 + " " + snp.allele2 + string(10)
		bimWriter.WriteString(str)
	}

	if err = bimWriter.Flush(); err != nil {
		log.Fatal(err)
	}
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

func readEigenstratSnp(path string) []Snp {
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

		snps[i] = Snp{
			id:         fields[0],
			chromosome: fields[1],
			position:   fields[2],
			coordinate: fields[3],
			allele1:    fields[4],
			allele2:    fields[5]}
		i++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return snps
}

func readEigenstratInd(path string) []Individual {
	size := getRowsNumber(path)
	inds := make([]Individual, size)

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(strings.Trim(scanner.Text(), " "))

		inds[i] = Individual{
			id:    fields[0],
			sex:   fields[1],
			label: fields[2]}
		i++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return inds
}

// Snp is a type for storing information about single nucleotide polymorphysm
type Snp struct {
	id, chromosome, position, allele1, allele2, coordinate string
	// allele1, allele2         byte
	// coordinate               int
}

// Individual is a type for storing information about sample
type Individual struct {
	id, sex, label string
}
