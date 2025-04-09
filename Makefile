.PHONY: update

update:
	go build -o gfp cmd/main.go && ./gfp
