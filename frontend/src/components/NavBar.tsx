
import { Dropdown } from "./Dropdown";
import { CgProfile } from "react-icons/cg";
import { useAuthStore } from "../stores/useAuthStore";


export function NavBar() {
  return (
    <div className="w-full h-[10vh] bg-[#252321] text-white px-6 flex items-center justify-between">
     
      <div className="flex items-center space-x-1 h-full w-1/3">
        <Dropdown url="/play/quick" label="Play" items={[
          <a href="/play/quick">Quick Play</a>,
          <a href="/play/custom">Create a custom game</a>,
        ]} />

        <Dropdown label="Learn" items={[
          <a>How to Play</a>,
          <a>Tips and tricks</a>,
        ]} />

        <Dropdown label="Community" items={[
          <a>Forums</a>,
          <a>Discord</a>,
          <a>GitHub</a>,
        ]} />
      </div>

      <div className="flex justify-center items-center w-1/3">
        <a href="/">
          <img src="/gopher.jpg" alt="logo" className="h-12" />
        </a>
      </div>

      <div className="flex items-center justify-end space-x-4 w-1/3">
        <UserAuth />
      </div>
    </div>
  );
}

function UserAuth() {
  const { isAuthenticated, user, logout } = useAuthStore();

  return (
    <>
      {isAuthenticated ? (
        <>
          <CgProfile className="text-4xl cursor-pointer" />
          <Dropdown label={user?.username || "annonymous"} items={[
          <button onClick={() =>logout()}>Log out</button>
          ]} />

        </>
      ) : (
        <>
          <a
            className="px-4 py-1 border border-white rounded cursor-pointer"
            href="/login"
          >
            Log In
          </a>
          <a
            className="px-4 py-1 border border-white rounded cursor-pointer"
            href="/signup"
          >
            Sign Up
          </a>
        </>
      )}
    </>
  );
}