# 1. Definimos dinámicamente dónde guarda Go los ejecutables
GOBIN := $(shell go env GOPATH)/bin

## Instala las herramientas de desarrollo local necesarias
setup:
	@echo "Instalando dependencias de Go (sqlc)..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

	@echo "Instalando dependencias de Go (templ)..."
	go install github.com/a-h/templ/cmd/templ@latest
	
sqlc:
	@echo "Generando código Go con sqlc..."
	$(GOBIN)/sqlc generate

templ:    
	@echo "Generando vistas con templ..."
	$(GOBIN)/templ generate 

## Construye y levanta los contenedores (api, db)
build: sqlc templ
	@echo "Construyendo y levantando contenedores (api, db)..."
	docker compose up --build db api

## Detiene los contenedores (api, db)
stop:
	@echo "Deteniendo contenedores..."
	docker compose stop

clean:
	@echo "Limpiando contenedores y volúmenes (reseteando la base de datos)..."
	docker compose down -v

## Alias
up: build
down: stop
full-reset: clean build