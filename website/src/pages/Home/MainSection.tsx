
export function MainSection() {
  return (
    <section id="home" className="relative h-screen flex items-center justify-center text-center bg-gray-900 overflow-hidden">
      {/* Background Grid/Pattern */}
      <div className="absolute inset-0 z-0 opacity-10 pointer-events-none" style={{
      backgroundImage: 'radial-gradient(circle at 50% 50%, rgba(79, 70, 229, 0.1) 0%, transparent 70%)',
      }}></div>

      <div className="max-w-4xl mx-auto px-4 z-10">
        <h1 className="text-6xl sm:text-7xl lg:text-8xl font-black text-white mb-6 leading-tight">
            Your Next Game Night Starts
            <span className="block text-transparent bg-clip-text bg-gradient-to-r from-green-400 to-indigo-400 mt-2">Here</span>
        </h1>
        <p className="text-xl sm:text-2xl text-gray-300 mb-10 max-w-2xl mx-auto">
            Discover, review, and connect over the best tabletop and board games the world has to offer. Never be bored again.
        </p>
        <div className="flex justify-center space-x-4">
            <button className="px-8 py-3 text-lg font-bold text-gray-900 bg-green-400 rounded-xl shadow-lg hover:bg-green-300 transition-transform duration-300 transform hover:scale-105">
            Explore Games
            </button>
            <button className="px-8 py-3 text-lg font-bold text-white bg-indigo-600 rounded-xl shadow-lg border border-indigo-500 hover:bg-indigo-500 transition-transform duration-300 transform hover:scale-105">
            Join the Community
            </button>
        </div>
      </div>
    </section>
  )

}
  

