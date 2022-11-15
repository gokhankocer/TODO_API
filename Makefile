SHELL = /bin/sh
TODO_API = todo_api/api.yaml

tools/swagger:
	$(call print-target)
	GOBIN=$(CURDIR)/tools go install github.com/go-swagger/go-swagger/cmd/swagger@v0.27.0

.PHONY: install_tools
install_tools: tools/swagger tools/golang-migrate

.PHONY: models
models: tools/swagger
	$(call print-target)
	find ./models -type f -not -name '*_test.go' -delete
	./tools/swagger  generate model DEBUG=1 --spec=docs/$(TODO_API)

tools/golang-migrate:
	$(call print-target)
	GOBIN=$(CURDIR)/tools go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.14.1

.PHONY: migrate_up
migrate_up: tools/migrate
	$(call print-target)
	./tools/migrate -path migrations/retail_api -database ${MIGRATION_DB_URL} -verbose up ${step}

.PHONY: migrate_down
migrate_down: tools/migrate
	./tools/migrate -path migrations/retail_api -database ${MIGRATION_DB_URL} -verbose down ${step}

.PHONY: migration_file 
migration_file: tools/migrate
	./tools/migrate create -ext sql -dir migrations/$(app) -seq $(title)

define print-target
	@printf "Executing target: \033[36m$@\033[0m\n"
endef
