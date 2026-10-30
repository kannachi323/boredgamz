import { Link } from "react-router-dom";

export function Wordmark() {
  return (
    <Link className="wordmark" to="/" aria-label="BoredGamz home">
      <span className="wordmark__mark" aria-hidden="true"><i /><i /><i /><i /></span>
      boredgamz
    </Link>
  );
}
