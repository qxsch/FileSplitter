package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type FileSplitInfo struct {
	PartCount uint
	FilePath  string
}

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

type FileMerger struct {
	PartsDirPath  string
	FilePath      string
	fileInfo      FileSplitInfo
	hasSplitInfo  bool
	WriteToStdOut bool
}

func NewFileMerger(PartsDirPath string, FilePath string) (*FileMerger, error) {
	return &FileMerger{
		PartsDirPath:  PartsDirPath,
		FilePath:      FilePath,
		WriteToStdOut: true,
		hasSplitInfo:  false,
		fileInfo: FileSplitInfo{
			PartCount: 0,
			FilePath:  "",
		},
	}, nil
}

func (fm *FileMerger) CheckRequiredFields() error {
	fm.hasSplitInfo = false
	// check if the parts dir path is empty
	if fm.PartsDirPath == "" {
		fm.PartsDirPath = "splitted"
	}
	// check if the parts dir is a directory or not
	fileInfo, err := os.Stat(fm.PartsDirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("ERROR: '" + fm.PartsDirPath + "' does not exist")
	} else {
		if !fileInfo.IsDir() {
			return fmt.Errorf("ERROR: '" + fm.PartsDirPath + "' is not a directory")
		}
	}

	// check if the info file exists
	infoFilePath := filepath.Join(fm.PartsDirPath, fmt.Sprintf("%s%s%s", "splitted_", "info", ".json"))
	_, err = os.Stat(infoFilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("ERROR: Info file '" + infoFilePath + "' does not exist")
	}

	// read the info file
	infoFile, err := os.Open(infoFilePath)
	if err != nil {
		return fmt.Errorf("ERROR: Could not open info file '" + infoFilePath + "'")
	}
	defer infoFile.Close()
	jsonParser := json.NewDecoder(infoFile)
	if err = jsonParser.Decode(&fm.fileInfo); err != nil {
		return fmt.Errorf("ERROR: Could not parse the info file '" + infoFilePath + "'")
	}

	fm.hasSplitInfo = true

	return nil
}

func (fm *FileMerger) Merge() (string, error) {
	err := fm.CheckRequiredFields()
	if err != nil {
		if fm.WriteToStdOut {
			fmt.Println(err)
		}
	}
	if fm.FilePath == "" && fm.fileInfo.FilePath != "" {
		fm.FilePath = fm.fileInfo.FilePath
	}
	if fm.FilePath == "" {
		return "", fmt.Errorf("ERROR: File path is required")
	}

	// open the destination file
	destFile, err := os.Create(fm.FilePath)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not create destination file '" + fm.FilePath + "'")
	}
	defer destFile.Close()

	var maxPartCount uint = fm.fileInfo.PartCount
	var silentlyFailOnOpenError bool = false
	if maxPartCount == 0 {
		if !fm.hasSplitInfo {
			if fm.WriteToStdOut {
				fmt.Println("WARNING: Could not find any part count, trying to merge until we fail")
			}
			silentlyFailOnOpenError = true
			maxPartCount = ^uint(0)
		} else {
			return fm.FilePath, nil
		}
	}

	for i := uint(1); i <= maxPartCount; i++ {
		partFilePath := filepath.Join(fm.PartsDirPath, fmt.Sprintf("%s%d%s", "splitted_", i, ".bin"))
		partFile, err := os.Open(partFilePath)
		if err != nil {
			if silentlyFailOnOpenError {
				break
			} else {
				return fm.FilePath, fmt.Errorf("ERROR: Could not open part file '" + partFilePath + "'")
			}
		}
		defer partFile.Close()

		reader := bufio.NewReader(partFile)
		buffer := make([]byte, 8192)
		for {
			bytesRead, err := reader.Read(buffer)
			if err != nil {
				break
			}
			if bytesRead > 0 {
				bytesWritten, err := destFile.Write(buffer[:bytesRead])
				if err != nil {
					return fm.FilePath, fmt.Errorf("ERROR: Could not write to destination file '" + fm.FilePath + "'")
				}
				if bytesWritten != bytesRead {
					return fm.FilePath, fmt.Errorf("ERROR: Could not write the whole part to the destination file '" + fm.FilePath + "'")
				}
			}
		}
	}

	return fm.FilePath, nil
}
