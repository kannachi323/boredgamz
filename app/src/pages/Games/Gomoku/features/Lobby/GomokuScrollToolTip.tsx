import { useState } from "react";
import { X, ArrowDown, HelpCircle } from "lucide-react";

export function GomokuScrollTooltip() {
  const [visible, setVisible] = useState(true);

  if (!visible) return (
    <div
      className="w-[32px] aspect-square flex justify-center items-center cursor-pointer"
      onClick={() => setVisible(true)}
    >
      <HelpCircle className=" hover:text-[#817670]" size={32} />
    </div>
  );

  return (
    <>
      {/* Text + arrows */}
      <div className="flex items-center gap-2 animate-bounce">
        <span>Scroll down for more options</span>

        {/* Down arrows */}
        <div className="flex flex-col leading-none">
          <ArrowDown size={16} />
        </div>
      </div>

      {/* Close button */}
      <button
        onClick={() => setVisible(false)}
        className="text-[#C3B299] hover:text-[#817670] transition cursor-pointer"
      >
        <X size={18} />
      </button>
    </>
  );
}
