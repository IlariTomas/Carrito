#   Etapa única: Configuración del entorno de pruebas con 'hurl'
FROM curlimages/curl:latest

USER root

#    Instalamos 'hurl' Y reinstalamos 'curl'
RUN apk add --no-cache hurl curl --repository=https://dl-cdn.alpinelinux.org/alpine/edge/testing

#   Configura el directorio de trabajo a /app
WORKDIR /app

#   Copia el código fuente y estáticos desde la raíz del contexto
COPY requests.hurl .
COPY cargar_productos.sh .

#   Corre script para cargar productos
RUN chmod +x ./cargar_productos.sh

#   Cambiamos a un usuario no root
USER curl_user