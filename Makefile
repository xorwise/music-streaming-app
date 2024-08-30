gen:
	@protoc -I proto ./auth/proto/*.proto --go_out=./auth/gen/ --go_opt=paths=source_relative --go-grpc_out=./auth/gen/ --go-grpc_opt=paths=source_relative

run:
	@docker compose --env-file ./auth/.env up -d --build

stop:
	@docker compose down

test:
	@docker exec -i auth sh -c "go test -v ./tests/auth/..."
