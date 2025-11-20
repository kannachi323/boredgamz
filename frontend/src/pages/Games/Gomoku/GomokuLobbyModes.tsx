import { Brain, Trophy, User, Settings } from "lucide-react";

export function GomokuLobbyModes({ onSelect, onOpen }: { onSelect: (mode: string) => void, onOpen: () => void }) {
  const modes = [
    { label: "Casual", value: "casual", icon: <User className="w-8 h-8 text-[#C3B299]" /> },
    { label: "Ranked", value: "ranked", icon: <Trophy className="w-8 h-8 text-[#C3B299]" /> },
    { label: "Custom", value: "custom", icon: <Settings className="w-8 h-8 text-[#C3B299]" /> },
    { label: "Bots",   value: "bots",   icon: <Brain className="w-8 h-8 text-[#C3B299]" /> },
  ];

  function handleSelect(value: string) {
    onSelect(value)
    onOpen()
  }

  return (
    <div className="flex flex-row justify-evenly w-full gap-6">
      {modes.map((mode) => (
        <div
          key={mode.value}
          onClick={() => handleSelect(mode.value)}
          className="
            bg-[#302e2e] 
            flex flex-col items-center justify-center
            gap-2 py-6 px-8
            rounded-lg border border-[#1b1918] 
            text-[#C3B299] text-2xl font-semibold
            hover:bg-[#524b4b] cursor-pointer 
            transition-all duration-200
          "
        >
          {mode.icon}
          {mode.label}
        </div>
      ))}
    </div>
  );
}
