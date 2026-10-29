import { LogOut } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { useAuthStore } from "../stores/useAuthStore";
import { ThemeToggle } from "./ThemeToggle";
import { Wordmark } from "./Wordmark";

export function TopBar() {
  const status = useAuthStore((state) => state.status);
  const user = useAuthStore((state) => state.user);
  const location = useLocation();
  const next = encodeURIComponent(location.pathname + location.search);
  return (
    <header className="topbar">
      <div className="topbar__inner">
        <Wordmark />
        <span className="topbar__spacer" />
        <ThemeToggle />
        {status === "authenticated" && user ? (
          <AccountMenu name={user.username} />
        ) : status === "anonymous" ? (
          <>
            <Link className="btn btn--ghost hide-sm" to={`/login?next=${next}`}>Sign in</Link>
            <Link className="btn btn--ink" to={`/signup?next=${next}`}>Create account</Link>
          </>
        ) : null}
      </div>
    </header>
  );
}

function AccountMenu({ name }: { name: string }) {
  const [open, setOpen] = useState(false);
  const logout = useAuthStore((state) => state.logout);
  const ref = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!open) return;
    const onPointer = (event: PointerEvent) => {
      if (!ref.current?.contains(event.target as Node)) setOpen(false);
    };
    const onKey = (event: KeyboardEvent) => event.key === "Escape" && setOpen(false);
    document.addEventListener("pointerdown", onPointer);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("pointerdown", onPointer);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);
  return (
    <div className="menu" ref={ref}>
      <button className="btn btn--ghost" type="button" aria-label={`Account: ${name}`} aria-expanded={open} aria-haspopup="menu" onClick={() => setOpen((value) => !value)}>
        <span className="avatar" aria-hidden="true">{name.slice(0, 1)}</span>
        <span className="hide-sm">{name}</span>
      </button>
      {open && (
        <div className="menu__panel" role="menu">
          <p>Signed in as {name}</p>
          <button className="btn btn--ghost" role="menuitem" type="button" onClick={() => { setOpen(false); void logout(); }}>
            <LogOut size={16} /> Sign out
          </button>
        </div>
      )}
    </div>
  );
}
