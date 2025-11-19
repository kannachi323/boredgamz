
import { SearchCategory } from "./SearchCategory";

interface SearchAdvancedProps {
  className?: string;
}
interface SearchAdvancedProps {
  className?: string;
}

export function SearchAdvanced({ className }: SearchAdvancedProps) {
  return (
    <div
      className={`
        bg-gray-900 text-white
        rounded-2xl shadow-2xl p-6 
        w-full
        animate-fade-slide
        ${className}
      `}
    >
      <SearchCategory title="Board Size" options={["3x3", "5x5", "7x7", "Custom"]} />
      <SearchCategory title="Difficulty" options={["Easy", "Medium", "Hard"]} />
      <SearchCategory title="Game Type" options={["Board", "Card", "Puzzle"]} />
      <SearchCategory title="Players" options={["1 Player", "2 Players", "Online"]} />
    </div>
  );
}
