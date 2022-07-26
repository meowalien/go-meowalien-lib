
lint:
	golangci-lint run

format:
	gofmt -s -w .

mock :
	mockgen -source ./arangodb/Arangodb.go -destination ./arangodb/Arangodb_mock.go -imports driver=github.com/arangodb/go-driver -package arangodb

# make mockgen file=<file_path>
mockgen:
	./internal/scripts/mockgen.sh