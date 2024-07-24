# FileSplitter
Splits files into parts and merges them back

```
FileSplitter splits files into parts and merges them back.
Usage:
  -b uint
        Size of the parts in bytes. Default is 25MB (default 25000000)
  -d string
        Destination directory, where to save the splitted files. Default is 'splitted' (default "splitted")
  -f string
        Required. The source file to split.
  -merge
        Merge the splitted files back to the original file
  -split
        Split the file into parts
```

## How to install?
### Linux
```bash
wget https://github.com/qxsch/FileSplitter/raw/main/filesplitter -O FileSplitter
chmod +x FileSplitter
```
### Windows
```powershell
Invoke-WebRequest -Uri "https://github.com/qxsch/FileSplitter/raw/main/filesplitter.exe" -OutFile "FileSplitter.exe"
Unblock-File -Path ".\FileSplitter.exe"
```


## Splitting a file
To split a file into parts, use the `-split` flag. The `-f` flag is required to specify the source file to split. The `-b` flag can be used to specify the size of the parts in bytes. The default size is 25MB.

Example Linux:
```bash
./FileSplitter -split -f /path/to/file.txt -d /path/to/parts
# or simply
./FileSplitter -f /path/to/file.txt
```
Example Windows:
```powershell
.\FileSplitter.exe -split -f C:\path\to\file.txt -d C:\path\to\parts
# or simply
.\FileSplitter.exe -f C:\path\to\file.txt
```

## Merging the splitted files
To merge the splitted files back to the original file, use the `-merge` flag. The `-f` flag can be used the file to merge. The `-d` flag can be used to specify the parts directory, where the splitted files are located.

Example Linux:
```bash
./FileSplitter -merge -f /path/to/file.txt -d /path/to/parts
# or simply
./FileSplitter -merge
```
Example Windows:
```powershell
.\FileSplitter.exe -merge -f C:\path\to\file.txt -d C:\path\to\parts
# or simply
.\FileSplitter.exe -merge
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
