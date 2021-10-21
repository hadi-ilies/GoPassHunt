package hunter

import (
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"os"
	"strings"

	//The path/filepath stdlib package provides the handy Walk function. It automatically scans subdirectories
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

type Hunter struct {
	folderpath string
	verbose    bool
	files      []string
	words      []string //words to hunt
}

func NewHunter(folderPath string, verbose bool, words []string) *Hunter {
	return &Hunter{folderpath: folderPath, verbose: verbose, words: words}
}

func (h *Hunter) readTxtFile(path string) error {
	fmt.Println("starting to read:", path)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	suffix := suffixarray.New(b) // accepts []byte
	for _, word := range h.words {
		indexList := suffix.Lookup([]byte(word), -1)
		if len(indexList) == 0 {
			if h.verbose {
				fmt.Println("the word ", "\"", word, "\" ", "has not been detected inside this file")
			}
			//if word not detected pass to the other one directly
			continue
		}
		s := string(b)
		// loop through the word indices
		for _, idx := range indexList {
			nbLines := strings.Count(s[:idx+len(word)], "\n")
			fmt.Println(path, fmt.Sprint(nbLines)+":"+fmt.Sprint(idx), "word:", "\""+string(s[idx:idx+len(word)])+"\"", "detected")
		}
	}
	return nil
}

func (h *Hunter) readXslxFile(path string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return err
	}
	for i, row := range rows {
		for j, colCell := range row {
			for _, word := range h.words {
				if strings.Contains(colCell, word) {
					//row:column
					fmt.Println(path, "row="+fmt.Sprint(i)+":"+"col=", fmt.Sprint(j), "word:", "\""+word+"\"", "detected")
				}
				// fmt.Print(colCell, "\t")
			}
		}
		// fmt.Println()
	}
	return nil
}

func (h *Hunter) readGdocFile(path string) error {
	return nil
}

func (h *Hunter) readGsheetFile(path string) error {
	return nil
}

func (h *Hunter) readMsgFile(path string) error {
	return nil
}

func (h *Hunter) readDocxFile(path string) error {
	return nil
}

func (h *Hunter) browsePC(path string, info os.FileInfo, err error) error {
	h.files = append(h.files, path)
	//check extension file
	ext := filepath.Ext(path)
	if h.verbose {
		fmt.Println("Handle", ext)
	}
	switch ext {
	case ".txt":
		err = h.readTxtFile(path)
	case ".xlsx":
		err = h.readXslxFile(path)
	case ".gdoc":
		err = h.readGdocFile(path)
	case ".gsheet":
		err = h.readGsheetFile(path)
	case ".msg":
		err = h.readMsgFile(path)
	case ".docx":
		err = h.readDocxFile(path)
	default:
		if h.verbose {
			fmt.Println("no need to read:", path)
		}
	}
	return err
}

func (h *Hunter) processFolder() error {
	//filepath.Walk accepts a string pointing to the root folder thanks to this we get all files inside the targeted folder
	err := filepath.Walk(h.folderpath, h.browsePC)
	if err != nil {
		panic(err)
	}
	for _, file := range h.files {
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
		if h.verbose {
			fmt.Println("Is a Folder")
		}
		err = h.processFolder()
	default:
		if h.verbose {
			fmt.Println("Is a File")
		}
		err = h.processFile()
	}
	return err
}
