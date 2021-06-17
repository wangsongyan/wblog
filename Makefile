build:
	echo "Compiling wblog for Linux OS"
	GOOS=linux GOARCH=arm64 go build -o wblog