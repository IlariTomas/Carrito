// La constante API_HOST se define en el HTML
const SALES_URL = `${API_HOST}/sales`;
const SALE_URL = `${API_HOST}/sale`;
const PRODUCT_URL = `${API_HOST}/products`;
const USER_URL = `${API_HOST}/users`;

const salesList = document.getElementById('sales-list');
const saleForm = document.getElementById('create-sale-form');
const saleMessage = document.getElementById('sale-message');

function displayMessage(element, msg, isSuccess) {
    element.textContent = msg;
    element.className = 'message ' + (isSuccess ? 'success' : 'error');
}

async function handleCreation(e, url, data, msgElement) {
    e.preventDefault();

    if (data.cantidad !== undefined) data.cantidad = parseInt(data.cantidad);
    if (data.id_producto !== undefined) data.id_producto = parseInt(data.id_producto);
    if (data.id_usuario !== undefined) data.id_usuario = parseInt(data.id_usuario);
    if (data.fecha !== undefined) data.fecha = new Date(data.fecha);
    
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
            throw new Error(`Error al Generar venta: ${errorText}`);
        }
    } catch (error) {
        console.error('Error de creación:', error);
        displayMessage(msgElement, `Error al Generar venta, verifica los campos`, false);
    }
}

saleForm.addEventListener('submit', (e) => {
    const data = {
        id_producto: document.getElementById('select-product').value,
        id_usuario: document.getElementById('buyer-Select').value,
        cantidad: document.getElementById('productStock').value,
        fecha: document.getElementById('purchaseDate').value,
        total: document.getElementById('totalAmount').textContent.replace('$', ''),
    };
    handleCreation(e, SALES_URL, data, saleMessage);
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
    salesList.innerHTML = '<p>Cargando ventas...</p>';
    const sales = await fetchEntities(SALES_URL);
    renderSales(sales);
}

async function deleteSale(saleId) {
    try {
        const response = await fetch(`${SALE_URL}/${saleId}`, {
            method: 'DELETE',
        });

        if (response.ok) {
            await loadAndRenderEntities();
        } else {
            console.error(`Error al eliminar venta con ID ${saleId}:`, await response.text());
        }
    } catch (error) {
        console.error('Error de eliminación:', error);
    }
}

function renderSales(sales) {
    const list = Array.isArray(sales) ? sales : [];
    salesList.innerHTML = '';
    if (list.length === 0) {
        salesList.innerHTML = '<p>No hay ventas registradas.</p>';
        return;
    }
    
    list.forEach(v => {
        const item = document.createElement('div');
        item.className = 'entity-item';
        item.innerHTML = `
            <div class="entity-info">
                <div class="key-data">
                    <strong>ID: ${v.id_venta}-</strong>
                    <strong>Producto: ${v.id_producto}</strong>
                </div>
                <div class="secondary-data">
                    <p>id_usuario: ${v.id_usuario}</p>
                    <p>total: ${v.total}</p>
                </div>
            </div>
            <button class="delete-btn">Eliminar</button>
        `;
        const deleteButton = item.querySelector('.delete-btn');
        deleteButton.addEventListener('click', () => {
            deleteSale(v.id_venta);
        });

        salesList.appendChild(item);
    });
}

document.addEventListener('DOMContentLoaded', loadAndRenderEntities);
// Actualizar total al cambiar cantidad o producto
const productSelect = document.getElementById('select-product');
const quantityInput = document.getElementById('productStock');
const totalAmount = document.getElementById('totalAmount');

function updateTotal() {
    const productId = productSelect.value;
    const quantity = parseInt(quantityInput.value) || 0;
    const price = PrecioProduct.get(productId) || 0;
    const total = price * quantity;
    totalAmount.textContent = total.toFixed(2);
}

// Escuchamos cambios en producto y cantidad
productSelect.addEventListener('change', updateTotal);
quantityInput.addEventListener('input', updateTotal);

// Cargar productos y compradores para los selectores
async function loadAndRenderSelectors() {
    const products = await fetchEntities(PRODUCT_URL);
    const users = await fetchEntities(USER_URL);
    renderProductOptions(products);
    renderUserOptions(users);
}

const PrecioProduct = new Map();
function renderProductOptions(products) {
    const select = document.getElementById('select-product');
    
    products.forEach(p => {
        const option = document.createElement('option');
        option.value = p.id_producto;
        option.textContent = `${p.nombre_producto}`;
        select.appendChild(option);
        PrecioProduct.set(String(p.id_producto), Number(p.precio));
    });
}

function renderUserOptions(users) {
    const select = document.getElementById('buyer-Select');
    // Crear un Map {id -> precio} y exponerlo globalmente para usarlo desde otros módulos
    
    users.forEach(u => {
        const option = document.createElement('option');
        option.value = u.id_usuario;
        option.textContent = u.nombre_usuario;
        select.appendChild(option);
    });
}

document.addEventListener('DOMContentLoaded', loadAndRenderSelectors);
