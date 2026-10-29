import { ArrowRight, LockKeyhole } from "lucide-react";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuthStore } from "../../stores/useAuthStore";

export default function Login() {
  const login = useAuthStore((state) => state.login);
  const navigate = useNavigate();
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setBusy(true);
    setError("");
    const data = new FormData(event.currentTarget);
    try {
      await login(String(data.get("email")), String(data.get("password")));
      navigate("/games");
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Sign in failed.");
    } finally {
      setBusy(false);
    }
  }

  return <section className="auth-page">
    <div className="auth-card">
      <div className="auth-card__intro"><span><LockKeyhole size={15} /> PLAYER ACCOUNT</span><h1>Welcome back<br />to the table.</h1><p>Your account name follows you into every live room while the games remain server-authoritative.</p></div>
      <form onSubmit={submit}>
        <header><small>Sign in</small><h2>Continue playing</h2></header>
        <label><span>Email</span><input name="email" type="email" autoComplete="email" required /></label>
        <label><span>Password</span><input name="password" type="password" autoComplete="current-password" required minLength={12} maxLength={128} /></label>
        {error && <p className="auth-error" role="alert">{error}</p>}
        <button className="button button--primary" disabled={busy} type="submit">{busy ? "Signing in…" : "Sign in"}<ArrowRight size={17} /></button>
        <p className="auth-switch">New to BoredGamz? <Link to="/signup">Create an account</Link></p>
      </form>
    </div>
  </section>;
}
