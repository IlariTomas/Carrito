// La constante API_HOST se define en el HTML
const USER_URL = `${API_HOST}/users`;

const usersList = document.getElementById('users-list');
const userForm = document.getElementById('user-insert-section');
const userMessage = document.getElementById('user-message');

function displayMessage(element, msg, isSuccess) {
    element.textContent = msg;
    element.className = 'message ' + (isSuccess ? 'success' : 'error');
}

async function handleCreation(e, url, data, msgElement) {
    e.preventDefault();

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        });

        if (response.status === 201) {
            displayMessage(msgElement, 'usuario creado correctamente!', true);
            e.target.reset();
            await loadAndRenderEntities();
        } else {
            const errorText = await response.text();
            throw new Error(`Error al crear usuario: ${errorText}`);
        }
    } catch (error) {
        console.error('Error de creación:', error);
        displayMessage(msgElement, `Error al crear usuario, verifica los campos`, false);

    }
}

userForm.addEventListener('submit', (e) => {
    const data = {
        nombre_usuario: document.getElementById('userName').value,
        email: document.getElementById('userEmail').value,
    };
    handleCreation(e, USER_URL, data, userMessage);
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
    usersList.innerHTML = '<p>Cargando usuarios...</p>';
    const users = await fetchEntities(USER_URL);
    renderUsers(users);
}

function renderUsers(users) {
    const list = Array.isArray(users) ? users : [];
    usersList.innerHTML = '';
    if (list.length === 0) {
        usersList.innerHTML = '<p>No hay usuarios registrados.</p>';
        return;
    }

    list.forEach(u => {
        const item = document.createElement('div');
        item.className = 'entity-item';
        item.innerHTML = `
            <div class="entity-info">
                <div class="key-data">
                    <strong>ID: ${u.id_usuario}-</strong>
                    <strong>${u.nombre_usuario}</strong>
                </div>
                <div class="secondary-data">
                    <p>Email: ${u.email}</p>
                </div>
            </div>
            <button class="delete-btn" data-id="${u.id_usuario}">Eliminar</button>
        `;
        usersList.appendChild(item);
    });
}

document.addEventListener('DOMContentLoaded', loadAndRenderEntities);
