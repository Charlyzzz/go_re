.PHONY: deps
.PHONY: clean
.PHONY: build
.PHONY: test
.PHONY: deploy

deps:
	go get -u ./...

clean:
	rm -rf ./redirect/redirect
	rm -rf ./create/create

build:
	GOOS=linux GOARCH=amd64 go build -o redirect/redirect ./redirect
	GOOS=linux GOARCH=amd64 go build -o create/create ./create

test:
	go test ./create ./redirect

deploy:
	sam build
	sam deploy