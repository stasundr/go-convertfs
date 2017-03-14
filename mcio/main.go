package mcio

import (
	"bufio"
	"convertfs/admutils"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

// HashFileFirstColumn takes a file path string as argument
// and calculates hash for the file's first column (see *.snp or *.ind file)
// File must have last empty line
func HashFileFirstColumn(path string) uint32 {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var hash, thash uint32
	reader := bufio.NewReader(file)
	firstColumnRegexp, _ := regexp.Compile(`\S+`)
	for {
		str, err := reader.ReadString(10)
		if err != nil {
			return hash
		}

		thash = admutils.HashIt(firstColumnRegexp.FindString(str))
		hash *= 17
		hash ^= thash
	}
}

// Calcishash takes ...
// and calculate hashes on individuals and SNPs (to compare with file values).
// (see https://github.com/DReichLab/EIG/blob/master/src/mcio.c#L2697)
func Calcishash(genoPath, indPath, snpPath string) {
	file, err := os.Open(genoPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ok := scanner.Scan()
	if !ok {
		log.Fatal("Can't read geno file")
	}
	header := scanner.Text()
	indHash := strconv.FormatInt(int64(HashFileFirstColumn(indPath)), 16)
	snpHash := strconv.FormatInt(int64(HashFileFirstColumn(snpPath)), 16)
	indRegexp, _ := regexp.Compile(indHash)
	snpRegexp, _ := regexp.Compile(snpHash)
	hashOk := indRegexp.MatchString(header) && snpRegexp.MatchString(header) && (indPath != snpPath)
	fmt.Println("Hash OK: ", hashOk)
}
