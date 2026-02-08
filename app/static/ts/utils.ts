namespace Mockster {
  export function escapeHtml(value: unknown): string {
    if (value === null || value === undefined) {
      return "";
    }
    return String(value)
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll('"', "&quot;")
      .replaceAll("'", "&#039;");
  }

  export function formatErrorMessage(details: unknown): string {
    if (Array.isArray(details)) {
      return details.join("; ");
    }
    if (details === null || details === undefined) {
      return "";
    }
    return String(details);
  }

  export function parseJSON(data: string): Record<string, string> | null {
    if (!data) {
      return null;
    }
    return JSON.parse(data) as Record<string, string>;
  }

  export function serializeJSON(value: Record<string, string> | null | undefined): string {
    if (!value) {
      return "";
    }
    return JSON.stringify(value, null, 2);
  }
}
