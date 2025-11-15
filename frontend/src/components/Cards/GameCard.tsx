
interface GameCardProps {
  title: string;
  description: string;
  icon: React.ElementType;
  bgColor: string;
  pattern?: string;
  playLink: string;
}

export function GameCard({ title, description, icon: Icon, bgColor, pattern, playLink } : GameCardProps) {
  return (
    <div className={`relative p-8 rounded-3xl shadow-2xl overflow-hidden group h-80 flex flex-col justify-between ${bgColor}`}>
        {/* Background Pattern */}
        {pattern && (
            <div className="absolute inset-0 opacity-10 group-hover:opacity-20 transition-opacity duration-300" style={{
                backgroundImage: pattern,
                backgroundSize: '30px 30px',
                maskImage: 'linear-gradient(to bottom, black 50%, transparent 100%)',
            }}></div>
        )}

        <div className="relative z-10">
            <Icon className="h-12 w-12 text-white mb-4 transition-transform duration-300 group-hover:-translate-y-1" />
            <h3 className="text-4xl font-extrabold text-white mb-2">{title}</h3>
            <p className="text-lg text-gray-200">{description}</p>
        </div>

        <div className="relative z-10 flex space-x-4 mt-6">
            <a href={playLink} className="px-6 py-2 bg-white text-gray-900 font-bold rounded-full shadow-lg hover:bg-gray-200 transition-transform duration-300 transform group-hover:scale-105">
                Play Now
            </a>
        </div>
    </div>
  )
}