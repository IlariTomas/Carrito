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
	# 2. Usamos la ruta completa para ejecutar sqlc
	$(GOBIN)/sqlc generate

templ:    
	@echo "Generando vistas con templ..."
	# 3. Usamos la ruta completa para ejecutar templ
	$(GOBIN)/templ generate 

## Construye y levanta los contenedores (api, db)
# Agregamos 'setup' aquí si quieres asegurar que se instalen, 
# pero idealmente 'setup' se corre una sola vez manualmente.
build: sqlc templ
	@echo "Construyendo y levantando contenedores (api, db)..."
	docker compose up --build db api

## Ejecuta los tests de Hurl (levanta 'tester' y depende de 'api')
test:
	@echo "Ejecutando tests de Hurl..."
	docker compose up --build tester

## Detiene los contenedores (api, db)
stop:
	@echo "Deteniendo contenedores..."
	docker compose stop

clean:
	@echo "Limpiando contenedores y volúmenes..."
	docker compose down -v

## Alias
up: build
down: stop
full-reset: clean build