    # Etapa 1: Compilación del binario
    FROM golang:1.23-alpine AS builder

    # 🔹 Configura el directorio de trabajo a /api (como pediste)
    WORKDIR /api

    # 🔹 Copia go.mod y go.sum desde la raíz del contexto
    COPY go.mod go.sum ./
    RUN go mod download

    # 🔹 Copia el código fuente y estáticos desde la raíz del contexto
    COPY main.go .
    COPY index.html .
    COPY static ./static
    COPY db ./db

    # Compila el binario
    RUN go build -o carrito .

    # Etapa 2: Imagen final (más liviana)
    FROM alpine:latest
    RUN apk --no-cache add ca-certificates

    # 🔹 Crea el directorio /api
    WORKDIR /api

    # 🔹 Copia los archivos compilados desde la etapa anterior
    COPY --from=builder /api/carrito .
    COPY --from=builder /api/index.html .
    COPY --from=builder /api/static ./static

    # Expone el puerto 8080
    EXPOSE 8080

    # Comando por defecto
    CMD ["./carrito"]