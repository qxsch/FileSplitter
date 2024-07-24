package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/qxsch/FileSplitter/shared"
)

func main() {
	var (
		partsSize    = flag.Uint("b", 25000000, "Size of the parts in bytes. Default is 25MB.")
		partsDirPath = flag.String("d", "splitted", "Destination directory, where to save the splitted files.")
		filenPath    = flag.String("f", "", "Required for splitting. The source file to split or the destination file to merge.")
		merge        = flag.Bool("merge", false, "Merge the splitted files back to the original file")
		split        = flag.Bool("split", false, "Split the file into parts")
	)

	flag.Parse()

	// in case nothing is selected, use split as default
	if !*split && !*merge {
		*split = true
	}
	// in case both are selected, show an error
	if *split && *merge {
		fmt.Println("ERROR: Either -split or -merge flag is required")
		flag.PrintDefaults()
		os.Exit(2)
	}

	if !*merge {
		// check if the file path is empty
		if *filenPath == "" {
			fmt.Println("ERROR: File path is required")
			flag.PrintDefaults()
			os.Exit(2)
		}
		// check if the parts size is greater than 0
		if *partsSize <= 0 {
			fmt.Println("Setting parts size to 25MB")
			*partsSize = 25000000
		}
	}

	// check if the parts dir path is empty
	if *partsDirPath == "" {
		fmt.Println("Setting parts dir path to 'splitted'")
		*partsDirPath = "splitted"
	}

	if *merge {
		fmt.Println("Merging the splitted files back to the original file")
		fmt.Println("Parts dir path:", *partsDirPath)
		if *filenPath != "" {
			fmt.Println("Source file path:", *filenPath)
		}
		fm, _ := shared.NewFileMerger(*partsDirPath, *filenPath)
		fm.WriteToStdOut = true
		fm.Merge()
	} else {
		fmt.Println("Splitting the file into parts")
		// print the values
		fmt.Println("Parts size:", *partsSize)
		fmt.Println("Parts dir path:", *partsDirPath)
		fmt.Println("Source file path:", *filenPath)

		fs, _ := shared.NewFileSplitter(*partsSize, *partsDirPath, *filenPath)
		if fs == nil {
			fmt.Println("ERROR: Could not create file splitter")
		} else {
			fs.WriteToStdOut = true
			_, err := fs.CreateDirectoryIfRequired()
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
			partCount, err := fs.Split()
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
			fmt.Println("File splitted successfully into ", partCount, " parts")
		}
	}

}
