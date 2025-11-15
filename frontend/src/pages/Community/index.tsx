
import { MessageSquare, Camera, Zap } from "lucide-react";

const socialPlatforms: SocialPlatform[] = [
    { name: 'Discord', icon: MessageSquare, description: 'Join our official server for real-time LFG, announcements, and chat with developers.', color: 'bg-indigo-600/50 border-indigo-500', link: 'https://discord.gg/boredgamz' },
    { name: 'YouTube', icon: undefined, description: 'Watch game tutorials, review breakdowns, and tournament highlights.', color: 'bg-red-600/50 border-red-500', link: 'https://youtube.com/boredgamz' },
    { name: 'Instagram', icon: Camera, description: 'See gorgeous game photography, community fan art, and featured collections.', color: 'bg-pink-600/50 border-pink-500', link: 'https://instagram.com/boredgamz' },
    { name: 'Twitch', icon: Zap, description: 'Stream with us! Live gameplay and Q&A sessions every week.', color: 'bg-purple-600/50 border-purple-500', link: 'https://twitch.tv/boredgamz' },
];

interface SocialPlatform {
    name: string;
    icon?: React.ElementType;
    description: string;
    color: string;
    link: string;
}

export default function Community() {
  return (
    <section id="community" className="py-20 bg-gray-900 min-h-[80vh]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <h2 className="text-5xl font-extrabold text-white text-center mb-4">
            Join the Global Gamers Hub
        </h2>
        <p className="text-xl text-gray-400 text-center mb-16 max-w-3xl mx-auto">
            Connect with players, share your best games, and stay up-to-date with announcements across all our platforms.
        </p>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
            {socialPlatforms.map((platform) => (
                <a 
                    key={platform.name} 
                    href={platform.link} 
                    target="_blank" 
                    rel="noopener noreferrer"
                    className={`p-6 rounded-2xl shadow-2xl transition duration-300 flex flex-col items-center text-center hover:scale-[1.03] ${platform.color} border border-transparent hover:border-white/50`}
                    style={{ backdropFilter: 'blur(5px)' }} // Adds a slight frosted effect
                >
                    {platform.icon && <platform.icon className="w-10 h-10 mb-4" />}
                    <h3 className="text-3xl font-bold text-white mb-2">{platform.name}</h3>
                    <p className="text-gray-200 text-sm flex-grow">{platform.description}</p>
                    <span className="mt-4 text-sm font-semibold text-green-400 hover:text-green-300 transition">
                        Connect Now &rarr;
                    </span>
                </a>
            ))}
        </div>

        <div className="mt-20 p-8 bg-gray-800 rounded-2xl border border-indigo-700 text-center">
            <h3 className="text-2xl font-bold text-white mb-2">Can't find a group?</h3>
            <p className="text-gray-400 mb-4">Jump straight into our Discord to use the 'looking-for-game' channels and get rolling instantly!</p>
            <a href="https://discord.gg/boredgamz" target="_blank" rel="noopener noreferrer" className="inline-flex items-center px-6 py-3 border border-transparent text-base font-bold rounded-xl text-white bg-indigo-600 hover:bg-indigo-700 shadow-xl transition-transform duration-300 transform hover:-translate-y-1">
                <MessageSquare className="w-5 h-5 mr-2" />
                Join Discord
            </a>
        </div>
      </div>
    </section>
  );
}