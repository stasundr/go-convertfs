package main

import (
	"convertfs/admutils"
	"convertfs/mcio"
	"fmt"
)

func main() {
	fmt.Println(fmt.Sprintf("%x", admutils.Hashfile("hash_is_6111c60d.txt")))
	fmt.Println(fmt.Sprintf("%x", mcio.HashFileFirstColumn("full230.ind")))
}
