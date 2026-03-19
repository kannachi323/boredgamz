
import { Dropdown } from "../Dropdown";
import { CgProfile } from "react-icons/cg";
import { useAuthStore } from "../../stores/useAuthStore";
import { useNavigate } from "react-router-dom";


export function NavBarV1() {
  const { isAuthenticated, user, logout } = useAuthStore();
  const navigate = useNavigate();

  return (
    <>
        <div className="flex justify-between items-center h-16 px-8">
          <a href="/" className="flex items-center space-x-2">
            <span className="text-2xl font-extrabold tracking-tight text-white">Bored<span className="text-gray-300">Gamz</span></span>
          </a>

          <nav className="flex flex-row justify-center items-center space-x-8">
            {['Games', 'Community', 'Donate'].map((item) => (
              <a key={item} href={`/${item.toLowerCase()}`} className="text-gray-300 hover:text-white transition duration-150 font-semibold">
                {item}
              </a>
            ))}

            <Dropdown
              align="center"
              openOnHover
              trigger={
                <span className="text-gray-300 hover:text-white transition duration-150 font-semibold">
                  Support
                </span>
              }
              items={[
                {
                  key: 'help',
                  content: 'Help',
                  onClick: () => {
                    window.location.href = '/help'
                  },
                },
                {
                  key: 'about',
                  content: 'About',
                  onClick: () => {
                    window.location.href = '/about'
                  },
                },
              ]}
            />

            <div className="flex flex-row justify-center items-center gap-2">
              {isAuthenticated ? (
                <Dropdown
                  align="right"
                  trigger={
                    <span
                      className="
                        flex items-center gap-2 px-3 py-1.5 
                        bg-gray-800 text-gray-200 
                        rounded-full border border-gray-700
                        hover:bg-gray-700 
                        transition
                      "
                    >
                      <CgProfile className="text-2xl" />
                      <span className="font-semibold text-sm">
                        {user?.username || "Player"}
                      </span>
                    </span>
                  }
                  items={[
                    {
                      key: 'logout',
                      content: 'Log out',
                      onClick: () => {
                        void logout(() => navigate('/'))
                      },
                    },
                  ]}
                />
              ) : (
                <div className="flex items-center gap-3">
                  <a
                    className="
                      px-4 py-1.5 
                      bg-gray-800 text-gray-100 
                      font-semibold rounded-md 
                      border border-gray-700
                      hover:bg-gray-700 
                      transition
                    "
                    href="/login"
                  >
                    Log In
                  </a>
                  <a
                    className="
                      px-4 py-1.5 
                      bg-white text-gray-900 
                      font-semibold rounded-md 
                      hover:bg-gray-200 
                      transition
                    "
                    href="/signup"
                  >
                    Sign Up
                  </a>
                </div>
              )}
            </div>
          </nav>
        </div>
    </>
  );
}