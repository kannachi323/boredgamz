export function GomokuModeModal({
  mode,
  onClose,
  selectedBotDifficulty,
  onSelectBotDifficulty,
}: {
  mode: string;
  selectedBotDifficulty: "beginner" | "intermediate" | "advanced";
  onSelectBotDifficulty: (difficulty: "beginner" | "intermediate" | "advanced") => void;
  onClose: () => void;
}) {

  function renderContent() {
    console.log(mode);
    switch (mode) {
      case "custom":
        return <CustomModeContent onClose={onClose} />;
      case "bots":
        return (
          <BotModeContent
            onClose={onClose}
            selectedBotDifficulty={selectedBotDifficulty}
            onSelectBotDifficulty={onSelectBotDifficulty}
          />
        );
      default:
        return null;
    }
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex justify-center items-center z-50">
      <div className="bg-[#2c2826] p-8 rounded-xl text-[#C3B299] min-w-[400px] shadow-xl">
        {renderContent()}
      </div>
    </div>
  );
}

function CustomModeContent({ onClose }: { onClose: () => void }) {
  return (
    <div className="flex flex-col gap-4">
      <h2 className="text-2xl font-bold">Custom Mode</h2>
      <p className="text-[#d8c8b0]">Custom mode is coming soon.</p>
      <button
        onClick={onClose}
        className="bg-[#524b4b] px-4 py-2 rounded-lg text-lg hover:bg-[#6a605f]"
      >
        Close
      </button>
    </div>
  );
}

function BotModeContent({
  onClose,
  selectedBotDifficulty,
  onSelectBotDifficulty,
}: {
  onClose: () => void;
  selectedBotDifficulty: "beginner" | "intermediate" | "advanced";
  onSelectBotDifficulty: (difficulty: "beginner" | "intermediate" | "advanced") => void;
}) {
  const bots = [
    { key: "beginner", name: "Beginner", strength: 1 },
    { key: "intermediate", name: "Intermediate", strength: 2 },
    { key: "advanced", name: "Advanced", strength: 3 },
  ] as const;

  function handleSelectBot(difficulty: "beginner" | "intermediate" | "advanced") {
    onSelectBotDifficulty(difficulty);
    onClose();
  }

  function formatDifficulty(difficulty: "beginner" | "intermediate" | "advanced") {
    switch (difficulty) {
      case "beginner":
        return "Beginner";
      case "intermediate":
        return "Intermediate";
      case "advanced":
        return "Advanced";
    }
  }

  return (
    <div className="flex flex-col gap-4">
      <h2 className="text-2xl font-bold">Choose Bot Difficulty</h2>
      <p className="text-[#d8c8b0]">
        Selected: <span className="font-bold text-white">{formatDifficulty(selectedBotDifficulty)}</span>
      </p>

      <div className="overflow-y-auto pr-2 flex flex-col gap-3">
        {bots.map((bot) => (
          <div
            key={bot.key}
            className="
              bg-[#3b3735] flex flex-col items-center justify-center
              gap-2 p-4 rounded-lg border transition-all duration-200
              hover:bg-[#524b4b] cursor-pointer
              border-[#C3B299]
            "
            onClick={() => handleSelectBot(bot.key)}
          >
            <p className="text-xl">{bot.name}</p>
            <div className="flex gap-1">
              {Array(bot.strength).fill(0).map((_, i) => (
                <div key={i} className="w-3 h-3 bg-[#C3B299] rounded-full" />
              ))}
            </div>
          </div>
        ))}
      </div>

      <button
        onClick={onClose}
        className="bg-[#524b4b] px-4 py-2 rounded-lg text-lg hover:bg-[#6a605f]"
      >
        Close
      </button>
    </div>
  );
}
