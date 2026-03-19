import { useState } from "react";

interface Props {
  className?: string;
  items?: React.ReactNode[];
  label: string;
  url?: string;
}

export function Dropdown({ className, items, label, url }: Props) {
  const [showItems, setShowItems] = useState(false);

  return (
    <div
      className={`relative inline-block ${className} flex flex-col items-center justify-center`}
      onMouseEnter={() => setShowItems(true)}
      onMouseLeave={() => setShowItems(false)}
    >
      <a href={url} className="cursor-pointer p-5">{label}</a>
      {showItems && items && items.length > 0 && (
        <div className="absolute left-0 top-1/2 mt-5 w-32 bg-[#33312d] border-[#7f7c7b] border-1 rounded-lg z-50 shadow-md p-1"
          onClick={() => setShowItems(false)}
        >
          {items.map((item, index) => (
            <DropdownItem key={index}>{item}</DropdownItem>
          ))}
        </div>
      )}
    </div>
  );
}



export function DropdownItem({ children }: { children: React.ReactNode }) {
  return (
    <div className="p-2 hover:bg-[#474540] rounded-lg cursor-pointer text-white">
      {children}
    </div>
  );
}