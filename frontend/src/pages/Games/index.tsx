
import { Square, Grid, Rows,  } from "lucide-react"
import { GameCard } from "../../components/Cards/GameCard"

export default function Games() {
  return (
    <section id="games" className="py-20 bg-gray-800">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <h2 className="text-4xl font-extrabold text-white text-center mb-4">
            Our Games
        </h2>
        <p className="text-xl text-gray-400 text-center mb-16">
            Jump into classic strategy games, re-imagined for seamless multiplayer fun.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <GameCard
              title="Tic-Tac-Toe"
              description="The timeless classic. Perfect for quick matches and challenging friends."
              icon={Square}
              bgColor="bg-gradient-to-br from-indigo-700 to-purple-800"
              pattern="linear-gradient(45deg, rgba(255,255,255,0.1) 25%, transparent 25%, transparent 75%, rgba(255,255,255,0.1) 75%, rgba(255,255,255,0.1)), linear-gradient(45deg, rgba(255,255,255,0.1) 25%, transparent 25%, transparent 75%, rgba(255,255,255,0.1) 75%, rgba(255,255,255,0.1))"
              playLink="/games/tictactoe"
          />
          <GameCard
              title="Connect Four"
              description="Strategize your way to four in a row. A game of calculated moves."
              icon={Grid}
              bgColor="bg-gradient-to-br from-green-700 to-teal-800"
              pattern="linear-gradient(90deg, #ffffff1a 1px, transparent 1px), linear-gradient(#ffffff1a 1px, transparent 1px)"
              playLink="/games/connectfour"
          />
          <GameCard
              title="Gomoku"
              description="Also known as Five-in-a-Row. A deep strategic challenge from Japan."
              icon={Rows}
              bgColor="bg-gradient-to-br from-red-700 to-orange-800"
              pattern="radial-gradient(circle, rgba(255,255,255,0.1) 1px, transparent 1px)"
              playLink="/games/gomoku"
          />
        </div>
      </div>
    </section>
  )
}