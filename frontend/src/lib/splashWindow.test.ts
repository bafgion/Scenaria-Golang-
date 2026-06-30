import { describe, expect, it } from 'vitest'
import { MAIN_WINDOW, SPLASH_WINDOW } from './splashWindow'

describe('splashWindow constants', () => {
  it('uses compact splash dimensions', () => {
    expect(SPLASH_WINDOW.width).toBeLessThan(MAIN_WINDOW.width)
    expect(SPLASH_WINDOW.height).toBeLessThan(MAIN_WINDOW.height)
  })
})
