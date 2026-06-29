/** Catalog tree — parity with Python app/mvc/models/catalog_model.py */

export type RowKind = 'root' | 'dir' | 'file'
export type EmptyKind = 'no_project' | 'missing' | 'no_files' | 'no_match'

export type RunStatus = {
  success: boolean
  message: string
  at: string
  runner: string
}

export interface CatalogNode {
  kind: RowKind
  path: string
  name: string
  children: CatalogNode[]
  runSuccess: boolean | null
  runMessage: string
  runAt: string
  runRunner: string
}

export interface CatalogViewState {
  tree: CatalogNode | null
  emptyTitle: string | null
  emptyHint: string | null
  emptyKind: EmptyKind | null
  expandAll: boolean
  showEmptyMessage: boolean
}

function basename(path: string): string {
  const norm = path.replace(/\\/g, '/')
  const i = norm.lastIndexOf('/')
  return i >= 0 ? norm.slice(i + 1) : norm
}

function stem(path: string): string {
  const name = basename(path)
  return name.endsWith('.feature') ? name.slice(0, -8) : name
}

function normPath(path: string): string {
  return path.replace(/\\/g, '/').toLowerCase()
}

export function countFeatureFiles(node: CatalogNode | null): number {
  if (!node) return 0
  let count = node.kind === 'file' ? 1 : 0
  for (const child of node.children) count += countFeatureFiles(child)
  return count
}

export function normTag(tag: string): string {
  return tag.replace(/^@/, '').trim().toLowerCase()
}

export function parseCatalogFilter(filterText: string): { query: string; tag: string } {
  const raw = (filterText || '').trim()
  if (raw.startsWith('@')) {
    const tag = raw.slice(1).trim().split(/\s+/)[0]?.toLowerCase() || ''
    return { query: '', tag }
  }
  const lowered = raw.toLowerCase()
  if (lowered.startsWith('tag:')) {
    return { query: '', tag: raw.split(':', 2)[1]?.trim().toLowerCase() || '' }
  }
  return { query: raw.toLowerCase(), tag: '' }
}

function featureHasTag(fileTags: string[] | undefined, tag: string): boolean {
  if (!tag) return true
  if (!fileTags?.length) return false
  const needle = normTag(tag)
  return fileTags.some((item) => normTag(item) === needle)
}

function fileMatchesFilter(
  node: CatalogNode,
  query: string,
  tag: string,
  tagsByPath: Map<string, string[]>,
): boolean {
  if (tag) {
    return featureHasTag(tagsByPath.get(normPath(node.path)), tag)
  }
  if (!query) return true
  const rel = node.path.replace(/\\/g, '/').toLowerCase()
  return node.name.toLowerCase().includes(query) || rel.includes(query)
}

function filterTree(
  node: CatalogNode,
  query: string,
  tag: string,
  tagsByPath: Map<string, string[]>,
): CatalogNode | null {
  if (!query && !tag) return node

  if (node.kind === 'file') {
    return fileMatchesFilter(node, query, tag, tagsByPath) ? node : null
  }

  const children: CatalogNode[] = []
  for (const child of node.children) {
    const filtered = filterTree(child, query, tag, tagsByPath)
    if (filtered) children.push(filtered)
  }

  if (node.kind === 'root') {
    return { ...node, children }
  }
  if (!children.length) return null
  return { ...node, children }
}

export function buildCatalogTree(
  root: string,
  featurePaths: string[],
  runByPath: Map<string, RunStatus>,
): CatalogNode {
  const rootName = basename(root)
  const rootNode: CatalogNode = {
    kind: 'root',
    path: root,
    name: rootName,
    children: [],
    runSuccess: null,
    runMessage: '',
    runAt: '',
    runRunner: '',
  }
  const dirMap = new Map<string, CatalogNode>([['', rootNode]])

  const sorted = [...featurePaths].sort((a, b) => normPath(a).localeCompare(normPath(b)))
  const rootNorm = root.replace(/\\/g, '/').replace(/\/$/, '')

  for (const featurePath of sorted) {
    const rel = featurePath.replace(/\\/g, '/')
    if (!rel.toLowerCase().startsWith(rootNorm.toLowerCase())) continue
    const suffix = rel.slice(rootNorm.length).replace(/^\//, '')
    if (!suffix.endsWith('.feature')) continue

    const parts = suffix.split('/')
    let parent = rootNode
    const prefixParts: string[] = []

    for (const part of parts.slice(0, -1)) {
      prefixParts.push(part)
      const key = prefixParts.join('/')
      let dirNode = dirMap.get(key)
      if (!dirNode) {
        const dirPath = `${rootNorm}/${key}`.replace(/\//g, '\\')
        dirNode = {
          kind: 'dir',
          path: dirPath,
          name: part,
          children: [],
          runSuccess: null,
          runMessage: '',
          runAt: '',
          runRunner: '',
        }
        dirMap.set(key, dirNode)
        parent.children.push(dirNode)
      }
      parent = dirNode
    }

    const run = runByPath.get(normPath(featurePath))
    parent.children.push({
      kind: 'file',
      path: featurePath,
      name: stem(featurePath),
      children: [],
      runSuccess: run ? run.success : null,
      runMessage: run?.message || '',
      runAt: run?.at || '',
      runRunner: run?.runner || '',
    })
  }

  return rootNode
}

export function buildCatalogViewState(
  root: string | null,
  featurePaths: string[],
  filterText: string,
  runByPath: Map<string, RunStatus>,
  rootExists = true,
  tagsByPath: Map<string, string[]> = new Map(),
): CatalogViewState {
  if (!root) {
    return {
      tree: null,
      emptyTitle: 'Проект не открыт',
      emptyHint: 'Откройте папку с .feature сценариями.\nПроект → Открыть проект…',
      emptyKind: 'no_project',
      expandAll: false,
      showEmptyMessage: true,
    }
  }
  if (!rootExists) {
    return {
      tree: null,
      emptyTitle: 'Папка не найдена',
      emptyHint: `Путь недоступен:\n${root}\n\nВыберите другой проект.`,
      emptyKind: 'missing',
      expandAll: false,
      showEmptyMessage: true,
    }
  }

  const fullTree = buildCatalogTree(root, featurePaths, runByPath)
  const totalFiles = countFeatureFiles(fullTree)
  const { query, tag } = parseCatalogFilter(filterText)
  const tree = query || tag ? filterTree(fullTree, query, tag, tagsByPath) : fullTree
  const visibleFiles = countFeatureFiles(tree)

  if (totalFiles === 0) {
    return {
      tree,
      emptyTitle: 'Нет сценариев',
      emptyHint: `В «${basename(root)}» пока нет .feature файлов.\nНажмите + или Сценарий → Новый.\nДля пакетного запуска: режим «Выбор» или Ctrl+клик по файлу.`,
      emptyKind: 'no_files',
      expandAll: false,
      showEmptyMessage: true,
    }
  }

  if ((query || tag) && visibleFiles === 0) {
    const trimmed = filterText.trim()
    if (tag && !query) {
      return {
        tree,
        emptyTitle: 'Нет сценариев с тегом',
        emptyHint: `Тег «@${tag}» не найден ни в одном сценарии.\nОчистите поле поиска.`,
        emptyKind: 'no_match',
        expandAll: true,
        showEmptyMessage: true,
      }
    }
    return {
      tree,
      emptyTitle: 'Ничего не найдено',
      emptyHint: `Запрос «${trimmed}» не дал результатов.\nОчистите поле поиска.`,
      emptyKind: 'no_match',
      expandAll: true,
      showEmptyMessage: true,
    }
  }

  return {
    tree,
    emptyTitle: null,
    emptyHint: null,
    emptyKind: null,
    expandAll: Boolean(query || tag),
    showEmptyMessage: false,
  }
}

export function formatLastRunSummary(node: CatalogNode): string {
  if (node.runSuccess === null) return 'Последний прогон: нет данных'
  const status = node.runSuccess ? 'успех' : 'ошибка'
  const parts = [`Последний прогон: ${status}`]
  if (node.runAt) parts.push(`Время: ${node.runAt}`)
  if (node.runRunner) parts.push(`Runner: ${node.runRunner}`)
  if (!node.runSuccess && node.runMessage) parts.push(`Сообщение: ${node.runMessage}`)
  return parts.join('\n')
}

export function fileTreeLabel(
  node: CatalogNode,
  selectionMode: boolean,
  selected: boolean,
): string {
  let badge = '○'
  if (node.runSuccess === true) badge = '✓'
  else if (node.runSuccess === false) badge = '✗'

  let mark = ''
  if (selectionMode) {
    mark = selected ? '☑ ' : '☐ '
  } else if (selected) {
    mark = '◉ '
  }

  return `${mark}${badge} ${node.name}`
}

export function fileTreeTooltip(node: CatalogNode): string {
  const parts = [basename(node.path)]
  parts.push(formatLastRunSummary(node))
  return parts.join('\n')
}

export function collectFeaturePathsUnder(node: CatalogNode): string[] {
  if (node.kind === 'file') return [node.path]
  const out: string[] = []
  for (const child of node.children) out.push(...collectFeaturePathsUnder(child))
  return out
}
