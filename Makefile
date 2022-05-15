redeployed-at:=$(shell date +%s)

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

.PHONY: comment-server
comment-server:
	go run ./cmd/comment/

.PHONY: blog-server
blog-server:
	go run ./cmd/blog/

.PHONY: docker-build
docker-build:
	docker build -t blog/blog-server:latest -f ./build/docker/blog/Dockerfile .
	docker build -t blog/user-server:latest -f ./build/docker/user/Dockerfile .
	docker build -t blog/auth-server:latest -f ./build/docker/auth/Dockerfile .
	docker build -t blog/post-server:latest -f ./build/docker/post/Dockerfile .
	docker build -t blog/comment-server:latest -f ./build/docker/comment/Dockerfile .

.PHONY: kube-deploy
kube-deploy:
	kubectl apply -f ./deployments/
	kubectl apply -f ./deployments/dtm/
	kubectl apply -f ./deployments/blog/
	kubectl apply -f ./deployments/user/
	kubectl apply -f ./deployments/post/
	kubectl apply -f ./deployments/auth/
	kubectl apply -f ./deployments/comment/
	kubectl apply -f ./deployments/addons/

.PHONY: kube-delete
kube-delete:
	kubectl delete -f ./deployments/
	kubectl delete -f ./deployments/dtm/
	kubectl delete -f ./deployments/blog/
	kubectl delete -f ./deployments/user/
	kubectl delete -f ./deployments/post/
	kubectl delete -f ./deployments/auth/
	kubectl delete -f ./deployments/comment/
	#kubectl delete -f ./deployments/addons/

.PHONY: kube-redeploy
kube-redeploy:
	@echo "redeployed at ${redeployed-at}"
	kubectl patch deployment blog-server -p '{"spec": {"template": {"metadata": {"annotations": {"redeployed-at": "'${redeployed-at}'" }}}}}'
	kubectl patch deployment user-server -p '{"spec": {"template": {"metadata": {"annotations": {"redeployed-at": "'${redeployed-at}'" }}}}}'
	kubectl patch deployment auth-server -p '{"spec": {"template": {"metadata": {"annotations": {"redeployed-at": "'${redeployed-at}'" }}}}}'
	kubectl patch deployment post-server -p '{"spec": {"template": {"metadata": {"annotations": {"redeployed-at": "'${redeployed-at}'" }}}}}'
	kubectl patch deployment comment-server -p '{"spec": {"template": {"metadata": {"annotations": {"redeployed-at": "'${redeployed-at}'" }}}}}'