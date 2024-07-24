# FileSplitter
Splits files into parts and merges them back

```
FileSplitter splits files into parts and merges them back.
Usage: FileSplitter -split -f <file path> [-d <parts directory>] [-b <parts size>] [-mode <split mode>]
Usage: FileSplitter -merge [-f <file path>] [-d <parts directory>]
Options:
  -b uint
        Size of the parts in bytes. Default is 25MB. (default 25000000)
  -d string
        Destination directory, where to save the splitted files. (default "splitted")
  -f string
        Required for splitting. The source file to split or the destination file to merge.
  -m string
        Split mode. Allowed values: binary, newline. (default "binary")
  -merge
        Merge the splitted files back to the original file
  -split
        Split the file into parts
```

## How to install?
### Linux
```bash
wget https://github.com/qxsch/FileSplitter/raw/main/filesplitter -O fileSplitter
chmod +x filesplitter
```
### Windows
```powershell
Invoke-WebRequest -Uri "https://github.com/qxsch/FileSplitter/raw/main/filesplitter.exe" -OutFile "fileSplitter.exe"
Unblock-File -Path ".\filesplitter.exe"
```


## Splitting a file
To split a file into parts, use the `-split` flag. The `-f` flag is required to specify the source file to split. The `-b` flag can be used to specify the size of the parts in bytes. The default size is 25MB.

Example Linux:
```bash
./filesplitter -split -f /path/to/file.txt -d /path/to/parts
# or simply
./filesplitter -f /path/to/file.txt
```
Example Windows:
```powershell
.\filesplitter.exe -split -f C:\path\to\file.txt -d C:\path\to\parts
# or simply
.\filesplitter.exe -f C:\path\to\file.txt
```

## Merging the splitted files
To merge the splitted files back to the original file, use the `-merge` flag. The `-f` flag can be used the file to merge. The `-d` flag can be used to specify the parts directory, where the splitted files are located.

Example Linux:
```bash
./filesplitter -merge -f /path/to/file.txt -d /path/to/parts
# or simply
./filesplitter -merge
```
Example Windows:
```powershell
.\filesplitter.exe -merge -f C:\path\to\file.txt -d C:\path\to\parts
# or simply
.\filesplitter.exe -merge
```

## How to build?

### Linux
```bash
./build.sh
```
### Windows
```powershell
.\build.ps1
```
