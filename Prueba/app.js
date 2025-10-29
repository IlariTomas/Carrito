// La constante API_HOST se define en index.html
const PRODUCT_URL = `${API_HOST}/products`;
const SALES_URL = `${API_HOST}/sales`;
// NUEVA URL para usuarios
const USERS_URL = `${API_HOST}/users`;

const productsList = document.getElementById('products-list');
const salesList = document.getElementById('sales-list');
// NUEVO elemento de lista de usuarios
const usersList = document.getElementById('users-list'); 

const productForm = document.getElementById('create-product-form');
const saleForm = document.getElementById('create-sale-form');
// NUEVO formulario de usuario
const userForm = document.getElementById('create-user-form'); 

const productMessage = document.getElementById('product-message');
const saleMessage = document.getElementById('sale-message');
// NUEVO elemento de mensaje de usuario
const userMessage = document.getElementById('user-message'); 

/**
 * Muestra un mensaje de estado (éxito o error) en la interfaz.
 * @param {HTMLElement} element - El elemento DIV donde se mostrará el mensaje.
 * @param {string} msg - El texto del mensaje.
 * @param {boolean} isSuccess - Si es true, es éxito (verde); si es false, es error (rojo).
 */
function displayMessage(element, msg, isSuccess) {
    element.textContent = msg;
    element.className = 'message ' + (isSuccess ? 'success' : 'error');
    setTimeout(() => {
        element.textContent = '';
        element.className = 'message';
    }, 5000); // El mensaje desaparece después de 5 segundos
}

// ------------------------------------------------------------------
// LÓGICA DE LISTADO Y RENDERIZACIÓN
// ------------------------------------------------------------------

/**
 * Función genérica para obtener datos de la API.
 * Se ha mejorado para asegurar que siempre devuelve un Array, incluso si la API responde con null.
 * @param {string} url - URL del endpoint.
 * @returns {Promise<Array>} - Array de entidades (o array vacío si falla/es nulo).
 */
async function fetchEntities(url) {
    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Error HTTP: ${response.status}`);
        }
        
        const data = await response.json();
        
        // CORRECCIÓN: Si el backend retorna null o undefined, regresamos un array vacío.
        if (!data || typeof data !== 'object') {
            return [];
        }
        
        // Si el backend retorna un objeto JSON (como {}) en lugar de un array, 
        // asumimos que es un array vacío para evitar errores de .forEach.
        if (!Array.isArray(data)) {
            console.warn(`La API en ${url} devolvió un objeto, se tratará como un array vacío.`);
            return [];
        }

        return data;

    } catch (error) {
        console.error(`Error al obtener datos de ${url}:`, error);
        return []; // Siempre devolvemos un array vacío en caso de error
    }
}

/**
 * Función genérica para eliminar una entidad.
 * @param {string} url - URL completa de la entidad a eliminar (e.g., /products/1).
 */
async function deleteEntity(url, successCallback, errorElement) {
    try {
        const response = await fetch(url, {
            method: 'DELETE',
        });

        if (response.status === 204) {
            displayMessage(errorElement, 'Eliminación exitosa.', true);
            successCallback(); // Refresca la lista
        } else {
            // Se lee el cuerpo para reportar el error específico de la API (si no es 404)
            const errorText = await response.text();
            throw new Error(`Fallo al eliminar: ${errorText || response.statusText}`);
        }
    } catch (error) {
        console.error('Error al eliminar:', error);
        // Si la respuesta es HTML (el 404 del servidor Go), mostramos un mensaje genérico.
        const errorMessage = error.message.includes('404') ? 
                             'Error 404: La ruta de la API no se encontró. Verifique el servidor Go.' : 
                             `Error: ${error.message}`;
        displayMessage(errorElement, errorMessage, false);
    }
}


/**
 * Renderiza la lista de usuarios en el DOM.
 * @param {Array} users - Array de usuarios.
 */
function renderUsers(users) {
    const userListArray = Array.isArray(users) ? users : [];

    usersList.innerHTML = '';
    if (userListArray.length === 0) {
        usersList.innerHTML = '<p>No hay usuarios registrados.</p>';
        return;
    }

    userListArray.forEach(u => {
        const item = document.createElement('div');
        item.className = 'entity-item';
        
        // Asumiendo que el GET devuelve las claves en minúscula (id_usuario, nombre, email, rol)
        item.innerHTML = `
            <div class="entity-info">
                <strong>ID: ${u.id_usuario || 'N/A'}</strong> - ${u.nombre} 
                <br>
                Email: ${u.email} | Rol: ${u.rol}
            </div>
            <button class="delete-btn" data-id="${u.id_usuario}" data-type="user">Eliminar</button>
        `;
        usersList.appendChild(item);
    });
}

/**
 * Renderiza la lista de productos en el DOM.
 * @param {Array} products - Array de productos.
 */
function renderProducts(products) {
    // CORRECCIÓN: Si products no es un array (por ejemplo, es null), lo inicializamos como vacío.
    const productListArray = Array.isArray(products) ? products : [];

    productsList.innerHTML = '';
    if (productListArray.length === 0) {
        productsList.innerHTML = '<p>No hay productos registrados.</p>';
        return;
    }

    productListArray.forEach(p => {
        const item = document.createElement('div');
        item.className = 'entity-item';
        
        item.innerHTML = `
            <div class="entity-info">
                <strong>ID: ${p.id_producto || 'N/A'}</strong> - ${p.nombre_producto} 
                (Categoría: ${p.categoria}) <br>
                Precio: $${p.precio} | Stock: ${p.stock}
            </div>
            <button class="delete-btn" data-id="${p.id_producto}" data-type="product">Eliminar</button>
        `;
        productsList.appendChild(item);
    });
}

/**
 * Renderiza la lista de ventas en el DOM.
 * @param {Array} sales - Array de ventas.
 */
function renderSales(sales) {
    // CORRECCIÓN: Si sales no es un array (por ejemplo, es null), lo inicializamos como vacío.
    const salesListArray = Array.isArray(sales) ? sales : [];

    salesList.innerHTML = '';
    if (salesListArray.length === 0) {
        salesList.innerHTML = '<p>No hay ventas registradas.</p>';
        return;
    }

    salesListArray.forEach(s => {
        const item = document.createElement('div');
        item.className = 'entity-item';
        
        // Manejo de fecha
        const date = s.fecha ? new Date(s.fecha).toLocaleDateString() : 'Fecha N/A';

        item.innerHTML = `
            <div class="entity-info">
                <strong>Venta #${s.id_venta || 'N/A'}</strong> - Producto: ${s.id_producto} / Usuario: ${s.id_usuario} <br>
                Cantidad: ${s.cantidad} | Total: $${s.total} | Fecha: ${date}
            </div>
            <button class="delete-btn" data-id="${s.id_venta}" data-type="sale">Eliminar</button>
        `;
        salesList.appendChild(item);
    });
}

/**
 * Carga y renderiza todas las listas (usuarios, productos y ventas).
 */
async function loadAndRenderEntities() {
    usersList.innerHTML = '<p>Cargando usuarios...</p>';
    productsList.innerHTML = '<p>Cargando productos...</p>';
    salesList.innerHTML = '<p>Cargando ventas...</p>';
    
    // Ahora incluye la carga de usuarios
    const [users, products, sales] = await Promise.all([
        fetchEntities(USERS_URL),
        fetchEntities(PRODUCT_URL),
        fetchEntities(SALES_URL)
    ]);

    renderUsers(users);
    renderProducts(products);
    renderSales(sales);
}

// ------------------------------------------------------------------
// LÓGICA DE CREACIÓN (POST)
// ------------------------------------------------------------------

/**
 * Maneja el envío de formularios para crear productos o ventas.
 * @param {Event} e - Evento de envío del formulario.
 * @param {string} url - Endpoint de la API (PRODUCT_URL o SALES_URL).
 * @param {object} data - Objeto con los datos a enviar.
 * @param {HTMLElement} msgElement - Elemento para mostrar el resultado.
 */
async function handleCreation(e, url, data, msgElement) {
    e.preventDefault();
    
    // Conversión de tipos para el backend (precio/total)
    if (data.precio !== undefined) data.precio = String(data.precio);
    if (data.total !== undefined) data.total = String(data.total);

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });

        if (response.status === 201) {
            displayMessage(msgElement, 'Creación exitosa!', true);
            e.target.reset(); // Limpia el formulario
            await loadAndRenderEntities(); // Refresca la lista
        } else {
            const errorText = await response.text();
            throw new Error(`Fallo al crear: ${errorText}`);
        }
    } catch (error) {
        console.error('Error de creación:', error);
        displayMessage(msgElement, `Error: ${error.message}`, false);
    }
}

// ------------------------------------------------------------------
// CONFIGURACIÓN DE EVENT LISTENERS
// ------------------------------------------------------------------

// 1. Evento para crear Usuario
userForm.addEventListener('submit', (e) => {
    // CORRECCIÓN CLAVE: Se usan mayúsculas iniciales (Nombre, Email, Rol) en las claves JSON.
    // Esto es necesario porque el backend de Go generalmente requiere que los campos de
    // los structs (que son a los que se mapea el JSON) empiecen con mayúscula.
    const data = {
        Nombre: document.getElementById('u-nombre').value,
        Email: document.getElementById('u-email').value,
    };
    handleCreation(e, USERS_URL, data, userMessage);
});


// 2. Evento para crear Producto
productForm.addEventListener('submit', (e) => {
    const data = {
        nombre_producto: document.getElementById('p-nombre').value,
        descripcion: document.getElementById('p-descripcion').value,
        precio: parseFloat(document.getElementById('p-precio').value).toFixed(2),
        stock: parseInt(document.getElementById('p-stock').value),
        categoria: document.getElementById('p-categoria').value,
    };
    handleCreation(e, PRODUCT_URL, data, productMessage);
});


// 3. Evento para crear Venta
saleForm.addEventListener('submit', (e) => {
    
    // Convertir la fecha local a formato UTC (RFC3339) para la API de Go
    const localDateInput = document.getElementById('s-fecha').value;
    let isoDate = null;
    if (localDateInput) {
        // Convierte el valor local (e.g., "2025-10-25T15:00") a un objeto Date y luego a formato ISO (UTC)
        isoDate = new Date(localDateInput).toISOString();
    }
    
    const data = {
        id_producto: parseInt(document.getElementById('s-producto-id').value),
        id_usuario: parseInt(document.getElementById('s-usuario-id').value),
        cantidad: parseInt(document.getElementById('s-cantidad').value),
        total: parseFloat(document.getElementById('s-total').value).toFixed(2),
        // Si isoDate es null, la API de Go deberá manejarlo como nulo si el campo sql.NullTime lo permite.
        // Si no, podríamos forzar a que sea una fecha actual si no se provee.
        fecha: isoDate || new Date().toISOString()
    };
    handleCreation(e, SALES_URL, data, saleMessage);
});


// 4. Evento para los botones de Eliminar (Delegación de eventos)
document.getElementById('list-entities').addEventListener('click', (e) => {
    if (e.target.classList.contains('delete-btn')) {
        const id = e.target.getAttribute('data-id');
        const type = e.target.getAttribute('data-type');
        let url;
        let msgElement;

        // Lógica de eliminación corregida para usar rutas singulares (API_HOST/entity/ID)
        if (type === 'user') {
            url = `${API_HOST}/user/${id}`; // Coincide con la ruta /user/ en main.go
            msgElement = userMessage;
        } else if (type === 'product') {
            url = `${API_HOST}/product/${id}`; // Coincide con la ruta /product/ en main.go
            msgElement = productMessage;
        } else if (type === 'sale') {
            url = `${API_HOST}/sale/${id}`; // Coincide con la ruta /sale/ en main.go
            msgElement = saleMessage;
        }

        if (url) {
            if (confirm(`¿Está seguro de eliminar ${type} con ID ${id}?`)) {
                deleteEntity(url, loadAndRenderEntities, msgElement);
            }
        }
    }
});


// Inicialización: Cargar entidades al iniciar la página
window.onload = loadAndRenderEntities;
