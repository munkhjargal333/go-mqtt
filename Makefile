airapp: 
	@air -c ./.air.app.toml

runapp: buildapp
	@./bin/app/main

runadmin: buildadmin
	@./bin/admin/main

buildapp:
	@go build -o ./bin/app/main ./cmd/app/main.go

buildadmin:
	@go build -o ./bin/admin/main ./cmd/admin/main.go

migrate:
	@go run ./cmd/migrate/main.go --migrate

drop:
	@go run ./cmd/migrate/main.go --drop
	
seed:
	@go run ./cmd/seed/main.go
