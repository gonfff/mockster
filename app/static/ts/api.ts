namespace Mockster {
  export async function apiRequest<T>(url: string, options: RequestInit = {}): Promise<{ response: Response; payload: T | null }> {
    const response = await fetch(url, options);
    let payload: T | null = null;
    try {
      payload = (await response.json()) as T;
    } catch {
      payload = null;
    }

    if (!response.ok) {
      const body = payload as ApiMessage | null;
      const message = body?.message || `Request failed with status ${response.status}`;
      const details = formatErrorMessage(body?.details);
      throw new Error(details ? `${message}. ${details}` : message);
    }

    return { response, payload };
  }
}
