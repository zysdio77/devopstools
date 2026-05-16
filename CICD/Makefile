.PHONY: build run clean

build:
	go build -o bin/cicd ./cmd/cicd

run: build
	./bin/cicd examples/pipeline.yaml

clean:
	rm -rf bin
