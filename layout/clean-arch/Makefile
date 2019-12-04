fmt:
	goimports -w -srcdir ./ ./
run:
	PORT=8000 go run *.go

generate:
	../gen/bin/gen generate -c schema.yml -o .