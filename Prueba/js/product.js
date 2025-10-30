// La constante API_HOST se define en el HTML
const PRODUCT_URL = `${API_HOST}/products`;

const productsList = document.getElementById('products-list');
const productForm = document.getElementById('create-product-form');
const productMessage = document.getElementById('product-message');

function displayMessage(element, msg, isSuccess) {
    element.textContent = msg;
    element.className = 'message ' + (isSuccess ? 'success' : 'error');
    setTimeout(() => {
        element.textContent = '';
        element.className = 'message';
    }, 5000);
}

async function handleCreation(e, url, data, msgElement) {
    e.preventDefault();

    if (data.precio !== undefined) data.precio = String(data.precio);
    if (data.total !== undefined) data.total = String(data.total);

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        });

        if (response.status === 201) {
            displayMessage(msgElement, 'Producto creado correctamente!', true);
            e.target.reset();
            await loadAndRenderEntities();
        } else {
            const errorText = await response.text();
            throw new Error(`Error al crear producto: ${errorText}`);
        }
    } catch (error) {
        console.error('Error de creación:', error);
        displayMessage(msgElement, `Error al crear producto, verifica los campos`, false);
    }
}

productForm.addEventListener('submit', (e) => {
    const data = {
        nombre_producto: document.getElementById('p-nombre').value,
        descripcion: document.getElementById('p-descripcion').value,
        precio: parseFloat(document.getElementById('p-precio').value).toFixed(2),
        stock: parseInt(document.getElementById('p-stock').value),
        categoria: document.getElementById('p-categoria').value,
        imagen: document.getElementById('p-imagen').value,
    };
    handleCreation(e, PRODUCT_URL, data, productMessage);
});

async function fetchEntities(url) {
    try {
        const response = await fetch(url);
        if (!response.ok) throw new Error(`Error HTTP: ${response.status}`);
        const data = await response.json();
        return Array.isArray(data) ? data : [];
    } catch (error) {
        console.error(`Error al obtener datos de ${url}:`, error);
        return [];
    }
}

async function loadAndRenderEntities() {
    productsList.innerHTML = '<p>Cargando productos...</p>';
    const products = await fetchEntities(PRODUCT_URL);
    renderProducts(products);
}

function renderProducts(products) {
    const list = Array.isArray(products) ? products : [];
    productsList.innerHTML = '';
    if (list.length === 0) {
        productsList.innerHTML = '<p>No hay productos registrados.</p>';
        return;
    }

    list.forEach(p => {
        const item = document.createElement('div');
        item.className = 'entity-item';
        item.innerHTML = `
            <div class="entity-info">
                <div class="key-data">
                    <strong>ID: ${p.id_producto}-</strong>
                    <strong>${p.nombre_producto}</strong>
                </div>
                <div class="secondary-data">
                    <p>Precio: $${p.precio}</p>
                    <p>Stock: ${p.stock}</p>
                </div>
            </div>
            <button class="delete-btn" data-id="${p.id_producto}">Eliminar</button>
        `;
        productsList.appendChild(item);
    });
}

document.addEventListener('DOMContentLoaded', loadAndRenderEntities);
