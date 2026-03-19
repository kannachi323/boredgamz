import { useEffect } from "react";

export function useClickOutside(
  ref: React.RefObject<HTMLElement>,
  onOutsideClick: () => void
) {
  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        onOutsideClick();
      }
    }

    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, [ref, onOutsideClick]);
}
