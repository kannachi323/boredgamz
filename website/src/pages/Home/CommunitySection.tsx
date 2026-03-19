

export function CommunitySection() {
  return (
    <section id="community" className="py-20 bg-gray-900">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 bg-gray-800 p-10 rounded-3xl shadow-2xl border border-indigo-700">
        <div className="md:flex md:items-center md:justify-between">
          <div className="mb-6 md:mb-0">
            <h2 className="text-3xl font-extrabold text-white">
              Ready to Roll the Dice?
            </h2>
            <p className="mt-2 text-xl text-gray-400">
              Join thousands of gamers. Share reviews, find local meetups, and track your collection.
            </p>
          </div>
          <div className="flex-shrink-0">
            <a href="#" className="inline-flex items-center justify-center px-6 py-3 border border-transparent text-base font-bold rounded-xl text-white bg-green-600 hover:bg-green-700 shadow-xl transition-transform duration-300 transform hover:-translate-y-1">
              Sign Up Now
            </a>
          </div>
        </div>
      </div>
    </section>
  )
  
}