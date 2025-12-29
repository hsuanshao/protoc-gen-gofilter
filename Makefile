
gen-self:
	protoc --go_out=. --go_opt=paths=source_relative proto/filter/filter.proto
\
build: gen-self
	go build -o bin/protoc-gen-gofilter cmd/protoc-gen-gofilter/main.go

install: build
	go install ./cmd/protoc-gen-gofilter