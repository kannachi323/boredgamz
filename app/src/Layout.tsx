import type { ReactNode } from "react";
import { Link } from "react-router-dom";
import { TopBar } from "./components/TopBar";
import { useShell } from "./stores/useShell";

export function AppLayout({ children }: { children: ReactNode }) {
  const immersive = useShell((shell) => shell.immersive);
  return (
    <>
      <a className="skip-link" href="#main-content">Skip to content</a>
      {!immersive && <TopBar />}
      <main id="main-content" className={immersive ? undefined : "site-main"}>{children}</main>
      {!immersive && (
        <footer className="footer">
          <span>© {new Date().getFullYear()} BoredGamz</span>
          <Link to="/">Games</Link>
        </footer>
      )}
    </>
  );
}
