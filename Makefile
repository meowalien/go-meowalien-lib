
lint:
	golangci-lint run

format:
	gofmt -s -w .

# make mock file=<file_path>
mock:
	./internal/scripts/mockgen.sh