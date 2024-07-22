echo "Building filesplitter.exe for Windows"
GOOS="windows" GOARCH="amd64" go build -o filesplitter.exe main.go shared.go
echo "Building filesplitter for Linux"
GOOS="linux" GOARCH="amd64" go build -o filesplitter main.go shared.go
