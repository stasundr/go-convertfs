package mcio

import (
	"bufio"
	"convertfs/admutils"
	"fmt"
	"io"
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

// Calcishash takes paths for *.geno, *.ind and *.snp files (eigenstrat combo)
// and calculate hashes on individuals and SNPs (to compare with file values).
// (see https://github.com/DReichLab/EIG/blob/master/src/mcio.c#L2697)
func Calcishash(genoPath, indPath, snpPath string) bool {
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

	return hashOk
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
// (see http://stackoverflow.com/a/21067803)
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func main() {
	fmt.Printf("Copying %s to %s\n", os.Args[1], os.Args[2])
	err := CopyFile(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("CopyFile failed %q\n", err)
	} else {
		fmt.Printf("CopyFile succeeded\n")
	}
}
