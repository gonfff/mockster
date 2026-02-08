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
    if (data.length === 0) {
      return `
        <tr>
          <td colspan="6" class="empty-state-cell">
            <div class="empty-state-box">
              <p class="empty-state-title">No mocks found</p>
              <p class="empty-state-subtitle">Try changing filters or create a new mock.</p>
              <button type="button" class="btn btn-brand btn-sm" onclick="createModal()">Create first mock</button>
            </div>
          </td>
        </tr>
      `;
    }

    let tableHTML = "";

    data.forEach((item, index) => {
      const encodedName = encodeURIComponent(item.name);
      const method = String(item.method || "").toUpperCase();
      const methodClass = `method-badge method-${method.toLowerCase()}`;
      tableHTML += `
        <tr>
          <td class="align-middle">${index + 1}</td>
          <td class="text-primary align-middle">
            <span class="my-button" onClick="editModal(decodeURIComponent('${encodedName}'))">${escapeHtml(item.name)}</span>
          </td>
          <td class="align-middle"><span class="${methodClass}">${escapeHtml(method)}</span></td>
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
          <div class="modal-content modal-card">
            <div class="modal-header">
              <h1 class="modal-title fs-5" id="exampleModalLabel">Delete mock</h1>
              <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
              <p class="mb-1">Are you sure you want to delete <strong>${safeName}</strong>?</p>
              <p class="text-muted small mb-0">This action cannot be undone.</p>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-ghost" data-bs-dismiss="modal">Close</button>
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
    const toneClass = isError ? "toast-danger" : "toast-success";
    const icon = isError ? "!" : "OK";

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
    toastEl.className = `toast toast-shell ${toneClass} border-0`;
    toastEl.role = "alert";
    toastEl.ariaLive = "assertive";
    toastEl.ariaAtomic = "true";
    toastEl.innerHTML = `
      <div class="d-flex align-items-start">
        <div class="toast-icon">${icon}</div>
        <div class="toast-body">
          <div class="toast-message">${safeText}</div>
        </div>
        <button type="button" class="btn-close ms-2 mt-2" data-bs-dismiss="toast" aria-label="Close"></button>
      </div>
    `;

    stack.appendChild(toastEl);
    const toast = new bootstrap.Toast(toastEl, { delay: 4200 });
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
      <div class="modal fade" id="Editor" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
        <div class="modal-dialog modal-dialog-right">
          <div class="modal-content modal-content-right editor-shell">
            <div class="modal-header">
              <div>
                <h1 class="modal-title fs-5" id="exampleModalLabel">${textAction} Mock</h1>
                <p class="text-muted small mb-0">Configure request matching and response output</p>
              </div>
              <button type="button" class="btn-close" aria-label="Close" onclick="hideModal()"></button>
            </div>
            <div class="modal-body editor-body">
              <form id="jsonForm">
                <section class="form-section">
                  <h2 class="form-section-title">General</h2>
                  <div class="row g-3">
                    <div class="col-md-7"><label for="name" class="form-label">Name</label><input type="text" class="form-control" id="name" name="name" required></div>
                    <div class="col-md-5">
                      <label for="method" class="form-label">Method</label>
                      <select class="form-select" id="method" name="method" required>
                        <option value="GET">GET</option><option value="POST">POST</option><option value="PUT">PUT</option>
                        <option value="PATCH">PATCH</option><option value="DELETE">DELETE</option><option value="OPTIONS">OPTIONS</option><option value="HEAD">HEAD</option>
                      </select>
                    </div>
                    <div class="col-12"><label for="path" class="form-label">Path</label><input type="text" class="form-control" id="path" name="path" placeholder="/users/profile" required></div>
                  </div>
                </section>

                <section class="form-section">
                  <h2 class="form-section-title">Request Matchers</h2>
                  <div class="mb-3"><label for="requestBody" class="form-label">Request Body</label><textarea class="form-control" id="requestBody" name="requestBody" rows="3"></textarea></div>
                  <div class="row g-3">
                    <div class="col-md-4"><label for="requestHeaders" class="form-label">Headers (JSON)</label><textarea class="form-control code-area" id="requestHeaders" name="requestHeaders" rows="4"></textarea></div>
                    <div class="col-md-4"><label for="requestQueryParams" class="form-label">Query Params (JSON)</label><textarea class="form-control code-area" id="requestQueryParams" name="requestQueryParams" rows="4"></textarea></div>
                    <div class="col-md-4"><label for="requestCookies" class="form-label">Cookies (JSON)</label><textarea class="form-control code-area" id="requestCookies" name="requestCookies" rows="4"></textarea></div>
                  </div>
                </section>

                <section class="form-section">
                  <h2 class="form-section-title">Response</h2>
                  <div class="row g-3">
                    <div class="col-md-4"><label for="responseStatus" class="form-label">Status</label><input type="number" class="form-control" id="responseStatus" name="responseStatus" min="100" max="599" required></div>
                    <div class="col-md-4"><label for="responseHeaders" class="form-label">Headers (JSON)</label><textarea class="form-control code-area" id="responseHeaders" name="responseHeaders" rows="4"></textarea></div>
                    <div class="col-md-4"><label for="responseCookies" class="form-label">Cookies (JSON)</label><textarea class="form-control code-area" id="responseCookies" name="responseCookies" rows="4"></textarea></div>
                    <div class="col-12"><label for="responseBody" class="form-label">Body</label><textarea class="form-control" id="responseBody" name="responseBody" rows="4"></textarea></div>
                  </div>
                </section>

                <div class="fixed-footer">
                  <button type="button" class="btn btn-ghost" onclick="hideModal()">Cancel</button>
                  <button type="submit" class="btn btn-brand" onclick="${saveFunc}">${textAction}</button>
                </div>
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
