import { Link, isRouteErrorResponse, useRouteError } from "react-router-dom";
import { BoardArt } from "../games/BoardArt";

export function ErrorPage() {
  const error = useRouteError();
  const notFound = !error || (isRouteErrorResponse(error) && error.status === 404);
  return (
    <section className="error-page">
      <div className="waiting__art" style={{ animation: "none" }}><BoardArt slug="tictactoe" /></div>
      <h1>{notFound ? "No table here." : "Something broke."}</h1>
      <p>{notFound ? "That page doesn't exist." : "Reload the page or head back to the games."}</p>
      <Link className="btn btn--ink" to="/">All games</Link>
    </section>
  );
}
