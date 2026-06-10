#!make
include .env
deploy:
	go build .
	test ! -d "$(BIN_PATH)" || 	mv ./tumblr-dt.exe "${BIN_PATH}\tumblr-dt.exe"
build:
	go build .
test:
	go test  .\...
bench:
	go test -bench=. .\...
fix:
	go fmt  .\...
	go fix  .\...

update:
	go get -u
	
release-windows-amd64:
	env GOOS=windows 
	env GOARCH=amd64 
	go build -o ./release/tumblr-dt.exe main.go
	cd ./release && zip tumblr-dt-windows-amd64.zip ./tumblr-dt.exe
	rm ./release/tumblr-dt.exe
	
release-windows-arm64:
	env GOOS=windows 
	env GOARCH=amd64 
	go build -o ./release/tumblr-dt.exe main.go
	cd ./release && zip tumblr-dt-windows-arm64.zip ./tumblr-dt.exe
	rm ./release/tumblr-dt.exe
	
release-mac-amd64:
	env GOOS=darwin
	env GOARCH=amd64 
	go build -o ./release/tumblr-dt main.go
	cd ./release && zip tumblr-dt-mac-amd64.zip tumblr-dt
	rm ./release/tumblr-dt
	
release-mac-arm64:
	env GOOS=darwin
	env GOARCH=arm64 
	go build -o ./release/tumblr-dt main.go
	cd ./release && zip tumblr-dt-mac-arm64.zip tumblr-dt
	rm ./release/tumblr-dt
	
release-linux-amd64:
	env GOOS=linux
	env GOARCH=amd64 
	go build -o ./release/tumblr-dt main.go
	cd ./release && zip tumblr-dt-linux-amd64.zip tumblr-dt
	rm ./release/tumblr-dt
	
release-linux-arm64:
	env GOOS=linux
	env GOARCH=arm64 
	go build -o ./release/tumblr-dt main.go
	cd ./release && zip tumblr-dt-linux-arm64.zip tumblr-dt
	rm ./release/tumblr-dt
	
clean:
	if test -d ./release; then rm -r ./release; fi 
	
release:
	make clean
	mkdir ./release
	make release-windows-amd64
	make release-windows-arm64
	make release-mac-arm64
	make release-mac-amd64
	make release-linux-amd64
	make release-linux-arm64
