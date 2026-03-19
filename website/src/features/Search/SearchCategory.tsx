export function SearchCategory({ title, options }: { title: string; options: string[] }) {
  return (
    <div className="mb-6">
      <h4 className="text-sm font-semibold text-gray-400 mb-2">{title}</h4>
      <div className="flex flex-wrap gap-2">
        {options.map((opt) => (
          <Tag key={opt}>{opt}</Tag>
        ))}
      </div>
    </div>
  );
}

function Tag({ children }: { children: string }) {
  return (
    <div
      className="
        px-3 py-1 
        bg-gray-800 
        text-gray-300 
        rounded-full 
        text-sm 
        cursor-pointer
        hover:bg-gray-700 
        hover:text-white
        transition
      "
    >
      {children}
    </div>
  );
}
