
export function FooterV1() {
  return (
    <footer className="bg-gray-900 border-t border-gray-700">
      <div className="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-2 gap-8 md:grid-cols-4">
          <div>
            <h3 className="text-sm font-semibold text-gray-400 tracking-wider uppercase">Product</h3>
            <ul className="mt-4 space-y-4">
              {['All Games', 'Top Rated', 'New Releases', 'Board Game News'].map((item) => (
                <li key={item}>
                  <a href="#" className="text-base text-gray-500 hover:text-green-400 transition">{item}</a>
                </li>
              ))}
            </ul>
          </div>
          <div>
            <h3 className="text-sm font-semibold text-gray-400 tracking-wider uppercase">Company</h3>
            <ul className="mt-4 space-y-4">
              {['About Us', 'Careers', 'Press', 'Contact'].map((item) => (
                <li key={item}>
                  <a href="#" className="text-base text-gray-500 hover:text-green-400 transition">{item}</a>
                </li>
              ))}
            </ul>
          </div>
          <div>
            <h3 className="text-sm font-semibold text-gray-400 tracking-wider uppercase">Support</h3>
            <ul className="mt-4 space-y-4">
              {['Help Center', 'API', 'Terms of Service', 'Privacy Policy'].map((item) => (
                <li key={item}>
                  <a href="#" className="text-base text-gray-500 hover:text-green-400 transition">{item}</a>
                </li>
              ))}
            </ul>
          </div>
          <div className="col-span-2 md:col-span-1">
            <h3 className="text-sm font-semibold text-gray-400 tracking-wider uppercase">Follow Us</h3>
            <div className="mt-4 flex space-x-6">
              <a href="#" className="text-gray-400 hover:text-green-400 transition">
                {/* Dummy Icon for Twitter/X */}
                <svg className="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path d="M12 2C6.477 2 2 6.477 2 12s4.477 10 10 10 10-4.477 10-10S17.523 2 12 2zm3.707 10.707a.5.5 0 01-.707 0L12 9.414l-3.09 3.293a.5.5 0 01-.707-.707l3.5-3.5a.5.5 0 01.707 0l3.5 3.5a.5.5 0 010 .707z" /></svg>
              </a>
              <a href="#" className="text-gray-400 hover:text-green-400 transition">
                {/* Dummy Icon for Instagram */}
                <svg className="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path d="M12 2C6.477 2 2 6.477 2 12s4.477 10 10 10 10-4.477 10-10S17.523 2 12 2zm4.5 9.176l-4.5 2.5v-5l4.5 2.5z" /></svg>
              </a>
            </div>
          </div>
        </div>
        <div className="mt-12 border-t border-gray-700 pt-8">
          <p className="text-base text-gray-500 xl:text-center">
            &copy; {new Date().getFullYear()} BoredGamz, Inc. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  )
}