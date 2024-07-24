package shared

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileSplitInfo struct {
	PartCount uint
	FilePath  string
}

func WriteFileSplitInfo(fsi FileSplitInfo, filePath string) error {
	// create the file split info
	infoFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ERROR: Could not create the split info file. " + err.Error())
	}
	defer infoFile.Close()
	jsonData, err := json.Marshal(fsi)
	if err != nil {
		return fmt.Errorf("ERROR: Could not create the split information. " + err.Error())
	}
	_, err = infoFile.Write(jsonData)
	if err != nil {
		return fmt.Errorf("ERROR: Could not write the split information. " + err.Error())
	}

	return nil
}

func ReadFileSplitInfo(filePath string) (FileSplitInfo, error) {
	fileInfo := FileSplitInfo{
		PartCount: 0,
		FilePath:  "",
	}
	// read the info file
	infoFile, err := os.Open(filePath)
	if err != nil {
		return fileInfo, fmt.Errorf("ERROR: Could not open info file '" + filePath + "'")
	}
	defer infoFile.Close()
	jsonParser := json.NewDecoder(infoFile)
	if err = jsonParser.Decode(&fileInfo); err != nil {
		return fileInfo, fmt.Errorf("ERROR: Could not parse the info file '" + filePath + "'")
	}

	return fileInfo, nil
}
