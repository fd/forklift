
build:
	go get -v '.'

dist: build
	mkdir -p dist/forklift-linux-386/bin
	mkdir -p dist/forklift-linux-amd64/bin
	mkdir -p dist/forklift-darwin-386/bin
	mkdir -p dist/forklift-darwin-amd64/bin
	GOOS=linux GOARCH=386 go build -o dist/forklift-linux-386/bin/forklift main.go
	GOOS=linux GOARCH=amd64 go build -o dist/forklift-linux-amd64/bin/forklift main.go
	GOOS=darwin GOARCH=386 go build -o dist/forklift-darwin-386/bin/forklift main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/forklift-darwin-amd64/bin/forklift main.go
	tar -C dist -czf dist/forklift-linux-386.tar.gz  forklift-linux-386
	tar -C dist -czf dist/forklift-linux-amd64.tar.gz  forklift-linux-amd64
	tar -C dist -czf dist/forklift-darwin-386.tar.gz  forklift-darwin-386
	tar -C dist -czf dist/forklift-darwin-amd64.tar.gz  forklift-darwin-amd64

.PHONY: build dist
