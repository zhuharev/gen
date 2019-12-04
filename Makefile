tools:
	go get modernc.org/assets

build: generate
	@go build -o ./bin/gen ./cmd/gen/*.go

.PHONY: build

all: build
	@go run cmd/gen/gen.go

play: build 
	cd play && ../bin/gen generate -c cfg.yml -o generated && goimports -w -srcdir ./ ./ && cd ..

.PHONY: play

generate:
	assets -d templates -package gen -o gen_templates.go