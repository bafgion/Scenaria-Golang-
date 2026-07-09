export async function resolveProjectPathInput(
  path: string,
  bundledExamplesPath: () => Promise<string> | string,
): Promise<string> {
  const trimmed = path.trim()
  if (!trimmed) return ''
  const normalized = trimmed.replace(/\\/g, '/')
  if (normalized === 'examples') {
    const bundled = await Promise.resolve(bundledExamplesPath()).catch(() => '')
    if (bundled) return bundled
  }
  return trimmed
}
