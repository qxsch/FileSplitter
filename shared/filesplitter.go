package shared

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type FileSplitter struct {
	PartsSize     uint
	PartsDirPath  string
	FilePath      string
	WriteToStdOut bool
}

func NewFileSplitter(PartsSize uint, PartsDirPath string, FilePath string) (*FileSplitter, error) {
	return &FileSplitter{
		PartsSize:     PartsSize,
		PartsDirPath:  PartsDirPath,
		FilePath:      FilePath,
		WriteToStdOut: true,
	}, nil
}

func (fs *FileSplitter) CreateDirectoryIfRequired() (bool, error) {
	var created bool = false
	// check if the parts dir is a directory or not
	fileInfo, err := os.Stat(fs.PartsDirPath)
	if os.IsNotExist(err) {
		if fs.WriteToStdOut {
			fmt.Println("Creating parts dir '" + fs.PartsDirPath + "'")
		}
		err = os.Mkdir(fs.PartsDirPath, 0755)
		if err != nil {
			return false, fmt.Errorf("ERROR: Could not create parts dir '" + fs.PartsDirPath + "'")
		} else {
			created = true
		}
	} else {
		if !fileInfo.IsDir() {
			return false, fmt.Errorf("ERROR: '" + fs.PartsDirPath + "' is not a directory")
		}
	}
	return created, nil
}

func (fs *FileSplitter) CheckRequiredFields() error {
	// check if the file path is empty
	if fs.FilePath == "" {
		return fmt.Errorf("ERROR: File path is required")
	}
	// check if the parts dir path is empty
	if fs.PartsDirPath == "" {
		fs.PartsDirPath = "splitted"
	}
	// check if the parts size is greater than 0
	if fs.PartsSize <= 0 {
		fs.PartsSize = 25000000
	}
	// check if the file exists
	_, err := os.Stat(fs.FilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("ERROR: File '" + fs.FilePath + "' does not exist")
	}
	// check if the parts dir is a directory or not
	fileInfo, err := os.Stat(fs.PartsDirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("ERROR: '" + fs.PartsDirPath + "' does not exist")
	} else {
		if !fileInfo.IsDir() {
			return fmt.Errorf("ERROR: '" + fs.PartsDirPath + "' is not a directory")
		}
	}
	return nil
}

func (fs *FileSplitter) Split() (uint, error) {
	var partCount uint = 0
	if err := fs.CheckRequiredFields(); err != nil {
		return 0, err
	}
	// open the source file
	sourceFile, err := os.Open(fs.FilePath)
	if err != nil {
		return 0, fmt.Errorf("ERROR: Could not open source file '" + fs.FilePath + "'")
	}
	defer sourceFile.Close()

	// split the file
	reader := bufio.NewReader(sourceFile)
	buffer := make([]byte, fs.PartsSize)
	for {
		partCount++
		bytesRead, err := reader.Read(buffer)
		if err != nil {
			break
		}
		if bytesRead > 0 {
			// create the part file
			var partFilePath string = filepath.Join(fs.PartsDirPath, fmt.Sprintf("%s%d%s", "splitted_", partCount, ".bin"))
			partFile, err := os.Create(partFilePath)
			if err != nil {
				fmt.Println("ERROR: Could not create part file '" + partFilePath + "'")
				os.Exit(2)
			}
			// writing to the part file
			bytesWritten, err := partFile.Write(buffer[:bytesRead])
			if err != nil {
				fmt.Println("ERROR: Could not write to part file '" + partFilePath + "'")
				os.Exit(2)
			}
			if bytesWritten != bytesRead {
				fmt.Println("ERROR: Could not write the whole part to the part file '" + partFilePath + "'")
				os.Exit(2)
			}
			partFile.Close()
		}
	}

	// create the file split info
	infoFile, err := os.Create(filepath.Join(fs.PartsDirPath, fmt.Sprintf("%s%s%s", "splitted_", "info", ".json")))
	if err != nil {
		return partCount, fmt.Errorf("ERROR: Could not create the split info file. " + err.Error())
	}
	defer infoFile.Close()
	jsonData, err := json.Marshal(&FileSplitInfo{
		PartCount: partCount,
		FilePath:  fs.FilePath,
	})
	if err != nil {
		return partCount, fmt.Errorf("ERROR: Could not create the split information. " + err.Error())
	}
	_, err = infoFile.Write(jsonData)
	if err != nil {
		return partCount, fmt.Errorf("ERROR: Could not write the split information. " + err.Error())
	}

	return partCount, nil
}
