declare const bootstrap: any;

interface Window {
  requestTable: () => Promise<void>;
  orderBy: (field: Mockster.SortField) => void;
  deleteModal: (name: string) => void;
  deleteMock: (name: string) => Promise<void>;
  importYaml: (e: Event) => Promise<void>;
  createModal: () => void;
  editModal: (name: string) => void;
  saveCreateForm: (e: Event) => Promise<void>;
  saveEditForm: (e: Event, name: string) => Promise<void>;
  searchMocks: (e: Event) => void;
  hideModal: () => void;
}
