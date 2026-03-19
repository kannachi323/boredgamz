export function GomokuModeModal({
  mode,
  onClose
}: {
  mode: string;
  onClose: () => void;
}) {

  function renderContent() {
    console.log(mode);
    switch (mode) {
      case "custom":
        return <BotModeContent onClose={onClose} />;
      case "bots":
        return <BotModeContent onClose={onClose} />;
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

function BotModeContent({ onClose }: { onClose: () => void }) {
  const bots = [
    { name: "Beginner Bot", strength: 1 },
    { name: "Intermediate Bot", strength: 2 },
    { name: "Advanced Bot", strength: 3 },
    { name: "Master Bot", strength: 4 },
  ];

  return (
    <div className="flex flex-col gap-4">
      <h2 className="text-2xl font-bold">Choose Bot Difficulty</h2>

      <div className="max-h-64 overflow-y-auto pr-2 flex flex-col gap-3">
        {bots.map((bot) => (
          <div
            key={bot.name}
            className="
              bg-[#3b3735] flex flex-col items-center justify-center
              gap-2 p-4 rounded-lg border border-[#1b1918]
              hover:bg-[#524b4b] cursor-pointer transition-all duration-200
            "
            onClick={onClose}
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
