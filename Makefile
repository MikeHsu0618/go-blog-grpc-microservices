
.PHONY: protoc
protoc:
	for file in $$(find api -name '*.proto'); do \
		protoc \
		-I $$(dirname $$file) \
		-I ./third_party \
		--go_out=:$$(dirname $$file) --go_opt=paths=source_relative \
		--go-grpc_out=:$$(dirname $$file) --go-grpc_opt=paths=source_relative \
		--validate_out="lang=go:$$(dirname $$file)" --validate_opt=paths=source_relative \
		--grpc-gateway_out=:$$(dirname $$file) --grpc-gateway_opt=paths=source_relative \
		$$file; \
	done

.PHONY: migrate-up
migrate-up:
	migrate -path ./migrations/user -database "postgres://postgres:postgres@localhost:5432/users?sslmode=disable" -verbose up
	migrate -path ./migrations/post -database "postgres://postgres:postgres@localhost:5432/posts?sslmode=disable" -verbose up
	migrate -path ./migrations/comment -database "postgres://postgres:postgres@localhost:5432/comments?sslmode=disable" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path ./migrations/user -database "postgres://postgres:postgres@localhost:5432/users?sslmode=disable" -verbose down --all
	migrate -path ./migrations/post -database "postgres://postgres:postgres@localhost:5432/posts?sslmode=disable" -verbose down --all
	migrate -path ./migrations/comment -database "postgres://postgres:postgres@localhost:5432/comments?sslmode=disable" -verbose down --all

.PHONY: wire
wire:
	wire ./...

.PHONY: auth-server
auth-server:
	go run ./cmd/auth/

.PHONY: user-server
user-server:
	go run ./cmd/user/

.PHONY: post-server
post-server:
	go run ./cmd/post/
