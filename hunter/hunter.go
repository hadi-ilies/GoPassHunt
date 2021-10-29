package hunter

import (
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"os"
	"strings"

	//The path/filepath stdlib package provides the handy Walk function. It automatically scans subdirectories
	"path/filepath"

	"github.com/iafan/cwalk"
	"github.com/nguyenthenguyen/docx"
	"github.com/xuri/excelize/v2"
	"google.golang.org/api/drive/v3"
)

// type readFile func(path string) error
type Hunter struct {
	folderpath string
	verbose    bool
	Gdrive     bool
	files      []string
	driveData  []Gdrive
	//words to hunt
	words []string
	//map of compatiblefiles (files that we gonna read) [extensionfile] -> func
	// compatibleFiles map[string]readFile
	service *drive.Service
}

func NewHunter(folderPath string, verbose bool, gDrive bool, words []string) *Hunter {
	return &Hunter{folderpath: folderPath, verbose: verbose, Gdrive: gDrive, words: words}
}

func (h *Hunter) readTxtFile(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		// if errors.Is(err, os.ErrPermission) {
		if h.verbose {
			fmt.Println("Error:", path, err)
		}
		return nil
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
		if h.verbose {
			fmt.Println("Error:", path, err)
		}
		return nil
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

func (h *Hunter) displayGdriveFiles() error {
	for _, data := range h.driveData {
		fmt.Printf("FileName=%s | FileType=(%s) | WordDetected=\"%s\"\n", data.Name, data.Type, data.WordDetected)
	}
	return nil
}

func (h *Hunter) searchInGdrive(fileType string) error {
	if h.verbose {
		fmt.Println("Looking for", "\""+fileType+"\"", "files:")
	}
	for _, word := range h.words {
		r, err := h.service.Files.List().Q("mimeType = '" + fileType + "' and fullText contains '" + word + "'").Fields("nextPageToken, files(id, name)").Do()
		if err != nil {
			return fmt.Errorf("Unable to retrieve files: %v", err)
		}
		if len(r.Files) == 0 {
			if h.verbose {
				fmt.Println("No files found.")
			}
		} else {
			for _, i := range r.Files {
				h.driveData = append(h.driveData, Gdrive{Link: i.IconLink, Type: fileType, Name: i.Name, Id: i.Id, WordDetected: word})
			}
		}
	}
	return nil
}

func (h *Hunter) readMsgFile(path string) error {
	return nil
}

func (h *Hunter) readDocxFile(path string) error {
	// Read from docx file
	r, err := docx.ReadDocxFile(path)
	if err != nil {
		if h.verbose {
			fmt.Println("docx Error:", err)
		}
		return nil
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
	path = h.folderpath + "/" + path
	//TODO(Hadi): should we save paths ?
	// h.files = append(h.files, path)
	//check extension file
	ext := filepath.Ext(path)
	if h.verbose {
		fmt.Println("Handle", ext)
	}
	err = nil
	if h.verbose && (ext == ".txt" || ext == ".xlsx" || ext == ".docx" || ext == "msg") {
		fmt.Println("starting to read:", path)
	}
	switch ext {
	case ".txt":
		err = h.readTxtFile(path)
	case ".xlsx":
		err = h.readXslxFile(path)
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
	//we use walk with goroutines, Check: https://github.com/iafan/cwalk for more info
	err := cwalk.Walk(h.folderpath, h.browsePC)
	return err
}

//connect to gdrive + search all files that contains words inserted inside Hunter struct
//it takes the type of file we want to look at
func (h *Hunter) lookAtGdrive(filesType string) error {
	var err error
	//connect to gdrive
	if h.service == nil {
		if h.service, err = ConnectToGdrive(); err != nil {
			fmt.Println("Error: while connecting to google drive")
			return err
		}
	}
	err = h.searchInGdrive(filesType)
	if err != nil {
		return err
	}
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

func (h *Hunter) searchInGdriveAllFiles() error {
	if err := h.searchInGdrive(GdriveDocxType); err != nil {
		return err
	}
	if err := h.searchInGdrive(GdriveGdocType); err != nil {
		return err
	}
	if err := h.searchInGdrive(GdriveGsheetType); err != nil {
		return err
	}
	if err := h.searchInGdrive(GdriveXlsxType); err != nil {
		return err
	}
	return nil
}

func (h *Hunter) readAllGdrive() error {
	if err := h.displayGdriveFiles(); err != nil {
		return err
	}
	return nil
}

func (h *Hunter) Start() error {
	if h.Gdrive {
		var err error

		fmt.Println("Google Drive: ")
		//connect once to google drive
		if h.service, err = ConnectToGdrive(); err != nil {
			fmt.Println("Error: while connecting to google drive")
			return err
		}
		//look for the wanted types
		if err := h.searchInGdriveAllFiles(); err != nil {
			return err
		}
		if err := h.readAllGdrive(); err != nil {
			return err
		}
	}

	//look for file/folder on current computer
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
