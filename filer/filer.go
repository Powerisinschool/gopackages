package gopackages

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/Powerisinschool/gopackages/select"
)

type FileOptions struct {
	Filter string
}

func SelectFile(rootPath string, options ...FileOptions) (string, fs.DirEntry, error) {
	filesDir, err := os.ReadDir(rootPath)
	if err != nil {
		log.Fatal(err)
	}
	var files []string
	for _, file := range filesDir {
		info, err := os.Stat(rootPath + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		if info.IsDir() {
			files = append(files, file.Name()+"/")
			continue
		}
		files = append(files, file.Name())
	}

	i, err := gopackages.Select(files)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(filesDir[i].Name())
	info, err := os.Stat(rootPath + filesDir[i].Name())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(options))
	if len(options) > 0 {
		fmt.Println(options[0].Filter)
	}

	if info.IsDir() {
		return SelectFile(rootPath + filesDir[i].Name() + "/")
	} else {
		return filesDir[i].Name(), filesDir[i], nil
	}
}
