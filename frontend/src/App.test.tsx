import { render, screen } from '@testing-library/react'
import { afterEach, expect, test, vi } from 'vitest'
import App from './App'

afterEach(() => {
  vi.unstubAllGlobals()
})

test('shows the member count and each question', async () => {
  vi.stubGlobal(
    'fetch',
    vi.fn().mockResolvedValue(
      new Response(
        JSON.stringify({
          member_count: 3,
          questions: [
            {
              question_id: 'q1',
              label: 'Graduation year',
              answered: 3,
              values: [
                { value: '2027', count: 2 },
                { value: '2028', count: 1 },
              ],
            },
          ],
        }),
      ),
    ),
  )

  render(<App />)

  expect(await screen.findByText('3 members')).toBeInTheDocument()
  expect(screen.getByText('Graduation year')).toBeInTheDocument()
  expect(screen.getByText('2027')).toBeInTheDocument()
})

test('shows an error when the API fails', async () => {
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('', { status: 502 })))

  render(<App />)

  expect(await screen.findByText('Loading demographics failed (502)')).toBeInTheDocument()
})
