package admutils

import (
	"bufio"
	"log"
	"os"
)

// Hashit takes a string as argument
// and calculates hash for it
// (see https://github.com/DReichLab/EIG/blob/master/src/admutils.c :683)
func Hashit(str string) int32 {
	var hash, length int32 = 0, int32(len(str))

	for j := int32(0); j < length; j++ {
		hash *= 23
		hash += int32(str[j])
	}

	return hash
}

// Hasharr takes an array of strings as argument
// and calculates hash for it
// (see https://github.com/DReichLab/EIG/blob/master/src/admutils.c :666)
func Hasharr(xarr []string) int32 {
	var hash, thash, nxarr int32
	nxarr = int32(len(xarr))

	for i := int32(0); i < nxarr; i++ {
		thash = Hashit(xarr[i])
		hash *= 17
		hash ^= thash
	}

	return hash
}

// Hashfile takes a file path string as argument
// and calculates hash for the file
func Hashfile(path string) int32 {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var hash, thash int32
	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadString(10)
		if err != nil {
			thash = Hashit(str)
			hash *= 17
			hash ^= thash
			return hash
		}

		thash = Hashit(str[0 : len(str)-1])
		hash *= 17
		hash ^= thash
	}
}
