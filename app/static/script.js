document.addEventListener('DOMContentLoaded', async function () {
    await requestTable();
})
    
async function requestTable(){
    const tableBody = document.getElementById('table-body');
    
    try {
        const response = await fetch('/management/mocks');
        if (response.ok) {
            const mocks = await response.json();
            
            const listData = mocks.items.map(item => ({
                ...item,
            }));
            const mapData = {};
            for (const item of listData) {
                mapData[item.name] = item;
            }

            localStorage.setItem('tableData', JSON.stringify(mapData));
            ordering = localStorage.getItem('ordering');
            if (ordering) {
                orderBy(ordering);
            } else {
                orderBy('name');
            }
            localStorage.setItem('ordering', 'name');
            
            const tableHTML = generateTableHTML(listData);
            tableBody.innerHTML = tableHTML;
        } else {
            console.error('Failed to fetch JSON data');
        }
    } catch (error) {
        console.error('An error occurred:', error);
    }
}

function generateTableHTML(data) {
    let tableHTML = '';

    let i = 0;
    data.forEach(item => {
        i += 1;
        tableHTML += `
            <tr>
                <td class="align-middle">${i}</td>
                <td class="text-primary align-middle"><span class="my-button" onClick="editModal('${item.name}')">${item.name}</span></td>
                <td class="align-middle">${item.method}</td>
                <td class="align-middle">${item.path}</td>
                <td class="align-middle">${item.response.status}</td>
                <td class="align-middle">
                    <button type="button" class="btn btn-outline-danger" onClick="deleteModal('${item.name}')">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash3" viewBox="0 0 16 16">
                            <path d="M6.5 1h3a.5.5 0 0 1 .5.5v1H6v-1a.5.5 0 0 1 .5-.5ZM11 2.5v-1A1.5 1.5 0 0 0 9.5 0h-3A1.5 1.5 0 0 0 5 1.5v1H2.506a.58.58 0 0 0-.01 0H1.5a.5.5 0 0 0 0 1h.538l.853 10.66A2 2 0 0 0 4.885 16h6.23a2 2 0 0 0 1.994-1.84l.853-10.66h.538a.5.5 0 0 0 0-1h-.995a.59.59 0 0 0-.01 0H11Zm1.958 1-.846 10.58a1 1 0 0 1-.997.92h-6.23a1 1 0 0 1-.997-.92L3.042 3.5h9.916Zm-7.487 1a.5.5 0 0 1 .528.47l.5 8.5a.5.5 0 0 1-.998.06L5 5.03a.5.5 0 0 1 .47-.53Zm5.058 0a.5.5 0 0 1 .47.53l-.5 8.5a.5.5 0 1 1-.998-.06l.5-8.5a.5.5 0 0 1 .528-.47ZM8 4.5a.5.5 0 0 1 .5.5v8.5a.5.5 0 0 1-1 0V5a.5.5 0 0 1 .5-.5Z"/>
                        </svg>
                    </button>
                </td>        
            </tr>
        `;
    });
    return tableHTML;
}

function orderBy(field){
    const tableData = JSON.parse(localStorage.getItem('tableData'));
    const listData = Object.values(tableData);
    let sortedData;
    if (field === 'status') {
    sortedData = listData.sort((a, b) => {
        if (a['response'][field] > b['response'][field]) {
            return 1;
        } else if (a['response'][field] < b['response'][field]) {
            return -1;
        } else {
            return 0;
        }
    });
    } else {
    sortedData = listData.sort((a, b) => {
        if (a[field] > b[field]) {
            return 1;
        } else if (a[field] < b[field]) {
            return -1;
        } else {
            return 0;
        }
    });
    }
    localStorage.setItem('ordering', field);
    const tableHTML = generateTableHTML(sortedData);
    document.getElementById('table-body').innerHTML = tableHTML;
}

function deleteModal(name) {
    const modalHTML = `
    <div class="modal fade" id="deleteModal" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h1 class="modal-title fs-5" id="exampleModalLabel">Delete mock</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                    Are you sure you want to delete ${name}?
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                    <button type="button" class="btn btn-danger" data-bs-dismiss="modal" onClick="deleteMock('${name}')">Delete</button>
                </div>
            </div>
        </div>
    </div>`;
    document.getElementById('delete-modal-container').innerHTML = modalHTML; 
    const deleteModal = new bootstrap.Modal(document.getElementById('deleteModal')); 
    deleteModal.show(); 
}

function toasts(title, status, text) {
    const toastHTML = `
    <div class="toast-container position-fixed bottom-0 end-0 p-3">
        <div id="liveToast" class="toast" role="alert" aria-live="assertive" aria-atomic="true">
            <div class="toast-header">
                <strong class="me-auto">${title}</strong>
                <small>${status}</small>
                <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
            </div>
            <div class="toast-body">
                ${text}
            </div>
        </div>
    </div>`;
    document.getElementById('toast-container').innerHTML = toastHTML;
    const toastEl = document.getElementById('liveToast');
    const toast = new bootstrap.Toast(toastEl);
    toast.show();
}

function deleteMock(name) {
    const deleteModal = new bootstrap.Modal(document.getElementById('deleteModal'));


    fetch(`/management/mocks/${name}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
        },
    })
    .then(response => {
        if (response.ok) {
            return response.json().then(data => {
                toasts(`delete ${name}`, response.status, `Mock ${name} deleted`);
            }); 
        } else {
            return response.json().then(data => {
                toasts(`delete ${name}`, response.status, data.message+'. '+data.details);
            });
        }
    })
    .catch(error => {
        toasts('delete mock', 'ERROR', error.message);
    });

    const tableData = JSON.parse(localStorage.getItem('tableData'));
    delete tableData[name];
    localStorage.setItem('tableData', JSON.stringify(tableData));

    const ordering = localStorage.getItem('ordering');
    orderBy(ordering);

}

async function importYaml(e) {
    e.preventDefault();
  
    const form = document.getElementById('fileInput'); 
    const formData = new FormData(form);

    try {
        const response = await fetch('/management/mocks/actions/import', {
            method: 'POST',
            body: formData,
        });

        form.reset();
        
        const data = await response.json();
        toasts('import', response.status, data.message);

        await requestTable();
        const ordering = localStorage.getItem('ordering');
        orderBy(ordering);

    } catch (error) {
        toasts('import', 'ERROR', error.message);
    }
}

function editorForm(action, name=""){
    const textAction = action === 'edit' ? 'Edit' : 'Create';
    const saveFunc = action === 'edit' ? `saveEditForm('${name}')` : 'saveCreateForm()';


    const editorFormHTML = `
    <div class="modal" id="Editor" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
        <div class="modal-dialog modal-dialog-right">
            <div class="modal-content modal-content-right">
                <div class="modal-header">
                    <h1 class="modal-title fs-5" id="exampleModalLabel">${textAction} Mock</h1>
                    <button type="button" class="btn-close" aria-label="Close" onclick="hideModal()"></button>
                </div>
                <div class="modal-body" style="overflow-y: auto;">
                    <form id="jsonForm">
                        <div class="mb-3">
                        <label for="name" class="form-label">Name</label>
                        <input type="text" class="form-control" id="name" name="name" required>

                        </div>
                        <div class="mb-3">
                        <label for="path" class="form-label">Path</label>
                        <input type="text" class="form-control" id="path" name="path" required>
                        </div>
                        <div class="mb-3">
                            <label for="method" class="form-label">Method</label>
                            <select class="form-select" id="method" name="method" required>
                                <option value="GET">GET</option>
                                <option value="POST">POST</option>
                                <option value="PUT">PUT</option>
                                <option value="PATCH">PATCH</option>
                                <option value="DELETE">DELETE</option>
                                <option value="OPTIONS">OPTIONS</option>
                                <option value="HEAD">HEAD</option>
                            </select>
                        </div>
                        <div class="mb-3">
                            <label for="requestHeaders" class="form-label">Request Headers (JSON)</label>
                            <textarea class="form-control" id="requestHeaders" name="requestHeaders"></textarea>
                        </div>
                        <div class="mb-3">
                            <label for="requestQueryParams" class="form-label">Request Query Params (JSON)</label>
                            <textarea class="form-control" id="requestQueryParams" name="requestQueryParams"></textarea>
                        </div>
                        <div class="mb-3">
                            <label for="requestCookies" class="form-label">Request Cookies (JSON)</label>
                            <textarea class="form-control" id="requestCookies" name="requestCookies"></textarea>
                        </div>
                        <div class="mb-3">
                            <label for="requestBody" class="form-label">Request Body</label>
                            <textarea class="form-control" id="requestBody" name="requestBody"></textarea>
                        </div>
                        <div class="mb-3">
                            <label for="responseStatus" class="form-label">Response Status</label>
                            <input type="number" class="form-control" id="responseStatus" name="responseStatus"
                                required>
                        </div>
                        <div class="mb-3">
                            <label for="responseHeaders" class="form-label">Response Headers (JSON)</label>
                            <textarea class="form-control" id="responseHeaders" name="responseHeaders"></textarea>
                        </div>
                        <div class="mb-3">
                            <label for="responseCookies" class="form-label">Response Cookies (JSON)</label>
                            <textarea class="form-control" id="responseCookies" name="responseCookies"></textarea>
                        </div>
                        <div class="mb-3">
                            <label for="responseBody" class="form-label">Response Body</label>
                            <textarea class="form-control" id="responseBody" name="responseBody"></textarea>
                        </div>
                        <div class="fixed-footer">
                            <button type="submit" class="btn btn-primary" onclick="${saveFunc}">${textAction}</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>`;
    return editorFormHTML;
}

function createModal() {
    const editorFormHTML = editorForm('create');
    document.getElementById('editor-modal-container').innerHTML = editorFormHTML; 
    const Modal = new bootstrap.Modal(document.getElementById('Editor')); 
    Modal.show(); 
}

function makeEditModal(mock){
    const editorFormHTML = editorForm('edit', mock.name);
    document.getElementById('editor-modal-container').innerHTML = editorFormHTML; 
    const Modal = new bootstrap.Modal(document.getElementById('Editor')); 
    Modal.show(); 

    document.getElementById('name').value = mock.name;
    document.getElementById('path').value = mock.path;
    document.getElementById('method').value = mock.method;

    requestHeaders = JSON.stringify(mock.request.headers, null, 2);
    document.getElementById('requestHeaders').value = requestHeaders === 'null' ? '' : requestHeaders;

    requestQueryParams = JSON.stringify(mock.request.query_params, null, 2);
    document.getElementById('requestQueryParams').value = requestQueryParams === 'null' ? '' : requestQueryParams;

    requestCookies = JSON.stringify(mock.request.cookies, null, 2);
    document.getElementById('requestCookies').value = requestCookies === 'null' ? '' : requestCookies;
    
    document.getElementById('requestBody').value = mock.request.body;

    document.getElementById('responseStatus').value = mock.response.status;

    responseHeaders = JSON.stringify(mock.response.headers, null, 2);
    document.getElementById('responseHeaders').value = responseHeaders === 'null' ? '' : responseHeaders;

    responseCookies = JSON.stringify(mock.response.cookies, null, 2);
    document.getElementById('responseCookies').value = responseCookies === 'null' ? '' : responseCookies;

    document.getElementById('responseBody').value = mock.response.body;   
}

function hideModal() {
    const Modal = bootstrap.Modal.getInstance(document.getElementById('Editor'));
    Modal.hide();
}

function plainToMockObj(data){
    const obj = {};
    obj.response = {};
    obj.request = {};

    obj.name = data.name;
    obj.path = data.path;
    obj.method = data.method;

    obj.request.headers = parseJSON(data.requestHeaders);
    obj.request.query_params = parseJSON(data.requestQueryParams);
    obj.request.cookies = parseJSON(data.requestCookies);
    obj.request.body = data.requestBody || null;

    obj.response.status = parseInt(data.responseStatus);
    obj.response.headers = parseJSON(data.responseHeaders);
    obj.response.cookies = parseJSON(data.responseCookies);
    obj.response.body = data.responseBody || null;

    return obj;
}

function parseJSON(data) {
    if (!data) {
        return null;
    }

    try {
        return JSON.parse(data);
    } catch (error) {
        throw error;
    }
}

async function saveCreateForm(name) {
    event.preventDefault();
    const form = document.getElementById('jsonForm'); 

    let formData;
    try {
        formData = new FormData(form);
    } catch (error) {
        toasts('create', 'InvalidJson', error.message);
        return;
    }

    const jsonObject = {};
    formData.forEach((value, key) => {
        jsonObject[key] = value;
    });

    let sendObj;
    try {
        sendObj = plainToMockObj(jsonObject);
        } catch (error) {
            if (error instanceof SyntaxError) {
            toasts('create', 'ValidationFailed', error.message);
            return;
            }
        }
    
    try {
        const response = await fetch("/management/mocks", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(sendObj)
        });

        if (response.ok) {
            hideModal();
            const data = await response.json();
            toasts('create', response.status, data.message);
            await requestTable();

        } else {
            const data = await response.json();
            toasts('create', response.status, data.message +'. '+ data.details);
        }
    } catch (error) {
        toasts('create', 'Error', error.message);
    }
}

function editModal(name) {
    const tableData = JSON.parse(localStorage.getItem('tableData'));
    const item = tableData[name];
    makeEditModal(item);
}

async function saveEditForm(name){
    event.preventDefault();
    const form = document.getElementById('jsonForm'); 

    let formData;
    try {
        formData = new FormData(form);
    } catch (error) {
        toasts('edit', 'InvalidJson', error.message);
        return;
    }

    const jsonObject = {};
    formData.forEach((value, key) => {
        jsonObject[key] = value;
    });

    let sendObj;
    try {
        sendObj = plainToMockObj(jsonObject);
        } catch (error) {
            if (error instanceof SyntaxError) {
            toasts('edit', 'ValidationFailed', error.message);
            return;
            }
        }
    
    try {
        const response = await fetch(`/management/mocks/${name}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(sendObj)
        });

        if (response.ok) {
            hideModal();
            const data = await response.json();
            toasts('edit', response.status, data.message);
            await requestTable();

        } else {
            const data = await response.json();
            toasts('edit', response.status, data.message +'. '+ data.details);
        }
    } catch (error) {
        toasts('edit', 'Error', error.message);
    }
}

function searchMocks(){
    event.preventDefault();
    const form = document.getElementById('search-form'); 
    const formData = new FormData(form);
    const query = formData.get('query'); 

    if (!query) {
        requestTable();
    }

    const tableBody = document.getElementById('table-body');

    const tableData = JSON.parse(localStorage.getItem('tableData'));
    const listData = Object.values(tableData);
    const filteredData = listData.filter(item => item.name.includes(query));

    const tableHTML = generateTableHTML(filteredData);
    tableBody.innerHTML = tableHTML;
}