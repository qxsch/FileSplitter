Write-Host "Building filesplitter.exe for Windows"
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o filesplitter.exe main.go shared.go
Write-Host "Building filesplitter for Linux"
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o filesplitter main.go shared.go
