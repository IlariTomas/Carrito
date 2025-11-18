## Instala las herramientas de desarrollo local necesarias
setup:
	@echo "Instalando dependencias de Go (sqlc)..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

	@echo "Instalando dependencias de Go (templ)..."
	go install github.com/valyala/templ/cmd/templ@latest
	
sqlc:
	@echo "Generando código Go con sqlc..."
	sqlc generate

templ:    
	@echo "Generando vistas con templ..."
	templ generate 

## Construye y levanta los contenedores (api, db)
build: sqlc templ
	@echo "Construyendo y levantando contenedores (api, db)..."
	docker compose up --build  db api

## Ejecuta los tests de Hurl (levanta 'tester' y depende de 'api')
test:
	@echo "Ejecutando tests de Hurl..."
	docker compose up --build tester

## Detiene los contenedores (api, db)
stop:
	@echo "Deteniendo contenedores..."
	docker compose stop

clean:
	@echo "Limpiando contenedores y volúmenes (reseteando la base de datos)..."
	docker compose down -v

## Alias para 'make build'
up: build

## Alias para 'make stop'
down: stop

## Ejecuta el ciclo completo: Limpia, construye y prueba.
full-reset: clean build 
