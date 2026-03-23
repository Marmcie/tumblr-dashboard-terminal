#!make
include .env
deploy:
	go build .
	mv ./tumblr-dt.exe "${BIN_PATH}\tumblr-dt.exe"
test:
	go test -v .\...
bench:
	go test -v -bench=. .\...
