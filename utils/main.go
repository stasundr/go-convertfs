package utils

import "fmt"

// ShowHelp показывает справочную информацию о параметрах запуска программы
func ShowHelp() {
	fmt.Println(`
Version: 0.1.0
Options:
    -f, --format   Output file format - PACKEDPED/EIGENSTRAT. Default is PACKEDPED.
    -m, --mtdna    Mitochondrial chromosome code. Default is 25.
    -p, --prefix   Prefix for input/output files.
    -h, --help     Prints this help info.
    `)
}
