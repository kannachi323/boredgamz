
import TICTACTOE from "../../assets/tictactoe.png"
import CONNECTFOUR from "../../assets/connect4.png"
import GOMOKU from "../../assets/small-board.jpg"
import MAHJONG from "../../assets/mahjong.png"
import POKER from "../../assets/poker.png"
import CHESS from "../../assets/chess.png"

import { GameCard } from "../../components/Cards/GameCard"
import { SearchBar } from "../../features/Search/SearchBar"

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

       
        <SearchBar />


        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <GameCard
              title="Tic-Tac-Toe"
              description="The timeless classic. Perfect for quick matches and challenging friends."
              bgImg={`url(${TICTACTOE})`}
              bgSize='120% auto'
              playLink="/games/tictactoe"
          />
          <GameCard
              title="Connect Four"
              description="Strategize your way to four in a row. A game of calculated moves."
              bgImg={`url(${CONNECTFOUR})`}
              bgSize='100% auto'
              playLink="/games/connectfour"
          />
          <GameCard
            title="Gomoku"
            description="Also known as Five-in-a-Row. A deep strategic challenge from Japan."
            bgImg={`url(${GOMOKU})`}
            bgSize='80% auto'
            playLink="/games/gomoku"
          />
          <GameCard
            title="Mahjong"
            description="A timeless tile-matching classic that blends strategy, memory, and quick thinking."
            bgImg={`url(${MAHJONG})`}
            bgSize='100% auto'
            playLink="/games/mahjong"
          />
          <GameCard
            title="Poker (Texas Hold 'Em)"
            description="High stakes, sharp strategy, and endless mind games — welcome to Poker"
            bgImg={`url(${POKER})`}
            bgSize='100% auto'
            playLink="/games/poker"
          />
          <GameCard
            title="Chess"
            description="Master openings, control the board, and checkmate your opponent in this legendary game."
            bgImg={`url(${CHESS})`}
            bgSize='80% auto'
            playLink="/games/chess"
          />
        </div>
      </div>
    </section>
  )
}