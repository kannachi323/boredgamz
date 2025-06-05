import { NavBar } from "./components/NavBar";

interface LayoutProps {
  children : React.ReactNode
}

export function AppLayout({children} : LayoutProps) {
  return (
    <>
      <NavBar />
      {children}
    </>
  )
}