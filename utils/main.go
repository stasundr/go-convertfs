package utils

import "fmt"

// ShowHelp показывает справочную информацию о параметрах запуска программы
func ShowHelp() {
	fmt.Println(`
Version: 0.01a
Options:
    -f, --format   Output file format - PACKEDPED/EIGENSTRAT. Default is PACKEDPED.
    -p, --prefix   Prefix for input/output files.
    -h, --help     Prints this help info.
    `)
}
