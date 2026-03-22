.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -buildvcs=false -trimpath -o yamubot
