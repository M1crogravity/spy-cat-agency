test:
	go test ./...
run:
	go run ./cmd/api
migration:
	@read -p "Enter migration name: " MIGRATION; \
	docker run -v ./migrations:/migrations migrate/migrate create -seq -ext=.sql -dir=/migrations $$MIGRATION
migrate:
	docker run -v ./migrations:/migrations --network host migrate/migrate -path=/migrations/ -database $(PG_DSN) up

PG_DSN:=postgres://head_agent:pa55word@localhost:5432/sca?sslmode=disable
