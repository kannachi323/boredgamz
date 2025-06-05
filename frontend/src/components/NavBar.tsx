
export function NavBar() {
  return (
    <div id="modes" className="flex w-full h-[10vh] bg-blue-500 text-white">
      <section className="flex w-1/3 p-5"> 
        <a href="/play">Play</a>
        <div>
          <a>Quick Match</a>
          <a>Create a custom game</a>
          
        </div>
        <a href="/Learn">Learn</a>
        <div>
          <a>How to Play</a>
          <a>Tips and Tricks</a>
        </div>
      </section>

      <section className="flex w-1/3 bg-red-400 max-h-full justify-center items-center">
        <a href="/">
          <img src="/gopher.jpg" alt="logo" className="h-18" />
        </a>
      </section>

      <section className="flex w-1/3 items-center justify-end gap-2 p-5">
        <button>Sign Up</button>
        <button>Log In</button>
      </section>
    
    </div>
  )
}