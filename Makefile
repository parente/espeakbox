.PHONY: build dev-env kill

REPO:=parente/espeakbox
TAG?=latest
CMD?=

build: server
	@docker build -t $(REPO):$(TAG) .
	@rm server

server: server.go
	@GOARCH=amd64 GOOS=linux go build -o server server.go

dev-env:
	@docker run -it --rm -p 8080:8080 $(REPO):$(TAG) $(CMD)

