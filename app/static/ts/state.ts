namespace Mockster {
  const STORAGE_KEYS = {
    ordering: "ordering",
    tableData: "tableData",
  };

  export const state: {
    ordering: SortField;
    mocksByName: Record<string, Mock>;
  } = {
    ordering: (localStorage.getItem(STORAGE_KEYS.ordering) as SortField) || "name",
    mocksByName: {},
  };

  export function persistState(): void {
    localStorage.setItem(STORAGE_KEYS.ordering, state.ordering);
    localStorage.setItem(STORAGE_KEYS.tableData, JSON.stringify(state.mocksByName));
  }

  export function setMocks(items: Mock[]): void {
    state.mocksByName = {};
    for (const item of items) {
      state.mocksByName[item.name] = item;
    }
    persistState();
  }

  export function removeMock(name: string): void {
    delete state.mocksByName[name];
    persistState();
  }

  export function getMocksList(): Mock[] {
    return Object.values(state.mocksByName);
  }
}
