/**
 * Apply theme variables as CSS custom properties on :root
 */
export function applyTheme(vars: Record<string, string>): void {
  for (const [key, value] of Object.entries(vars)) {
    document.documentElement.style.setProperty(`--${key}`, value);
  }
}

/**
 * Resolve icon URL - returns the URL as-is if present, undefined otherwise
 */
export function resolveIcon(src?: string): string | undefined {
  if (!src) return undefined;
  return src;
}
