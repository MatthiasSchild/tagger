.PHONY: all linux_amd64 linux_386 linux_arm windows_amd64 windows_386

all: linux_amd64 linux_386 linux_arm windows_amd64 windows_386

linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o out/tagger-linux_amd64 .

linux_386:
	GOOS=linux GOARCH=386 go build -o out/tagger-linux_386 .

linux_arm:
	GOOS=linux GOARCH=arm go build -o out/tagger-linux_arm .

windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o out/tagger-windows_amd64.exe .

windows_386:
	GOOS=windows GOARCH=386 go build -o out/tagger-windows_386.exe .
