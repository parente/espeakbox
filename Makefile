.PHONY: build dev-env

REPO:=parente/espeakbox
TAG?=latest
CMD?=

build:
	@docker build -t $(REPO):$(TAG) .

dev-env:
	@docker run -it --rm -p 8080:8080 $(REPO):$(TAG) $(CMD)

