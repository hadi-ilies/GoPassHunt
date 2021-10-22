package hunter

import (
	"errors"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"os"
	"strings"

	//The path/filepath stdlib package provides the handy Walk function. It automatically scans subdirectories
	"path/filepath"

	"github.com/nguyenthenguyen/docx"
	"github.com/xuri/excelize/v2"
	"google.golang.org/api/drive/v3"
)

// type readFile func(path string) error
type Hunter struct {
	folderpath string
	verbose    bool
	files      []string
	//words to hunt
	words []string
	//map of compatiblefiles (files that we gonna read) [extensionfile] -> func
	// compatibleFiles map[string]readFile
	service *drive.Service
}

func NewHunter(folderPath string, verbose bool, words []string) *Hunter {
	return &Hunter{folderpath: folderPath, verbose: verbose, words: words}
}

func (h *Hunter) readTxtFile(path string) error {
	fmt.Println("starting to read:", path)
	b, err := ioutil.ReadFile(path)
	//CHECK RIGHTS access
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			if h.verbose {
				fmt.Println(path, "access denied")
			}
			return nil
		} else {
			return err
		}
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
	// get all sheets of excel
	for _, sheet := range f.WorkBook.Sheets.Sheet {
		rows, err := f.GetRows(sheet.Name)
		// Get all the rows in the Sheet1.
		if err != nil {
			fmt.Println(err)
			return err
		}
		for i, row := range rows {
			for j, colCell := range row {
				for _, word := range h.words {
					if strings.Contains(colCell, word) {
						//row:column
						fmt.Println(path, "sheetName="+sheet.Name+":"+"row="+fmt.Sprint(i)+":"+"col="+fmt.Sprint(j), "word:", "\""+word+"\"", "detected")
					}
					// fmt.Print(colCell, "\t")
				}
			}
			// fmt.Println()
		}
	}
	return nil
}

func (h *Hunter) readGdocFile(path string) error {
	r, err := h.service.Files.List().Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		return fmt.Errorf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		if h.verbose {
			fmt.Println("No files found.")
		}
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
	return nil
}

func (h *Hunter) readGsheetFile(path string) error {
	return nil
}

func (h *Hunter) readMsgFile(path string) error {
	return nil
}

func (h *Hunter) readDocxFile(path string) error {
	// Read from docx file
	r, err := docx.ReadDocxFile(path)
	if err != nil {
		return err
	}
	docx1 := r.Editable()
	b := []byte(docx1.GetContent())
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
			//count lines
			nbLines := strings.Count(s[:idx+len(word)], "</w:t>")
			fmt.Println(path, fmt.Sprint(nbLines)+":"+fmt.Sprint(idx), "word:", "\""+string(s[idx:idx+len(word)])+"\"", "detected")
		}
	}
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
	//diplay all filespath
	// for _, file := range h.files {
	// 	fmt.Println(file)
	// }
	return nil
}

func (h *Hunter) processFile(path string) error {
	var err error
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
		//connect to gdrive
		if h.service == nil {
			if h.service, err = ConnectToGdrive(); err != nil {
				fmt.Println("Error: while connecting to google drive")
				return err
			}
		}
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
		err = h.processFile(h.folderpath)
	}
	return err
}
