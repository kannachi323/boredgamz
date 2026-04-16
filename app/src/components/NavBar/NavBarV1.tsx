
import { Dropdown } from "../Dropdown";
import { CgProfile } from "react-icons/cg";
import { useAuthStore } from "../../stores/useAuthStore";
import { useNavigate } from "react-router-dom";
import { SITE_LOCKDOWN_ENABLED } from "../../config/siteAccess";


export function NavBarV1() {
  const { isAuthenticated, user, logout } = useAuthStore();
  const navigate = useNavigate();
  const navItems = SITE_LOCKDOWN_ENABLED ? [] : ['Games', 'Community', 'About'];
  const showAuthActions = !SITE_LOCKDOWN_ENABLED;
  const authControls = isAuthenticated ? (
    <div className="relative">
      <button
        className="
          flex items-center gap-2 px-3 py-1.5 
          bg-[#2a2f37] text-gray-200 
          rounded-full border border-[#3a3f47]
          hover:bg-[#333842] 
          transition
        "
      >
        <CgProfile className="text-2xl" />
        <span className="font-semibold text-sm">
          {user?.username || "Player"}
        </span>
      </button>

      <div className="absolute right-0 mt-2">
        <Dropdown
          label=""
          items={[
            <button
              className="text-left w-full px-3 py-2 hover:bg-gray-100 transition"
              onClick={() => logout(() => navigate("/"))}
            >
              Log out
            </button>,
          ]}
        />
      </div>
    </div>
  ) : (
    <div className="flex items-center gap-3">
      <a
        className="
          px-4 py-1.5 
          bg-[#313d51] text-[#d1d7e3] 
          font-semibold rounded-md 
          hover:bg-[#43536d] 
          transition
        "
        href="/login"
      >
        Log In
      </a>
      <a
        className="
          px-4 py-1.5 
          bg-[#05df71] text-[#002123] 
          font-semibold rounded-md 
          hover:bg-[#1bfe8c] 
          transition
        "
        href="/signup"
      >
        Sign Up
      </a>
    </div>
  );

  return (
    <>
        <div className="flex justify-between items-center h-16 px-8">
          <a href="/" className="flex items-center space-x-2">
            <span className="text-2xl font-extrabold tracking-tight text-white">Bored<span className="text-green-400">Gamz</span></span>
          </a>

          <nav className="flex flex-row justify-center items-center space-x-8">
            {navItems.map((item) => (
              <a key={item} href={`/${item.toLowerCase()}`} className="text-gray-300 hover:text-green-400 transition duration-150 font-semibold">
                {item}
              </a>
            ))}
            <div className="flex flex-row justify-center items-center gap-2">
              {SITE_LOCKDOWN_ENABLED ? (
                <span className="rounded-full border border-green-500/40 bg-green-500/10 px-3 py-1 text-sm font-semibold text-green-300">
                  Home page only
                </span>
              ) : null}
              {showAuthActions ? authControls : null}
            </div>
          </nav>
        </div>
    </>
  );
}
