// The frontend talks only to its own backend, on the same origin (docs/adr/0001). No URLs are baked
// in at build time.

export type ValueCount = { value: string; count: number }

export type QuestionSummary = {
  question_id: string
  label: string
  answered: number
  values: ValueCount[]
}

export type Demographics = {
  member_count: number
  questions: QuestionSummary[]
}

export async function fetchDemographics(signal?: AbortSignal): Promise<Demographics> {
  const res = await fetch('/api/community/demographics', { signal })
  if (!res.ok) {
    throw new Error(`Loading demographics failed (${res.status})`)
  }
  return res.json()
}
