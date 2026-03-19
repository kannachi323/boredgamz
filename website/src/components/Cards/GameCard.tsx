
interface GameCardProps {
  title: string;
  description: string;
  bgImg?: string;
  bgSize?: string;
  playLink: string;
}

export function GameCard({ title, description, bgImg, bgSize, playLink }: GameCardProps) {
  return (
    <div className="group relative flex flex-col items-center space-y-6 bg-gray-900 rounded-2xl">
      <div className="w-48 aspect-square overflow-hidden">
        <div
          className="w-full h-full"
          style={{
            backgroundImage: bgImg,
            backgroundSize: bgSize,
            backgroundPosition: "center",
            backgroundRepeat: 'no-repeat'
          }}
        ></div>
        <a
          href={playLink}
          className="
            absolute inset-0 
            flex items-center justify-center 
            opacity-0 group-hover:opacity-100
            transition-opacity duration-300 rounded-2xl
            bg-black/25
          "
        >
          <svg
            className="w-16 h-16 text-white"
            fill="currentColor"
            viewBox="0 0 24 24"
          >
            <path d="M8 5v14l11-7z" />
          </svg>
        </a>
      </div>

      <div className="w-full p-6 rounded-3xl flex flex-col items-start">
        <h3 className="text-2xl font-extrabold text-white mb-2">{title}</h3>
        <p className="text-gray-300 text-md mb-4">{description}</p>
      </div>
    </div>
  );
}
