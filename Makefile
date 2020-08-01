OUT_DIR=bin

build:
	go build -o $(OUT_DIR) ./cmd/...

clean:
	rm $(OUT_DIR)/*

test:
	go test -v ./...
