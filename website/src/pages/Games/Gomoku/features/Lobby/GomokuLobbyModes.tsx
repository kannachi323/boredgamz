import { useGomokuStore } from "@/stores/Gomoku/useGomokuStore";
import { GomokuModeModal } from "./GomokuModeModal";
import { Brain, Trophy, User, Settings } from "lucide-react";
import { useState } from "react";

export function GomokuLobbyModes() {
  const [showModeModal, setShowModeModal] = useState(false);
  const { lobbyRequest, setLobbyRequest } = useGomokuStore();
  const modes = [
    { mode: "casual", label: "Casual", icon: <User className="w-8 h-8 text-[#C3B299]" /> },
    { mode: "ranked", label: "Ranked",  icon: <Trophy className="w-8 h-8 text-[#C3B299]" /> },
    { mode: "custom", label: "Custom",  icon: <Settings className="w-8 h-8 text-[#C3B299]" /> },
    { mode: "bots", label: "Bots",    icon: <Brain className="w-8 h-8 text-[#C3B299]" /> },
  ];

  function handleModeSelect(mode: string) {
    setLobbyRequest({...lobbyRequest, data: { ...lobbyRequest.data, mode: mode}})
    if (mode === "custom" || mode === "bots") {
      setShowModeModal(true);
    }    
    
  }

  return (
    <>
      {modes.map(({ label, icon, mode }) => (
        <div
          key={label}
          onClick={() => handleModeSelect(mode)}
          className={`
            bg-[#302e2e] aspect-square 
            w-[128px] h-auto
            flex flex-col items-center justify-center
            gap-2 
            rounded-lg border border-[#1b1918] 
            text-[#C3B299] text-2xl font-semibold
             cursor-pointer 
            transition-all duration-200
            ${lobbyRequest.data.mode === mode
              ? "border-[#7DCFB6] shadow-[0_0_8px_#7DCFB6]"
              : "border-[#C3B299] hover:border-[#7DCFB6] hover:bg-[#524b4b]"
            }
          `}
        >
          {icon}
          {label}
        </div>
      ))}

      
      {/* Mode popup */}
      {showModeModal && (
        <GomokuModeModal
          mode={lobbyRequest.data.mode}
          onClose={() => setShowModeModal(false)}
        />
      )}
    </>
  );
}
