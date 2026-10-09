import type { ValueCount } from '../api'

// A horizontal bar per value, scaled to the largest count. Enough for a first dashboard; swap in a
// chart library when the dashboards need more.
export function BarList({ values, total }: { values: ValueCount[]; total: number }) {
  const max = Math.max(1, ...values.map((v) => v.count))
  return (
    <ul className="space-y-2">
      {values.map((v) => (
        <li key={v.value} className="text-sm">
          <div className="flex justify-between gap-4">
            <span className="truncate">{v.value}</span>
            <span className="tabular-nums text-neutral-500">
              {v.count}
              {total > 0 && ` · ${Math.round((v.count / total) * 100)}%`}
            </span>
          </div>
          <div className="mt-1 h-2 rounded bg-neutral-100 dark:bg-neutral-800">
            <div
              className="h-2 rounded bg-indigo-500"
              style={{ width: `${(v.count / max) * 100}%` }}
            />
          </div>
        </li>
      ))}
    </ul>
  )
}
