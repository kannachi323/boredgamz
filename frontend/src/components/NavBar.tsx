import { Dropdown } from "./Dropdown";

export function NavBar() {
  return (
    <div className="w-full h-[10vh] bg-[#363430] text-white px-6 flex items-center justify-between">
     
      <div className="flex items-center space-x-1 h-full w-1/3">
        <Dropdown url="/play" label="Play" items={[
          <a>Quick Play</a>,
          <a>Create a custom game</a>,
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
        <button className="px-4 py-1 border border-white rounded">Log In</button>
        <button className="px-4 py-1 bg-white text-[#363430] rounded">Sign Up</button>
      </div>
    </div>
  );
}
