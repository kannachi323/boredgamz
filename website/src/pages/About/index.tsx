import { FooterV1 } from '../../features/Footer/FooterV1';

export default function About() {
  return (
    <div className="min-h-screen bg-zinc-900 text-gray-100">
      <section className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <h1 className="text-5xl font-extrabold text-white mb-6">About BoredGamz</h1>
        <p className="text-xl text-gray-300 max-w-3xl leading-relaxed">
          BoredGamz is built for people who love tabletop and social play. Our mission is simple:
          make it easier to discover great games, connect with other players, and spend more time
          actually playing.
        </p>
      </section>

      <section className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 pb-20">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <article className="rounded-2xl bg-zinc-800 border border-zinc-700 p-6 shadow-lg">
            <h2 className="text-2xl font-bold text-white mb-3">Discover</h2>
            <p className="text-gray-300">
              Browse curated game categories and find titles that match your group size, style, and mood.
            </p>
          </article>

          <article className="rounded-2xl bg-zinc-800 border border-zinc-700 p-6 shadow-lg">
            <h2 className="text-2xl font-bold text-white mb-3">Connect</h2>
            <p className="text-gray-300">
              Meet players in the community, exchange recommendations, and keep up with events and updates.
            </p>
          </article>

          <article className="rounded-2xl bg-zinc-800 border border-zinc-700 p-6 shadow-lg">
            <h2 className="text-2xl font-bold text-white mb-3">Play</h2>
            <p className="text-gray-300">
              Jump from discovery to game night quickly with a platform designed around clarity and speed.
            </p>
          </article>
        </div>
      </section>

      <FooterV1 />
    </div>
  );
}
