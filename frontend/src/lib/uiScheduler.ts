/** Отложить тяжёлую работу, чтобы UI успел отрисоваться (кнопки, переключатели). */
export function deferToNextFrame(work: () => void): void {
  if (typeof requestAnimationFrame === 'function') {
    requestAnimationFrame(work)
    return
  }
  setTimeout(work, 0)
}

export type Debounced<T extends (...args: never[]) => void> = T & {
  cancel: () => void
  flush: () => void
}

/** Debounce с cancel/flush — для поиска и фильтров без блокировки ввода. */
export function debounce<T extends (...args: never[]) => void>(
  fn: T,
  delayMs: number,
): Debounced<T> {
  let timer: ReturnType<typeof setTimeout> | null = null
  let lastArgs: Parameters<T> | null = null

  const debounced = ((...args: Parameters<T>) => {
    lastArgs = args
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      timer = null
      const callArgs = lastArgs
      lastArgs = null
      if (callArgs) fn(...callArgs)
    }, delayMs)
  }) as Debounced<T>

  debounced.cancel = () => {
    if (timer) clearTimeout(timer)
    timer = null
    lastArgs = null
  }

  debounced.flush = () => {
    if (!timer || !lastArgs) return
    clearTimeout(timer)
    timer = null
    const callArgs = lastArgs
    lastArgs = null
    fn(...callArgs)
  }

  return debounced
}
