.PHONY: swagger
swagger:
	swag init --parseDependency --parseInternal -g /server/server.go


.PHONY: user-proto
user-proto:
	protoc \
	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	proto/user/user.proto
	

.PHONY: compose
compose:
	docker compose -f docker-compose.yaml up -d