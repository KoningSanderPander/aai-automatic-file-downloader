package main

import (
	"autodownloader/internal/download"
	"fmt"
)

func main() {
	fmt.Println("Start downloading...")
	download.SaveFiles()
	fmt.Println("Finished downloading")
}
