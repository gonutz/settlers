cd images
go run pack.go -s=1024
cd ..

go build -ldflags -H=windowsgui

settlers.exe