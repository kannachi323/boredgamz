
import { Dropdown } from "../Dropdown";
import { CgProfile } from "react-icons/cg";
import { useAuthStore } from "../../stores/useAuthStore";
import { useNavigate } from "react-router-dom";


export function NavBarV1() {
  const { isAuthenticated, user, logout } = useAuthStore();
  const navigate = useNavigate();

  return (

    <header className="sticky top-0 z-50 bg-gray-900/90 backdrop-blur-sm shadow-xl border-b border-gray-700">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          <a href="/" className="flex items-center space-x-2">
            <span className="text-2xl font-extrabold tracking-tight text-white">Bored<span className="text-green-400">Gamz</span></span>
          </a>

          <nav className="flex flex-row justify-center items-center space-x-8">
            {['Games', 'Community', 'About'].map((item) => (
              <a key={item} href={`/${item.toLowerCase()}`} className="text-gray-300 hover:text-green-400 transition duration-150 font-semibold">
                {item}
              </a>
            ))}
            <div className="flex flex-row justify-center items-center space-x-2">
              {isAuthenticated ? (
                <>
                  <CgProfile className="text-4xl cursor-pointer" />
                  <Dropdown label={user?.username || "annonymous"} items={[
                  <button onClick={() =>logout(() => navigate('/'))}>Log out</button>
                  ]} />

                </>
                ) : (
                <>
                  <a
                    className="px-4 py-1 bg-[#313d51] text-[#d1d7e3] font-semibold rounded cursor-pointer hover:bg-[#43536d] transition"
                    href="/login"
                  >
                    Log In
                  </a>
                  <a
                    className="px-4 py-1 bg-[#05df71] text-[#002123] font-semibold rounded cursor-pointer hover:bg-[#1bfe8c] transition"
                    href="/signup"
                  >
                    Sign Up
                  </a>
                </>
              )}
            </div>
          </nav>
        </div>
      </div>
    </header>
  );
}