import { ArrowRight, UserRoundPlus } from "lucide-react";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuthStore } from "../../stores/useAuthStore";

export default function SignUp() {
  const signup = useAuthStore((state) => state.signup);
  const navigate = useNavigate();
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setBusy(true);
    setError("");
    const data = new FormData(event.currentTarget);
    try {
      await signup(String(data.get("username")), String(data.get("email")), String(data.get("password")));
      navigate("/games");
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Account creation failed.");
    } finally {
      setBusy(false);
    }
  }

  return <section className="auth-page auth-page--signup">
    <div className="auth-card">
      <div className="auth-card__intro"><span><UserRoundPlus size={15} /> JOIN THE CLUB</span><h1>Claim your<br />seat.</h1><p>Create one account for your table name, secure sessions, and future match history—without purchases or paywalls.</p></div>
      <form onSubmit={submit}>
        <header><small>Create account</small><h2>Meet your next rival</h2></header>
        <label><span>Display name</span><input name="username" autoComplete="nickname" required minLength={2} maxLength={24} /></label>
        <label><span>Email</span><input name="email" type="email" autoComplete="email" required /></label>
        <label><span>Password</span><input name="password" type="password" autoComplete="new-password" required minLength={12} maxLength={128} /><small>Use at least 12 characters.</small></label>
        {error && <p className="auth-error" role="alert">{error}</p>}
        <button className="button button--primary" disabled={busy} type="submit">{busy ? "Creating account…" : "Create account"}<ArrowRight size={17} /></button>
        <p className="auth-switch">Already registered? <Link to="/login">Sign in</Link></p>
      </form>
    </div>
  </section>;
}
