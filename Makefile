.PHONY: docs

docs:
	@docker run -v ${PWD}/:/docs pandoc/latex -f markdown /docs/README.md -o /docs/build/output/README.pdf

lint:
	@docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.24.0 golangci-lint run

test:
	@docker-compose up