
export function MainSection() {
  return (
    <section id="home" className="relative bg-black/70">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="rounded-3xl shadow-2xl px-6 py-14 sm:px-10 sm:py-16 text-center">
          <h1 className="text-white text-5xl sm:text-6xl lg:text-7xl font-bold mb-5 tracking-tight leading-tight">
            Welcome to{' '}
            <span className="bg-linear-to-r from-white via-zinc-200 to-zinc-400 bg-clip-text text-transparent">
              BoredGamz
            </span>
          </h1>
          <p className="text-zinc-200 text-lg sm:text-2xl mb-10 max-w-3xl mx-auto">
            Find your favorite board game and jump straight into play.
          </p>

          <div className="flex flex-wrap items-center justify-center gap-4">
            <a
              href="/signup"
              className="px-8 py-3 text-lg font-bold text-zinc-900 bg-white hover:bg-white/90 rounded-xl shadow-lg"
            >
              Get Started
            </a>
            <a
              href="/games"
              className="px-8 py-3 text-lg font-bold text-zinc-900 bg-gray-200 hover:bg-gray-200/90 rounded-xl shadow-lg"
            >
              Explore Games
            </a>
          </div>
        </div>
      </div>
    </section>
  )
}
