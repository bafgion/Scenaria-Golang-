/** Canonical feature path for frontend comparisons and model keys. */
export function canonicalFeaturePath(path: string): string {
  return path.trim().replace(/\\/g, '/')
}
