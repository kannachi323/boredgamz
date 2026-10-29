import { ArrowLeft, Dices } from 'lucide-react'
import { Link, useRouteError } from 'react-router-dom'

export function ErrorPage() {
  const error = useRouteError()
  const message = error instanceof Error ? error.message : 'That table could not be found.'
  return (
    <main className="error-page">
      <Dices size={42} aria-hidden="true" />
      <p className="eyebrow">Wrong turn</p>
      <h1>We lost the board.</h1>
      <p>{message}</p>
      <Link className="button button--primary" to="/"><ArrowLeft size={18} /> Back home</Link>
    </main>
  )
}
