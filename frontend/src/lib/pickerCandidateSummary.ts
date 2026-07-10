import type { gui } from '../../wailsjs/go/models'

export type PickerCandidateFacts = {
  matchesCount?: number
  score?: number
  unique: boolean
  visible: boolean
  warnings: string[]
}

export function pickerCandidateFacts(cand: Pick<gui.SelectorCandidate, 'matches_count' | 'score' | 'unique' | 'visible' | 'warnings'>): PickerCandidateFacts {
  return {
    matchesCount: cand.matches_count > 0 ? cand.matches_count : undefined,
    score: cand.score > 0 ? cand.score : undefined,
    unique: cand.unique,
    visible: cand.visible,
    warnings: [...(cand.warnings ?? [])],
  }
}
