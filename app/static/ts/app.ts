namespace Mockster {
  function resolveSuccessMessage(defaultMessage: string, message?: string): string {
    const normalized = String(message || "").trim().toLowerCase();
    if (!normalized || normalized === "ok" || normalized === "success") {
      return defaultMessage;
    }
    return String(message);
  }

  function formToMock(form: HTMLFormElement): Mock {
    const formData = new FormData(form);
    const field = (name: string): string => String(formData.get(name) || "");

    return {
      name: field("name"),
      path: field("path"),
      method: field("method"),
      request: {
        headers: parseJSON(field("requestHeaders")),
        query_params: parseJSON(field("requestQueryParams")),
        cookies: parseJSON(field("requestCookies")),
        body: field("requestBody") || null,
      },
      response: {
        status: Number.parseInt(field("responseStatus"), 10),
        headers: parseJSON(field("responseHeaders")),
        cookies: parseJSON(field("responseCookies")),
        body: field("responseBody") || null,
      },
    };
  }

  async function requestTable(): Promise<void> {
    try {
      const { payload } = await apiRequest<MocksPayload>("/management/mocks");
      setMocks(Array.isArray(payload?.items) ? payload.items : []);
      renderCurrent();
    } catch (error) {
      showToast("ERROR", (error as Error).message);
    }
  }

  function orderBy(field: SortField): void {
    state.ordering = field;
    persistState();
    renderCurrent();
  }

  function deleteModal(name: string): void {
    showDeleteModal(name);
  }

  async function deleteMock(name: string): Promise<void> {
    try {
      const { response } = await apiRequest<ApiMessage>(`/management/mocks/${encodeURIComponent(name)}`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
      });
      removeMock(name);
      renderCurrent();
      showToast(response.status, `Mock ${name} deleted`);
    } catch (error) {
      showToast("ERROR", (error as Error).message);
    }
  }

  async function importYaml(e: Event): Promise<void> {
    e.preventDefault();

    const form = document.getElementById("fileInput") as HTMLFormElement | null;
    if (!form) {
      return;
    }

    const formData = new FormData(form);
    const file = formData.get("file");
    if (!(file instanceof File) || file.size === 0) {
      showToast("ERROR", "Choose a YAML file first");
      return;
    }

    try {
      const { response, payload } = await apiRequest<ApiMessage>("/management/mocks/actions/import", {
        method: "POST",
        body: formData,
      });
      form.reset();
      showToast(response.status, resolveSuccessMessage("Imported", payload?.message));
      await requestTable();
    } catch (error) {
      showToast("ERROR", (error as Error).message);
    }
  }

  function createModal(): void {
    openCreateModal();
  }

  function editModal(name: string): void {
    const item = state.mocksByName[name];
    if (!item) {
      showToast("ERROR", `Mock ${name} not found`);
      return;
    }
    openEditModal(item);
  }

  function hideModal(): void {
    hideEditorModal();
  }

  async function saveCreateForm(e: Event): Promise<void> {
    e.preventDefault();
    const form = document.getElementById("jsonForm") as HTMLFormElement | null;
    if (!form) {
      return;
    }

    let payload: Mock;
    try {
      payload = formToMock(form);
    } catch (error) {
      showToast("ERROR", `Invalid JSON: ${(error as Error).message}`);
      return;
    }

    try {
      const { response, payload: data } = await apiRequest<ApiMessage>("/management/mocks", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      hideEditorModal();
      showToast(response.status, resolveSuccessMessage("Created", data?.message));
      await requestTable();
    } catch (error) {
      showToast("ERROR", (error as Error).message);
    }
  }

  async function saveEditForm(e: Event, name: string): Promise<void> {
    e.preventDefault();
    const form = document.getElementById("jsonForm") as HTMLFormElement | null;
    if (!form) {
      return;
    }

    let payload: Mock;
    try {
      payload = formToMock(form);
    } catch (error) {
      showToast("ERROR", `Invalid JSON: ${(error as Error).message}`);
      return;
    }

    try {
      const { response, payload: data } = await apiRequest<ApiMessage>(`/management/mocks/${encodeURIComponent(name)}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      hideEditorModal();
      showToast(response.status, resolveSuccessMessage("Updated", data?.message));
      await requestTable();
    } catch (error) {
      showToast("ERROR", (error as Error).message);
    }
  }

  function searchMocks(): void {
		const input = document.getElementById("search-query") as HTMLInputElement | null;
		if (!input) {
			return;
		}

		const query = input.value.trim().toLowerCase();
		if (!query) {
			renderCurrent();
			return;
		}

		const filtered = sortMocks(
			getMocksList().filter((item) => String(item.name).toLowerCase().includes(query)),
			state.ordering
		);
		renderTable(filtered);
	}

  function exposeGlobals(): void {
    window.requestTable = requestTable;
    window.orderBy = orderBy;
    window.deleteModal = deleteModal;
    window.deleteMock = deleteMock;
    window.importYaml = importYaml;
    window.createModal = createModal;
    window.editModal = editModal;
    window.saveCreateForm = saveCreateForm;
    window.saveEditForm = saveEditForm;
    window.searchMocks = searchMocks;
    window.hideModal = hideModal;
  }

  document.addEventListener("DOMContentLoaded", async () => {
    exposeGlobals();
    await requestTable();
  });
}
