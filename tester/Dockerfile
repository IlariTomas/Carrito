FROM curlimages/curl:latest

USER root

#    Instalamos 'hurl' Y reinstalamos 'curl'
RUN apk add --no-cache hurl curl --repository=https://dl-cdn.alpinelinux.org/alpine/edge/testing

WORKDIR /app

COPY requests.hurl .
COPY cargar_productos.sh .

RUN chmod +x ./cargar_productos.sh

USER curl_user