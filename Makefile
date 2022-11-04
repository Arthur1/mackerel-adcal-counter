.PHONY: build
build:
	go build cmd/mackerel-adcal-counter/main.go

build-for-lambda:
	GOOS=linux GOARCH=amd64 go build cmd/mackerel-adcal-counter/main.go
	zip function.zip main
