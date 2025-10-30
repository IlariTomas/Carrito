#!/bin/sh

# La URL de tu API
API_URL="http://localhost:8080/products"

# Producto 1: Pc de escritorio
curl -X POST -H "Content-Type: application/json" \
-d '{
    "nombre_producto": "Pc de escritorio",
    "descripcion": "Descripción del producto.",
    "precio": "10000.00",
    "stock": 100,
    "categoria": "Computadoras",
    "imagen": "https://www.crucial.mx/content/dam/crucial/articles/for-pc-builders/new025-how-to-upgrade-your-pc/modern-gaming-pc.jpg.transform/medium-jpg/img.jpg"
}' "$API_URL"

# Producto 2: Laptop Gamer
curl -X POST -H "Content-Type: application/json" \
-d '{
    "nombre_producto": "Laptop Gamer",
    "descripcion": "Laptop con alto rendimiento para juegos.",
    "precio": "200000.00",
    "stock": 50,
    "categoria": "Computadoras",
    "imagen": "https://m.media-amazon.com/images/I/811QpiYXe-L.jpg"
}' "$API_URL"

# Producto 3: Teclado Mecánico
curl -X POST -H "Content-Type: application/json" \
-d '{
    "nombre_producto": "Teclado Mecánico",
    "descripcion": "Teclado con switches mecánicos y retroiluminación.",
    "precio": "15000.00",
    "stock": 200,
    "categoria": "Perifericos",
    "imagen": "https://http2.mlstatic.com/D_960056-MLA95235561941_102025-C.jpg"
}' "$API_URL"

# Producto 4: Camara Logitech
curl -X POST -H "Content-Type: application/json" \
-d '{
    "nombre_producto": "Cámara web Logitech Brio 4K 90FPS color negro",
    "descripcion": "Cámara web HD para videoconferencias.",
    "precio": "15000.00",
    "stock": 120,
    "categoria": "Perifericos",
    "imagen": "https://http2.mlstatic.com/D_NQ_NP_2X_682671-MLA95663048448_102025-F.webp.jpg"
}' "$API_URL"


# Producto 5: Mouse Gamer
curl -X POST -H "Content-Type: application/json" \
-d '{
    "nombre_producto": "Mouse Gamer Logitech G203",
    "descripcion": "Mouse ergonómico para gamers.",
    "precio": "5000.00",
    "stock": 100,
    "categoria": "Perifericos",
    "imagen": "https://http2.mlstatic.com/D_NQ_NP_2X_849696-MLA95939215137_102025-F.webp.jpg"
}' "$API_URL"

echo "¡Productos cargados!"