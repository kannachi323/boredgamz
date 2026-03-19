import { FooterV1 } from '../../features/Footer/FooterV1'

const tiers = [
  { name: 'Supporter', amount: '$3/mo', perks: 'Keeps servers running and development moving.' },
  { name: 'Backer', amount: '$8/mo', perks: 'Everything in Supporter plus early feature previews.' },
  { name: 'Patron', amount: '$15/mo', perks: 'Everything in Backer plus priority community feedback review.' },
]

export default function Donate() {
  return (
    <div className="min-h-screen bg-zinc-900 text-gray-100">
      <section className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <h1 className="text-5xl font-extrabold text-white mb-4">Donate</h1>
        <p className="text-lg text-gray-300 mb-10 max-w-3xl">
          If BoredGamz makes your game nights better, you can support hosting costs and future development.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {tiers.map((tier) => (
            <article key={tier.name} className="rounded-2xl bg-zinc-800 border border-zinc-700 p-6 shadow-lg flex flex-col">
              <h2 className="text-2xl font-bold text-white mb-1">{tier.name}</h2>
              <p className="text-3xl font-black text-white mb-4">{tier.amount}</p>
              <p className="text-gray-300 mb-6">{tier.perks}</p>
              <button className="mt-auto px-5 py-2.5 bg-white text-zinc-900 rounded-lg font-semibold hover:bg-zinc-200 transition">
                Choose {tier.name}
              </button>
            </article>
          ))}
        </div>
      </section>

      <FooterV1 />
    </div>
  )
}
