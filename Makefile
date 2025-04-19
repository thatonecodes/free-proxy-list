.PHONY: update build

update:
	go build -o gfp cmd/main.go && ./gfp
build:
	go build -v -o gfp cmd/main.go
