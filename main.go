package main

import (
	"fmt"
	"go-packages/filer"
	"log"
)

func main() {

	_ = []string{"Hey", "You", "Select"}

	fileName, file, err := filer.SelectFile("./", filer.FileOptions{Filter: "Yooo"})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fileName)
	fmt.Println(file)
}
