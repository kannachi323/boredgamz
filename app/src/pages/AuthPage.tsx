import { useState } from "react";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { safeNext } from "../lib/safeNext";
import { useAuthStore } from "../stores/useAuthStore";

export function AuthPage({ mode }: { mode: "login" | "signup" }) {
  const login = useAuthStore((state) => state.login);
  const signup = useAuthStore((state) => state.signup);
  const navigate = useNavigate();
  const [params] = useSearchParams();
  const next = safeNext(params.get("next"));
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  const isSignup = mode === "signup";

  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const data = new FormData(event.currentTarget);
    setBusy(true);
    setError("");
    try {
      if (isSignup) await signup(String(data.get("username")).trim(), String(data.get("email")).trim(), String(data.get("password")));
      else await login(String(data.get("email")).trim(), String(data.get("password")));
      navigate(next, { replace: true });
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Something went wrong.");
    } finally {
      setBusy(false);
    }
  }

  const other = `${isSignup ? "/login" : "/signup"}${next !== "/" ? `?next=${encodeURIComponent(next)}` : ""}`;
  return (
    <section className="auth">
      <div className="auth__card">
        <div className="auth__art" aria-hidden="true"><i /><i /><i /></div>
        <h1>{isSignup ? "Create your account" : "Welcome back"}</h1>
        <form onSubmit={submit}>
          {isSignup && (
            <label>
              Display name
              <input className="field" name="username" autoComplete="nickname" required minLength={2} maxLength={24} />
            </label>
          )}
          <label>
            Email
            <input className="field" name="email" type="email" autoComplete="email" required maxLength={254} />
          </label>
          <label>
            <span>Password {isSignup && <small>· 12 or more characters</small>}</span>
            <input className="field" name="password" type="password" autoComplete={isSignup ? "new-password" : "current-password"} required minLength={isSignup ? 12 : 1} maxLength={72} />
          </label>
          {error && <p className="alert" role="alert">{error}</p>}
          <button className="btn btn--primary btn--lg btn--block" disabled={busy} type="submit">
            {busy ? "One moment…" : isSignup ? "Create account" : "Sign in"}
          </button>
        </form>
        <p className="auth__switch">
          {isSignup ? "Have an account? " : "New here? "}
          <Link to={other}>{isSignup ? "Sign in" : "Create one"}</Link>
        </p>
      </div>
    </section>
  );
}
