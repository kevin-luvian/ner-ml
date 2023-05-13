export NOW=$(shell date +"%Y/%m/%d %H:%M:%S")
export PKGS=$(shell go list ./... | grep -vE '(vendor|cmd|entity|pkg/assert)')

configure:
	@echo "${NOW} === CONFIGURING FILES ==="
	@cp ./conf/app.ini.example conf/app.ini
	@echo "=== CONFIGURED ==="

generate:
	@echo "${NOW} === GENERATING FILES ==="
	@go generate ./...
	@echo "=== GENERATED ==="

migrate:
	@echo "${NOW} === RUNNING MIGRATION ==="
	@go run ./server/cmd/migrator $(args)
	@echo "=== MIGRATED ==="

play:
	@echo "${NOW} 🚀 === RUNNING PLAYGROUND === 🚀"
	@go run ./server/cmd/playground
	@echo "=== DONE ==="

.PHONY: dev
dev:
	@echo "${NOW} === RUNNING DEVELOPMENT ENV ==="
	@docker-compose stop goauth-be goauth-fe && docker-compose up -d goauth-be goauth-fe
	@echo "click this link to open the backend http://localhost:8000"
	@echo "click this link to open the frontend http://localhost:8001"

dev-all:
	@echo "${NOW} 🛠 === RUNNING DEVELOPMENT ALL === 🛠"
	@cd tools/ && docker-compose stop && docker-compose up -d
	@docker-compose stop && docker-compose up -d
	@echo "🚀 === RAN === 🚀"

dev-db:
	@echo "🛠 ${NOW} === RUNNING DEVELOPMENT DB === 🛠"
	@docker-compose stop goauth-pg && docker-compose up -d goauth-pg
	@echo "🚀 === RAN === 🚀"

dev-fe:
	@echo "${NOW} === RUNNING DEVELOPMENT ENV ==="
	@docker-compose stop goauth-fe && docker-compose up -d goauth-fe
	@echo "click this link to open the page http://localhost:8001"

dev-tools:
	@echo "${NOW} === RUNNING DEVELOPMENT TOOLS ==="
	@cd tools/ && docker-compose stop && docker-compose up -d
	@echo "click this link to open yacht page http://localhost:5000"
	@echo "click this link to open prometheus http://localhost:5001"

clean:
	@echo "🛠 CLEANING MACHINE FOR DEVELOPMENT 🛠"
	@echo "1 REMOVING BIN FOLDER"
	@rm -f -r ./bin
	@echo "1 REMOVING TEST OUT"
	@rm -f ./test.out
	@echo "🚀 Done, You are ready to Go 🚀"

down:
	@docker-compose stop
	@cd tools/ && docker-compose stop
	@docker-compose down
	@cd tools/ && docker-compose down

down-tools:
	@cd tools/ && docker-compose stop
	@cd tools/ && docker-compose down

test:
	@echo "${NOW} === TESTING ==="
	@go test -cover -race ${PKGS} -short | tee ./test.out
	@echo "=== DONE ==="