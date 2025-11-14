import { NavBarV1 } from "./components/NavBar/NavBarV1";

interface LayoutProps {
  children : React.ReactNode
}

export function AppLayout({children} : LayoutProps) {
  return (
    <>
      <NavBarV1 />
      {children}
    </>
  )
}