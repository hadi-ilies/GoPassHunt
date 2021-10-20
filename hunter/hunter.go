package hunter

import (
	"fmt"
	"os"

	//The path/filepath stdlib package provides the handy Walk function. It automatically scans subdirectories
	"path/filepath"
)

type Hunter struct {
	folderpath string
	verbose    bool
}

func NewHunter(folderPath string, verbose bool) *Hunter {
	return &Hunter{folderpath: folderPath, verbose: verbose}
}

func (h *Hunter) processFolder() error {
	var files []string
	//filepath.Walk accepts a string pointing to the root folder thanks to this we get all files inside this folder
	err := filepath.Walk(h.folderpath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		//check extension file
		switch ext := filepath.Ext(path); ext {
		case ".txt":
			fmt.Println("handle = ", ext)
		case ".xlsx":
			fmt.Println("handle = ", ext)
		case ".gdoc":
			fmt.Println("handle = ", ext)
		case ".gsheet":
			fmt.Println("handle = ", ext)
		case ".msg":
			fmt.Println("handle = ", ext)
		case ".docx":
			fmt.Println("handle = ", ext)
		default:
			fmt.Println("no need to read this file")
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
	return nil
}

func (h *Hunter) processFile() error {
	return nil
}

func (h *Hunter) Start() error {
	file, err := os.Open(h.folderpath)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	switch {
	case err != nil:
		return err
	case fi.IsDir():
		fmt.Println("Is a Folder")
		err = h.processFolder()
	default:
		fmt.Println("Is a File")
		err = h.processFile()
	}
	return err
}
