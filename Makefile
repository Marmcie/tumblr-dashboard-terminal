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
