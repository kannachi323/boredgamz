import { FooterV1 } from '../../features/Footer/FooterV1'

const faqs = [
  {
    q: 'How do I start playing?',
    a: 'Pick a game from the home page and click Play. You can jump in as a guest or sign up to keep your progress.',
  },
  {
    q: 'Do I need an account?',
    a: 'No. You can play as a guest. Creating an account helps you save game history and personalize your profile.',
  },
  {
    q: 'My connection dropped in a match. What now?',
    a: 'Reconnect from the game screen. The app attempts to restore your match when possible.',
  },
]

export default function Help() {
  return (
    <div className="min-h-screen bg-zinc-900 text-gray-100">
      <section className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <h1 className="text-5xl font-extrabold text-white mb-4">Help</h1>
        <p className="text-lg text-gray-300 mb-10">
          Quick answers and support for getting the most out of BoredGamz.
        </p>

        <div className="space-y-5">
          {faqs.map((faq) => (
            <article key={faq.q} className="rounded-2xl bg-zinc-800 border border-zinc-700 p-6 shadow-lg">
              <h2 className="text-xl font-bold text-white mb-2">{faq.q}</h2>
              <p className="text-gray-300">{faq.a}</p>
            </article>
          ))}
        </div>
      </section>

      <FooterV1 />
    </div>
  )
}
