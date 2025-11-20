import { NavBarV1 } from "./components/NavBar/NavBarV1";

interface LayoutProps {
  children : React.ReactNode
}

export function AppLayout({children} : LayoutProps) {
  return (
    <>
      <header className="sticky top-0 z-50 bg-gray-900/90 backdrop-blur-sm shadow-xl border-b 
        border-gray-700 h-[8vh] max-h-[8vh] w-screen max-w-screen"
      >
        <NavBarV1 />
      </header>
      <div className="h-[92vh] max-h-[92vh] w-screen max-w-screen overflow-hidden">
        {children}
      </div>
    </>
  )
}