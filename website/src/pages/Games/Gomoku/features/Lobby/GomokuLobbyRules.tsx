import { useGomokuStore } from "@/stores/Gomoku/useGomokuStore";

export function GomokuLobbyRules() {
  const { lobbyRequest, setLobbyRequest } = useGomokuStore();

  function setOpeningRule(openingRule: "freestyle" | "standard" | "renju") {
    setLobbyRequest({ ...lobbyRequest, data: { ...lobbyRequest.data, openingRule } });
  }

  function setSwapRuleEnabled(value: boolean) {
    setLobbyRequest({ ...lobbyRequest, data: { ...lobbyRequest.data, swapRuleEnabled: value } });
  }

  function setFirstMoveCenterEnabled(value: boolean) {
    setLobbyRequest({ ...lobbyRequest, data: { ...lobbyRequest.data, firstMoveCenterEnabled: value } });
  }

  return (
    <div className="flex flex-col items-center gap-4">
      <div className="flex flex-row items-center justify-center gap-5 flex-wrap">
        <p className="text-2xl text-[#C3B299] font-bold">Opening Rule:</p>

        {[
          { label: "Freestyle", openingRule: "freestyle" as const },
          { label: "Standard", openingRule: "standard" as const },
          { label: "Renju", openingRule: "renju" as const },
        ].map(({ label, openingRule }) => (
          <button
            key={openingRule}
            className={`
              text-xl font-semibold cursor-pointer rounded-xl px-5 py-2 transition-all duration-300
              border-2 text-[#C3B299] bg-[#302e2e]
              ${
                lobbyRequest.data.openingRule === openingRule
                  ? "border-white shadow-[0_0_8px_#ffffff]"
                  : "border-[#C3B299] hover:border-white hover:bg-[#524b4b]"
              }
            `}
            onClick={() => setOpeningRule(openingRule)}
          >
            {label}
          </button>
        ))}
      </div>

      <div className="flex flex-row items-center justify-center gap-6 flex-wrap">
        <label className="flex items-center gap-3 text-[#C3B299] text-lg font-semibold cursor-pointer">
          <input
            type="checkbox"
            className="h-5 w-5 accent-[#C3B299]"
            checked={!!lobbyRequest.data.swapRuleEnabled}
            onChange={(e) => setSwapRuleEnabled(e.target.checked)}
          />
          Swap Rule
        </label>

        <label className="flex items-center gap-3 text-[#C3B299] text-lg font-semibold cursor-pointer">
          <input
            type="checkbox"
            className="h-5 w-5 accent-[#C3B299]"
            checked={!!lobbyRequest.data.firstMoveCenterEnabled}
            onChange={(e) => setFirstMoveCenterEnabled(e.target.checked)}
          />
          First move must be center
        </label>
      </div>
    </div>
  );
}
