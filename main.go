package main

import (
	"bufio"
	"convertfs/mcio"
	"convertfs/utils"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var formatOption string
var prefixOption string
var helpOption bool
var mtdnaChromosomeOption string

func main() {
	PACKEDPED := "PACKEDPED"
	EIGENSTRAT := "EIGENSTRAT"
	REDIS := "REDIS"

	flag.StringVar(&formatOption, "f", PACKEDPED, "")
	flag.StringVar(&formatOption, "format", PACKEDPED, "")

	flag.StringVar(&prefixOption, "p", "", "")
	flag.StringVar(&prefixOption, "prefix", "", "")

	flag.BoolVar(&helpOption, "h", false, "")
	flag.BoolVar(&helpOption, "help", false, "")

	flag.StringVar(&mtdnaChromosomeOption, "m", "25", "")
	flag.StringVar(&mtdnaChromosomeOption, "mtdna", "25", "")

	setFlag(flag.CommandLine)

	flag.Parse()

	if helpOption {
		utils.ShowHelp()
		return
	}

	genoPath := prefixOption + ".geno"
	indPath := prefixOption + ".ind"
	snpPath := prefixOption + ".snp"
	genoOutPath := prefixOption + ".unpacked.geno"
	indOutPath := prefixOption + ".unpacked.ind"
	snpOutPath := prefixOption + ".unpacked.snp"
	bedOutPath := prefixOption + ".bed"
	famOutPath := prefixOption + ".fam"
	bimOutPath := prefixOption + ".bim"

	if formatOption == PACKEDPED {
		packedAncestryMapToBed(genoPath, indPath, snpPath, bedOutPath, famOutPath, bimOutPath)
	} else if formatOption == EIGENSTRAT {
		packedAncestryMapToEigenstrat(genoPath, indPath, snpPath, genoOutPath, indOutPath, snpOutPath)
	} else if formatOption == REDIS {
		packedAncestryMapToRedis(genoPath, indPath, snpPath)
	}
}

func setFlag(flag *flag.FlagSet) {
	flag.Usage = func() {
		utils.ShowHelp()
	}
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

func packedAncestryMapToRedis(genoPath, indPath, snpPath string) {
	indNum := getRowsNumber(indPath)
	currentByteIndex := 0

	for i := 0; i < indNum; i++ {
		genoFile, err := os.Open(genoPath)
		if err != nil {
			log.Fatal(err)
		}
		defer genoFile.Close()

		genoReader := bufio.NewReader(genoFile)
		chunkSize := int(math.Ceil(float64(indNum) / 4))

		// Why 48? See original EIG/src/mcio.c
		genoChunkSize := chunkSize
		if genoChunkSize < 48 {
			genoChunkSize = 48
		}

		// Skip header
		rchunk := make([]byte, genoChunkSize)
		genoReader.Read(rchunk)

		var redisByte byte
		var redisGeno string

		byteOffset := int(math.Floor(float64(i) / 4))
		bitOffset := uint8(i%4) * 2
		var bitPosition uint8

		for {
			genoByte, err := genoReader.ReadByte()
			if err != nil {
				break
			}

			if currentByteIndex < chunkSize {
				if currentByteIndex == byteOffset {
					snp := (genoByte >> bitOffset) & 3
					redisByte = redisByte | (snp << bitPosition)
					bitPosition += 2

					if bitPosition == 8 {
						redisGeno += string(redisByte)
						bitPosition = 0
						redisByte = 0
					}
				}
			}

			currentByteIndex = (currentByteIndex + 1) % genoChunkSize
		}

		fmt.Println(i, redisGeno)
		// TODO: write redisGeno to redis (+ hashes from ind & snp files)
	}
}

func packedAncestryMapToBed(genoPath, indPath, snpPath, bedOutPath, famOutPath, bimOutPath string) {
	// .bed (see https://www.cog-genomics.org/plink2/formats#bed)
	// bed: magicNumbers + V blocks of math.Ceil(N/4) bytes each
	// V - snp number
	// N - ind number
	// The first block corresponds to the first marker in the .bim file, etc.
	// 0x6c, 0x1b, and 0x01

	bedOutFile, err := os.Create(bedOutPath)
	if err != nil {
		log.Fatal(err)
	}
	defer bedOutFile.Close()

	genoFile, err := os.Open(genoPath)
	if err != nil {
		log.Fatal(err)
	}
	defer genoFile.Close()

	bedWriter := bufio.NewWriter(bedOutFile)
	bedWriter.Write([]byte{0x6c, 0x1b, 0x01})

	genoReader := bufio.NewReader(genoFile)
	indNum := getRowsNumber(indPath)
	chunkSize := int(math.Ceil(float64(indNum) / 4))

	lastOffset := uint8(indNum % 4)
	var lastByteMask uint8
	switch lastOffset {
	case 0:
		lastByteMask = 255
	case 3:
		lastByteMask = 191
	case 2:
		lastByteMask = 175
	case 1:
		lastByteMask = 171
	}

	// Why 48? See original EIG/src/mcio.c
	genoChunkSize := chunkSize
	if genoChunkSize < 48 {
		genoChunkSize = 48
	}

	// Skip header
	rchunk := make([]byte, genoChunkSize)
	genoReader.Read(rchunk)

	currentByteIndex := 0
	for {
		genoByte, err := genoReader.ReadByte()
		if err != nil {
			break
		}

		if currentByteIndex < chunkSize {
			// byte in eigenstrat *.geno is not the same as byte in plinks *.bed
			// 00	Homozygous for first allele in .bim file
			// 01	Missing genotype
			// 10	Heterozygous
			// 11	Homozygous for second allele in .bim file

			// TODO: !!! Treat bytes properly !!!
			// geno -> bed
			//   00 -> 11
			//   01 -> 10
			//   10 -> 00
			//   11 -> 01
			var bedByte, t byte

			for i := uint(0); i < 4; i++ {
				switch d := (genoByte >> (i * 2)) & 3; d {
				case 0:
					t = 3
				case 1:
					t = 2
				case 2:
					t = 0
				case 3:
					t = 1
				}
				bedByte = bedByte | (t << (6 - i*2))
			}

			if currentByteIndex == chunkSize-1 {
				bedByte = bedByte & lastByteMask
			}

			err = bedWriter.WriteByte(bedByte)
			if err != nil {
				log.Fatal(err)
			}
		}

		currentByteIndex = (currentByteIndex + 1) % genoChunkSize
	}

	if err = bedWriter.Flush(); err != nil {
		log.Fatal(err)
	}

	// .fam (see https://www.cog-genomics.org/plink2/formats#fam)
	famOutFile, err := os.Create(famOutPath)
	if err != nil {
		log.Fatal(err)
	}
	defer famOutFile.Close()

	inds := readEigenstratInd(indPath)
	famWriter := bufio.NewWriter(famOutFile)
	for i, ind := range inds {
		var sex string
		switch ind.sex {
		case "M":
			sex = "1"
		case "F":
			sex = "2"
		default:
			sex = "0"
		}
		str := strconv.Itoa(i) + " " + ind.id + " 0 0 " + sex + " 1" + string(10)
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
		// TODO: ignore variants with negative position
		//       for now just replacing minus with empty string
		todoPosition := strings.Replace(snp.position, "-", "", -1)
		str := snp.chromosome + " " + snp.id + " " + todoPosition + " " + snp.coordinate + " " + snp.allele1 + " " + snp.allele2 + string(10)
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

		rawChromosome := fields[1]
		var chromosome string

		switch rawChromosome {
		// TODO:
		// X chromosome is encoded as 23.
		// Also, Y is encoded as 24, mtDNA is encoded as 90, and XY is encoded as 91.
		case "90":
			chromosome = mtdnaChromosomeOption
		default:
			chromosome = rawChromosome
		}

		snps[i] = Snp{
			id:         fields[0],
			chromosome: chromosome,
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
