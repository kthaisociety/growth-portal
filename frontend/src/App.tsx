import { useEffect, useState } from 'react'
import { fetchDemographics, type Demographics } from './api'
import { BarList } from './components/BarList'

type State =
  | { status: 'loading' }
  | { status: 'error'; message: string }
  | { status: 'ready'; data: Demographics }

export default function App() {
  const [state, setState] = useState<State>({ status: 'loading' })

  useEffect(() => {
    const controller = new AbortController()
    fetchDemographics(controller.signal)
      .then((data) => setState({ status: 'ready', data }))
      .catch((err: unknown) => {
        if (controller.signal.aborted) return
        setState({ status: 'error', message: err instanceof Error ? err.message : String(err) })
      })
    return () => controller.abort()
  }, [])

  return (
    <main className="mx-auto max-w-5xl p-6">
      <h1 className="text-2xl font-semibold">Community tier</h1>
      {state.status === 'loading' && <p className="mt-4 text-neutral-500">Loading…</p>}
      {state.status === 'error' && <p className="mt-4 text-red-600">{state.message}</p>}
      {state.status === 'ready' && (
        <>
          <p className="mt-1 text-neutral-500">{state.data.member_count} members</p>
          <div className="mt-6 grid gap-6 md:grid-cols-2">
            {state.data.questions.map((q) => (
              <section key={q.question_id} className="rounded-lg border border-neutral-200 p-4 dark:border-neutral-800">
                <h2 className="font-medium">{q.label}</h2>
                <p className="mb-3 text-xs text-neutral-500">{q.answered} answered</p>
                <BarList values={q.values} total={q.answered} />
              </section>
            ))}
          </div>
        </>
      )}
    </main>
  )
}
