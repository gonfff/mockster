namespace Mockster {
  export function sortMocks(mocks: Mock[], field: SortField): Mock[] {
    const copy = [...mocks];
    copy.sort((a, b) => {
      if (field === "status") {
        const left = Number(a.response?.status || 0);
        const right = Number(b.response?.status || 0);
        return left - right;
      }
      const left = String(a[field] || "");
      const right = String(b[field] || "");
      return left.localeCompare(right);
    });
    return copy;
  }

  export function renderTable(mocks: Mock[]): void {
    const tableBody = document.getElementById("table-body");
    if (!tableBody) {
      return;
    }
    tableBody.innerHTML = generateTableHTML(mocks);
  }

  export function renderCurrent(): void {
    const sorted = sortMocks(getMocksList(), state.ordering);
    renderTable(sorted);
  }

  function generateTableHTML(data: Mock[]): string {
    let tableHTML = "";

    data.forEach((item, index) => {
      const encodedName = encodeURIComponent(item.name);
      tableHTML += `
        <tr>
          <td class="align-middle">${index + 1}</td>
          <td class="text-primary align-middle">
            <span class="my-button" onClick="editModal(decodeURIComponent('${encodedName}'))">${escapeHtml(item.name)}</span>
          </td>
          <td class="align-middle">${escapeHtml(item.method)}</td>
          <td class="align-middle">${escapeHtml(item.path)}</td>
          <td class="align-middle">${escapeHtml(item.response?.status)}</td>
          <td class="align-middle">
            <button type="button" class="btn btn-outline-danger" onClick="deleteModal(decodeURIComponent('${encodedName}'))">
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

  export function showDeleteModal(name: string): void {
    const safeName = escapeHtml(name);
    const encodedName = encodeURIComponent(name);
    const modalHTML = `
      <div class="modal fade" id="deleteModal" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h1 class="modal-title fs-5" id="exampleModalLabel">Delete mock</h1>
              <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">Are you sure you want to delete ${safeName}?</div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
              <button type="button" class="btn btn-danger" data-bs-dismiss="modal" onClick="deleteMock(decodeURIComponent('${encodedName}'))">Delete</button>
            </div>
          </div>
        </div>
      </div>
    `;

    const container = document.getElementById("delete-modal-container");
    if (!container) {
      return;
    }
    container.innerHTML = modalHTML;
    const modal = new bootstrap.Modal(document.getElementById("deleteModal"));
    modal.show();
  }

  export function showToast(status: number | string, text: string): void {
    const safeText = escapeHtml(text);
    const isError = typeof status === "number" ? status >= 400 : String(status).toLowerCase().includes("error");
    const toneClass = isError ? "text-bg-danger" : "text-bg-success";

    let stack = document.getElementById("toast-stack");
    if (!stack) {
      const container = document.getElementById("toast-container");
      if (!container) {
        return;
      }
      container.innerHTML = '<div id="toast-stack" class="toast-container position-fixed bottom-0 end-0 p-3"></div>';
      stack = document.getElementById("toast-stack");
    }

    if (!stack) {
      return;
    }

    const toastEl = document.createElement("div");
    toastEl.className = `toast ${toneClass} border-0`;
    toastEl.role = "alert";
    toastEl.ariaLive = "assertive";
    toastEl.ariaAtomic = "true";
    toastEl.innerHTML = `
      <div class="d-flex align-items-center">
        <div class="toast-body">${safeText}</div>
        <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
      </div>
    `;

    stack.appendChild(toastEl);
    const toast = new bootstrap.Toast(toastEl, { delay: 3000 });
    toastEl.addEventListener("hidden.bs.toast", () => toastEl.remove());
    toast.show();
  }

  function editorForm(action: "create" | "edit", name = ""): string {
    const textAction = action === "edit" ? "Edit" : "Create";
    const encodedName = encodeURIComponent(name);
    const saveFunc = action === "edit"
      ? `saveEditForm(event, decodeURIComponent('${encodedName}'))`
      : "saveCreateForm(event)";

    return `
      <div class="modal" id="Editor" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
        <div class="modal-dialog modal-dialog-right">
          <div class="modal-content modal-content-right">
            <div class="modal-header">
              <h1 class="modal-title fs-5" id="exampleModalLabel">${textAction} Mock</h1>
              <button type="button" class="btn-close" aria-label="Close" onclick="hideModal()"></button>
            </div>
            <div class="modal-body" style="overflow-y: auto;">
              <form id="jsonForm">
                <div class="mb-3"><label for="name" class="form-label">Name</label><input type="text" class="form-control" id="name" name="name" required></div>
                <div class="mb-3"><label for="path" class="form-label">Path</label><input type="text" class="form-control" id="path" name="path" required></div>
                <div class="mb-3">
                  <label for="method" class="form-label">Method</label>
                  <select class="form-select" id="method" name="method" required>
                    <option value="GET">GET</option><option value="POST">POST</option><option value="PUT">PUT</option>
                    <option value="PATCH">PATCH</option><option value="DELETE">DELETE</option><option value="OPTIONS">OPTIONS</option><option value="HEAD">HEAD</option>
                  </select>
                </div>
                <div class="mb-3"><label for="requestHeaders" class="form-label">Request Headers (JSON)</label><textarea class="form-control" id="requestHeaders" name="requestHeaders"></textarea></div>
                <div class="mb-3"><label for="requestQueryParams" class="form-label">Request Query Params (JSON)</label><textarea class="form-control" id="requestQueryParams" name="requestQueryParams"></textarea></div>
                <div class="mb-3"><label for="requestCookies" class="form-label">Request Cookies (JSON)</label><textarea class="form-control" id="requestCookies" name="requestCookies"></textarea></div>
                <div class="mb-3"><label for="requestBody" class="form-label">Request Body</label><textarea class="form-control" id="requestBody" name="requestBody"></textarea></div>
                <div class="mb-3"><label for="responseStatus" class="form-label">Response Status</label><input type="number" class="form-control" id="responseStatus" name="responseStatus" required></div>
                <div class="mb-3"><label for="responseHeaders" class="form-label">Response Headers (JSON)</label><textarea class="form-control" id="responseHeaders" name="responseHeaders"></textarea></div>
                <div class="mb-3"><label for="responseCookies" class="form-label">Response Cookies (JSON)</label><textarea class="form-control" id="responseCookies" name="responseCookies"></textarea></div>
                <div class="mb-3"><label for="responseBody" class="form-label">Response Body</label><textarea class="form-control" id="responseBody" name="responseBody"></textarea></div>
                <div class="fixed-footer"><button type="submit" class="btn btn-primary" onclick="${saveFunc}">${textAction}</button></div>
              </form>
            </div>
          </div>
        </div>
      </div>
    `;
  }

  export function openCreateModal(): void {
    const container = document.getElementById("editor-modal-container");
    if (!container) {
      return;
    }
    container.innerHTML = editorForm("create");
    const modal = new bootstrap.Modal(document.getElementById("Editor"));
    modal.show();
  }

  export function openEditModal(mock: Mock): void {
    const container = document.getElementById("editor-modal-container");
    if (!container) {
      return;
    }

    container.innerHTML = editorForm("edit", mock.name);
    const modal = new bootstrap.Modal(document.getElementById("Editor"));
    modal.show();

    (document.getElementById("name") as HTMLInputElement).value = mock.name || "";
    (document.getElementById("path") as HTMLInputElement).value = mock.path || "";
    (document.getElementById("method") as HTMLSelectElement).value = mock.method || "GET";
    (document.getElementById("requestHeaders") as HTMLTextAreaElement).value = serializeJSON(mock.request?.headers);
    (document.getElementById("requestQueryParams") as HTMLTextAreaElement).value = serializeJSON(mock.request?.query_params);
    (document.getElementById("requestCookies") as HTMLTextAreaElement).value = serializeJSON(mock.request?.cookies);
    (document.getElementById("requestBody") as HTMLTextAreaElement).value = mock.request?.body || "";
    (document.getElementById("responseStatus") as HTMLInputElement).value = String(mock.response?.status || "");
    (document.getElementById("responseHeaders") as HTMLTextAreaElement).value = serializeJSON(mock.response?.headers);
    (document.getElementById("responseCookies") as HTMLTextAreaElement).value = serializeJSON(mock.response?.cookies);
    (document.getElementById("responseBody") as HTMLTextAreaElement).value = mock.response?.body || "";
  }

  export function hideEditorModal(): void {
    const modal = bootstrap.Modal.getInstance(document.getElementById("Editor"));
    if (modal) {
      modal.hide();
    }
  }
}
