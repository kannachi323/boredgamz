
export function GomokuLobbyWaiting({ onCancel } : { onCancel : () => void }) {
  return (

    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70">

      <div className="bg-[#433d3a] p-10 rounded-xl shadow-2xl flex flex-col items-center gap-6 max-w-sm w-full">
        <div className="flex flex-col items-center gap-4">
          <svg className="animate-spin h-10 w-10 text-[#C3B299]" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <h2 className="text-2xl text-[#C3B299] font-bold">Searching for Opponent...</h2>
        </div>

        <button
          className="px-8 py-2 bg-[#d7c9b8] text-[#433d3a] font-bold rounded-lg hover:bg-[#C3B299] transition"
          onClick={onCancel}
        >
          Cancel
        </button>
      </div>
    </div>
  );
}