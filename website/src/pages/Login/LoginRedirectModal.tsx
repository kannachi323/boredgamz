import { useNavigate } from "react-router-dom";

export function LoginRedirectModal() {
  const navigate = useNavigate();

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/50 z-50">
      <div className="bg-[#262322] text-white rounded-md p-10 w-1/3 flex flex-col items-center gap-6 border border-[#454340] shadow-lg">
        <h2 className="text-2xl font-bold text-center">PlayGomoku</h2>
        <p className="text-center text-gray-300">
          Hey there! You must have an account to play.
        </p>
        <button
          onClick={() => navigate("/login")}
          className="bg-[#363430] hover:bg-[#454340] transition-colors duration-300 text-white rounded-sm px-6 py-3 w-full"
        >
          Log in
        </button>
         <button
          onClick={() => navigate("/signup")}
          className="bg-[#363430] hover:bg-[#454340] transition-colors duration-300 text-white rounded-sm px-6 py-3 w-full"
        >
          Sign up
        </button>
        <p className="text-blue-300 underline cursor-pointer" onClick={() => navigate("/")}>
          Return to Home
        </p>
      </div>
    </div>
  );
}
