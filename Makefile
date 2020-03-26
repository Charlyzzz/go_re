.PHONY: deps
.PHONY: clean
.PHONY: build
.PHONY: test
.PHONY: deploy

deps:
	go get -u ./...

clean:
	rm -rf ./redirect/redirect

build:
	GOOS=linux GOARCH=amd64 go build -o redirect/redirect ./redirect

test:
	go test ./...

deploy:
	sam build
	sam deploy