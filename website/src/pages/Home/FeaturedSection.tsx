import { Rocket, Swords, Users, BookOpen } from 'lucide-react';

interface Category {
  id: number;
  name: string;
  icon: React.ElementType; // Used to correctly type the Lucide icon component
  description: string;
  color: string;
}

const featuredCategories: Category[] = [
  { id: 1, name: 'Strategy Epics', icon: Swords, description: 'Deep, complex games for the masterful tactician.', color: 'border-indigo-500 hover:bg-indigo-900/50' },
  { id: 2, name: 'Party Games', icon: Users, description: 'Quick, fun, and chaotic—perfect for groups.', color: 'border-green-500 hover:bg-green-900/50' },
  { id: 3, name: 'Cooperative Quests', icon: Rocket, description: 'Work together or fail together. High stakes, high reward.', color: 'border-rose-500 hover:bg-rose-900/50' },
  { id: 4, name: 'Solo Adventures', icon: BookOpen, description: 'Explore worlds and challenge yourself, all on your own time.', color: 'border-yellow-500 hover:bg-yellow-900/50' },
];


export function FeaturedSection() {
  return (

  
    <section id="games" className="py-20 bg-gray-800">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <h2 className="text-4xl font-extrabold text-white text-center mb-4">
          Discover New Worlds
        </h2>
        <p className="text-xl text-gray-400 text-center mb-16">
          Browse games by the type of experience you're craving tonight.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          {featuredCategories.map((category) => (
            <div key={category.id} className={`p-6 rounded-2xl bg-gray-900 border-t-4 ${category.color} shadow-2xl transition duration-300`}>
              {/* The icon is rendered dynamically using the 'icon' field from the Category interface */}
              <category.icon className="h-10 w-10 text-green-400 mb-4" />
              <h3 className="text-2xl font-bold text-white mb-2">{category.name}</h3>
              <p className="text-gray-400">{category.description}</p>
              <button className="mt-4 text-sm font-semibold text-green-400 hover:text-green-300 transition">
                View Collection &rarr;
              </button>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}