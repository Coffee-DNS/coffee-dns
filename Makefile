ALLDOC := $(shell find . \( -name "*.md" -o -name "*.yaml" \) \
                              -type f | sort)

.PHONY: install-tools
install-tools:
	cd internal/tools && go install github.com/mgechev/revive
	cd internal/tools && go install github.com/client9/misspell/cmd/misspell
	cd internal/tools && go install github.com/securego/gosec/v2/cmd/gosec
	cd internal/tools && go install golang.org/x/tools/cmd/goimports
	cd internal/tools && go install honnef.co/go/tools/cmd/staticcheck
	cd internal/tools && go install github.com/uw-labs/lichen

.PHONY: gomoddownload
gomoddownload:
	go mod download

.PHONY: quick
quick:
	cd cmd/controller && go build -o coffee

.PHONY: test
test:
	go test ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: tidy
tidy:
	rm -fr go.sum
	go mod tidy

.PHONY: protobuf.controller.generate
protobuf.controller.generate:
	cd controller/api/ && \
		protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		./controller.proto

.PHONY: protobuf.nameserver.generate
protobuf.nameserver.generate:
	cd nameserver/api/ && \
		protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		./nameserver.proto

.PHONY: protobuf.generate
protobuf.generate: protobuf.nameserver.generate protobuf.controller.generate

.PHONY: shellcheck
shellcheck:
	shellcheck scripts/*

.PHONY: build
build:
	goreleaser release --rm-dist --snapshot --skip-publish

.PHONY: run
run: build
	docker-compose up -d --remove-orphans --force-recreate

.PHONY: misspell
misspell:
	misspell -error $(ALLDOC)

.PHONY: misspell-fix
misspell-fix:
	misspell -w $(ALLDOC)

.PHONY: lint
lint:
	revive -config .revive.toml -formatter friendly ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...
