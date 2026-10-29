import { useEffect } from 'react'
import { Outlet } from 'react-router-dom'
import { AppLayout } from './Layout'
import { useAuthStore } from './stores/useAuthStore'

function App() {
  const initialize = useAuthStore((state) => state.initialize)
  const authStatus = useAuthStore((state) => state.status)
  useEffect(() => { void initialize() }, [initialize])
  if (authStatus === 'checking') {
    return <AppLayout><div className="account-loading" role="status"><span /><p>Restoring your table profile…</p></div></AppLayout>
  }
  return (
    <AppLayout>
      <Outlet />
    </AppLayout>
  )
}

export default App
