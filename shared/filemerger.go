package shared

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

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
	fm.fileInfo, err = ReadFileSplitInfo(infoFilePath)
	if err != nil {
		return err
	}
	fm.hasSplitInfo = true

	return nil
}

func (fm *FileMerger) Merge() (string, uint, error) {
	var restoredPartNum uint = 0

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
		return "", restoredPartNum, fmt.Errorf("ERROR: File path is required")
	}

	// open the destination file
	destFile, err := os.Create(fm.FilePath)
	if err != nil {
		return "", restoredPartNum, fmt.Errorf("ERROR: Could not create destination file '" + fm.FilePath + "'")
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
			return fm.FilePath, restoredPartNum, nil
		}
	}

	for restoredPartNum = uint(1); restoredPartNum <= maxPartCount; restoredPartNum++ {
		partFilePath := filepath.Join(fm.PartsDirPath, fmt.Sprintf("%s%d%s", "splitted_", restoredPartNum, ".bin"))
		partFile, err := os.Open(partFilePath)
		if err != nil {
			if silentlyFailOnOpenError {
				restoredPartNum--
				break
			} else {
				return fm.FilePath, restoredPartNum, fmt.Errorf("ERROR: Could not open part file '" + partFilePath + "'")
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
					return fm.FilePath, restoredPartNum, fmt.Errorf("ERROR: Could not write to destination file '" + fm.FilePath + "'")
				}
				if bytesWritten != bytesRead {
					return fm.FilePath, restoredPartNum, fmt.Errorf("ERROR: Could not write the whole part to the destination file '" + fm.FilePath + "'")
				}
			}
		}
	}
	return fm.FilePath, restoredPartNum, nil
}
