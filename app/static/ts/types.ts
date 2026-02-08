namespace Mockster {
  export type SortField = "name" | "method" | "path" | "status";

  export interface MockRequest {
    headers: Record<string, string> | null;
    query_params: Record<string, string> | null;
    cookies: Record<string, string> | null;
    body: string | null;
  }

  export interface MockResponse {
    status: number;
    headers: Record<string, string> | null;
    cookies: Record<string, string> | null;
    body: string | null;
  }

  export interface Mock {
    name: string;
    path: string;
    method: string;
    request: MockRequest;
    response: MockResponse;
  }

  export interface ApiMessage {
    message?: string;
    details?: unknown;
  }

  export interface MocksPayload {
    items?: Mock[];
  }
}
