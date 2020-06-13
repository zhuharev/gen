tools:
	go get modernc.org/assets

build: generate
	#rm ./bin/gen
	#rm cmd/gen/pkged.go
	@go install ./...
	@go build -o ./bin/gen ./cmd/gen/*.go

.PHONY: build

all: build
	@go install ./...
	@go run cmd/gen/gen.go

play: build 
	rm -rf play/generated
	pkger -o cmd/gen
	cd play && ../bin/gen generate -c cfg.yml -o generated && cd ..
	# && goimports -w -srcdir ./ ./ && cd ..

.PHONY: play

generate:
	assets -d templates -package gen -o gen_templates.go


.PHONY: generate