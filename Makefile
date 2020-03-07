build-docker:
	@docker build -t form3-client:interview .

test:
	@docker run -v ${PWD}:/app form3-client:interview

all: build-docker test