build:
	GOARCH=arm CGO_ENABLED=0 GOARM=7 go build -o webapp main.go
install:
	scp webapp root@10.21.8.77:/root/web
