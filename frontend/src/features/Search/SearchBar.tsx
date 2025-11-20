import { useRef, useState } from "react";
import { SearchAdvanced } from "./SearchAdvanced";
import { ChevronDown, ChevronUp, Search } from "lucide-react";
import { useClickOutside } from "../../hooks/useClickOutside";

export function SearchBar() {
  const wrapperRef = useRef<HTMLDivElement>(null);
  const textRef = useRef<HTMLTextAreaElement>(null);
  const [showAdvanced, setShowAdvanced] = useState(false);

  const handleInput = () => {
    const textarea = textRef.current;
    if (!textarea) return;

    textarea.style.height = "auto";
    textarea.style.height =
      textarea.value.trim() === "" ? "3rem" : textarea.scrollHeight + "px";
  };

  // close when clicking outside everything (search + dropdown)
  useClickOutside(wrapperRef, () => setShowAdvanced(false));

  return (
    <div className="flex justify-center mb-12">
      <div ref={wrapperRef} className="relative w-full max-w-xl h-[48px]">

        {/* Search Bar */}
       
          <textarea
            ref={textRef}
            onInput={handleInput}
            placeholder="Search games..."
            rows={1}
            className="
              w-full
              py-3 pl-12 pr-16
              max-h-64
              bg-white 
              text-gray-900 
              rounded-full 
              shadow-lg 
              outline-none 
              
              resize-none
              placeholder-gray-500
              transition-all duration-200
              overflow-hidden
            "
          />

          {/* FIXED: Center vertically */}
          <Search className="absolute top-1/2 -translate-y-1/2 left-3 z-[50] text-gray-900" size={20} />

          <button
            onClick={() => setShowAdvanced((s) => !s)}
            className="
              absolute right-3 top-1/2 -translate-y-1/2 
              text-black px-3 py-1 rounded-full
              transition
              z-20
            "
          >
            {showAdvanced ? (
              <ChevronUp className="w-6 h-6" />
            ) : (
              <ChevronDown className="w-6 h-6" />
            )}
          </button>
          
          {showAdvanced && (
            <div className="absolute left-0 w-full mt-2 z-50">
              <SearchAdvanced className="w-full" />
            </div>
          )}
        </div>
      </div>
    
  );
}
