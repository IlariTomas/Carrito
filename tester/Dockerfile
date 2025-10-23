FROM curlimages/curl:latest

# Cambiamos al usuario root para poder instalar paquetes
USER root

# 🔹 ARREGLO: Instala Hurl usando el gestor de paquetes de Alpine.
# Esto es mucho más confiable que descargar el binario con curl.
RUN apk add --no-cache hurl --repository=https://dl-cdn.alpinelinux.org/alpine/edge/testing

# Volvemos al usuario no privilegiado por defecto
USER curl_user

WORKDIR /app

# Copia tu archivo de requests a la carpeta actual (/app) del contenedor
COPY requests.hurl .