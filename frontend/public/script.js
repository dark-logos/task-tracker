let accessToken = localStorage.getItem('accessToken') || '';
let refreshToken = localStorage.getItem('refreshToken') || '';

function saveTokens(access, refresh) {
    accessToken = access;
    refreshToken = refresh;
    localStorage.setItem('accessToken', access);
    localStorage.setItem('refreshToken', refresh);
}

function setMessage(elementId, message, isError = false) {
    const element = document.getElementById(elementId);
    element.textContent = message;
    element.className = isError ? 'text-red-500' : 'text-green-500';
}

function formatDateTime(dateTime) {
    if (!dateTime) return null;
    const date = new Date(dateTime);
    return date.toISOString();
}

async function register() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const email = prompt('Enter email (optional):');
    try {
        const response = await fetch('http://localhost:8080/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password_hash: password, email })
        });
        const data = await response.json();
        setMessage('auth-message', data.message || data.error, !response.ok);
    } catch (error) {
        setMessage('auth-message', 'Error: ' + error.message, true);
    }
}

async function login() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    try {
        const response = await fetch('http://localhost:8080/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });
        const data = await response.json();
        if (response.ok) {
            saveTokens(data.access_token, data.refresh_token);
            document.getElementById('auth').style.display = 'none';
            document.getElementById('tasks').style.display = 'block';
            loadTasks();
        } else {
            setMessage('auth-message', data.error, true);
        }
    } catch (error) {
        setMessage('auth-message', 'Error: ' + error.message, true);
    }
}

async function refreshTokenIfNeeded() {
    if (!refreshToken) return false;
    try {
        const response = await fetch('http://localhost:8080/refresh', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refresh_token: refreshToken })
        });
        const data = await response.json();
        if (response.ok) {
            saveTokens(data.access_token, data.refresh_token);
            return true;
        }
        return false;
    } catch (error) {
        console.error('Token refresh failed:', error);
        return false;
    }
}

async function loadTasks() {
    if (!accessToken) {
        setMessage('task-message', 'Please login first', true);
        return;
    }
    try {
        const response = await fetch('http://localhost:8080/tasks', {
            headers: { 'Authorization': `Bearer ${accessToken}` }
        });
        if (response.status === 401) {
            if (await refreshTokenIfNeeded()) {
                return loadTasks();
            }
        }
        const data = await response.json();
        const pendingTasks = document.getElementById('pending-tasks');
        const inProgressTasks = document.getElementById('in_progress-tasks');
        const doneTasks = document.getElementById('done-tasks');
        pendingTasks.innerHTML = '';
        inProgressTasks.innerHTML = '';
        doneTasks.innerHTML = '';
        if (response.ok) {
            let tasks = Array.isArray(data) ? data : [];
            tasks.sort((a, b) => b.priority - a.priority);
            tasks.forEach(task => {
                const div = document.createElement('div');
                div.className = 'task border p-4 rounded mb-2';
                const title = document.createElement('h3');
                title.className = 'text-lg font-semibold';
                title.textContent = task.title;
                const desc = document.createElement('p');
                desc.textContent = task.description || 'No description';
                const status = document.createElement('p');
                status.textContent = `Status: ${task.status}`;
                const priority = document.createElement('p');
                priority.textContent = `Priority: ${task.priority}`;
                const due = document.createElement('p');
                due.textContent = `Due: ${task.due_date || 'N/A'}`;
                const updateBtn = document.createElement('button');
                updateBtn.textContent = 'Update';
                updateBtn.className = 'bg-yellow-500 text-white p-1 rounded hover:bg-yellow-600';
                updateBtn.onclick = () => openModal(task.id, task.title, task.description || '', task.status, task.priority, task.due_date || '');
                const deleteBtn = document.createElement('button');
                deleteBtn.textContent = 'Delete';
                deleteBtn.className = 'bg-red-500 text-white p-1 ml-2 rounded hover:bg-red-600';
                deleteBtn.onclick = () => deleteTask(task.id);
                div.appendChild(title);
                div.appendChild(desc);
                div.appendChild(status);
                div.appendChild(priority);
                div.appendChild(due);
                div.appendChild(updateBtn);
                div.appendChild(deleteBtn);
                if (task.status === 'pending') {
                    pendingTasks.appendChild(div);
                } else if (task.status === 'in_progress') {
                    inProgressTasks.appendChild(div);
                } else if (task.status === 'done') {
                    doneTasks.appendChild(div);
                }
            });
        } else {
            setMessage('task-message', data.error || 'Failed to fetch tasks', true);
        }
    } catch (error) {
        setMessage('task-message', 'Error: ' + error.message, true);
    }
}

async function createTask() {
    if (!accessToken) {
        setMessage('task-message', 'Please login first', true);
        return;
    }
    let priority = parseInt(document.getElementById('priority').value);
    if (isNaN(priority)) {
        priority = 1; 
    }
    const task = {
        title: document.getElementById('title').value,
        description: document.getElementById('description').value,
        status: "pending",
        priority: priority
    };
    try {
        const response = await fetch('http://localhost:8080/tasks', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${accessToken}`
            },
            body: JSON.stringify(task)
        });
        if (response.status === 401) {
            if (await refreshTokenIfNeeded()) {
                return createTask();
            }
        }
        const data = await response.json();
        if (response.ok) {
            setMessage('task-message', 'Task created!');
            loadTasks();
        } else {
            setMessage('task-message', data.error || 'Failed to create task', true);
        }
    } catch (error) {
        setMessage('task-message', 'Error: ' + error.message, true);
    }
}

function openModal(id, title, description, status, priority, due_date) {
    document.getElementById('update-task-id').value = id;
    document.getElementById('update-title').value = title;
    document.getElementById('update-description').value = description;
    document.getElementById('update-status').value = status;
    document.getElementById('update-priority').value = priority;
    document.getElementById('update-due_date').value = due_date ? due_date.slice(0, 16) : '';
    document.getElementById('modal').style.display = 'block';
}

function closeModal() {
    document.getElementById('modal').style.display = 'none';
}

async function saveUpdate() {
    const id = document.getElementById('update-task-id').value;
    const task = {
        title: document.getElementById('update-title').value,
        description: document.getElementById('update-description').value,
        status: document.getElementById('update-status').value,
        priority: parseInt(document.getElementById('update-priority').value),
        due_date: formatDateTime(document.getElementById('update-due_date').value)
    };
    try {
        const response = await fetch(`http://localhost:8080/tasks/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${accessToken}`
            },
            body: JSON.stringify(task)
        });
        if (response.status === 401) {
            if (await refreshTokenIfNeeded()) {
                return saveUpdate();
            }
        }
        const data = await response.json();
        if (response.ok) {
            closeModal();
            setMessage('task-message', 'Task updated!');
            loadTasks();
        } else {
            setMessage('task-message', data.error || 'Failed to update task', true);
        }
    } catch (error) {
        setMessage('task-message', 'Error: ' + error.message, true);
    }
}

async function deleteTask(id) {
    if (!accessToken) {
        setMessage('task-message', 'Please login first', true);
        return;
    }
    try {
        const response = await fetch(`http://localhost:8080/tasks/${id}`, {
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${accessToken}` }
        });
        if (response.status === 401) {
            if (await refreshTokenIfNeeded()) {
                return deleteTask(id);
            }
        }
        const data = await response.json();
        if (response.ok) {
            setMessage('task-message', 'Task deleted!');
            loadTasks();
        } else {
            setMessage('task-message', data.error || 'Failed to delete task', true);
        }
    } catch (error) {
        setMessage('task-message', 'Error: ' + error.message, true);
    }
}