import { describe, it, expect, vi } from 'vitest'
import { resolveProjectPathInput } from './projectPath'

describe('resolveProjectPathInput', () => {
  it('resolves bundled examples alias through provider', async () => {
    const bundledExamplesPath = vi.fn(async () => 'C:/repo/examples')
    await expect(resolveProjectPathInput('examples', bundledExamplesPath)).resolves.toBe('C:/repo/examples')
    expect(bundledExamplesPath).toHaveBeenCalledOnce()
  })

  it('preserves explicit project paths', async () => {
    const bundledExamplesPath = vi.fn(async () => 'C:/repo/examples')
    await expect(resolveProjectPathInput('C:/work/project', bundledExamplesPath)).resolves.toBe('C:/work/project')
    expect(bundledExamplesPath).not.toHaveBeenCalled()
  })
})
