package filer

import (
	"fmt"
	selective "go-packages/select"
	"io/fs"
	"log"
	"os"
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
		files = append(files, file.Name())
	}
	
	i, _ := selective.Select(files)
	fmt.Println(filesDir[i].Name())
	info, err := os.Stat(rootPath + filesDir[i].Name())
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(len(options))
	if len(options) > 0 {
		fmt.Println(options[0].Filter)
	}

	if(info.IsDir()) {
		return SelectFile(rootPath + filesDir[i].Name() + "/")
	} else {
		return filesDir[i].Name(), filesDir[i], nil
	}
}
