import { useRef, useState } from "react";

import { SearchAdvanced } from "./SearchAdvanced";
import { ChevronDown, ChevronUp } from "lucide-react";

export function SearchBar() {
  const textRef = useRef<HTMLTextAreaElement>(null);
  const [showAdvanced, setShowAdvanced] = useState(false);

  const handleInput = () => {
    const textarea = textRef.current;
    if (!textarea) return;

    textarea.style.height = "auto";
    if (textarea.value.trim() === "") {
      textarea.style.height = "3rem";
      return;
    }
    textarea.style.height = textarea.scrollHeight + "px";
  };

  return (
    <div className="flex justify-center mb-12">
      <div className="relative w-full max-w-xl">

        <textarea
          ref={textRef}
          onInput={handleInput}
          placeholder="Search games..."
          rows={1}
          style={{ height: "3rem" }}
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
            focus:ring-2 focus:ring-indigo-500
            placeholder-gray-500
            transition-all duration-200
            overflow-hidden
          "
        />

        {/* Search Icon */}
        <svg
          className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-500 pointer-events-none"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          viewBox="0 0 24 24"
        >
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>

        {/* Filters Button */}
        <button
          onClick={() => setShowAdvanced(!showAdvanced)}
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
          <div className="absolute left-0 mt-3 w-full z-50">
            <SearchAdvanced className="w-full" />

          </div>
        )}

      </div>
    </div>
  );
}

