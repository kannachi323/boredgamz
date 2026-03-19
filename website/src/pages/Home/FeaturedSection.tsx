import { useEffect, useState } from 'react'
import TICTACTOE from '../../assets/tictactoe.png'
import CONNECTFOUR from '../../assets/connect4.png'
import GOMOKU from '../../assets/small-board.jpg'
import MAHJONG from '../../assets/mahjong.png'
import POKER from '../../assets/poker.png'
import CHESS from '../../assets/chess.png'
import { GameCard } from '../../components/Cards/GameCard'

type PopularGame = {
  title: string
  description: string
  bgImg: string
  bgSize: string
  playLink: string
}

type PopularGameApiItem = {
  slug?: string
  key?: string
  game?: string
  title?: string
  name?: string
  description?: string
  playLink?: string
  path?: string
}

const defaultPopularGames: PopularGame[] = [
  {
    title: 'Tic Tac Toe',
    description: 'A quick classic for warm-ups and rematches.',
    bgImg: `url(${TICTACTOE})`,
    bgSize: '120% auto',
    playLink: '/games/tictactoe',
  },
  {
    title: 'Connect Four',
    description: 'Fast tactical gameplay with clean head-to-head rounds.',
    bgImg: `url(${CONNECTFOUR})`,
    bgSize: '100% auto',
    playLink: '/games/connectfour',
  },
  {
    title: 'Gomoku',
    description: 'Deep strategy and ranked matchmaking in one board.',
    bgImg: `url(${GOMOKU})`,
    bgSize: '80% auto',
    playLink: '/games/gomoku',
  },
  {
    title: 'Mahjong',
    description: 'Pattern memory and precision under pressure.',
    bgImg: `url(${MAHJONG})`,
    bgSize: '100% auto',
    playLink: '/games/mahjong',
  },
  {
    title: 'Poker (Texas Hold \'Em)',
    description: 'Read the table, manage risk, and outplay opponents.',
    bgImg: `url(${POKER})`,
    bgSize: '100% auto',
    playLink: '/games/poker',
  },
  {
    title: 'Chess',
    description: 'Openings, middle game pressure, and sharp endgames.',
    bgImg: `url(${CHESS})`,
    bgSize: '80% auto',
    playLink: '/games/chess',
  },
]

const imageBySlug: Record<string, { bgImg: string; bgSize: string; playLink: string }> = {
  tictactoe: { bgImg: `url(${TICTACTOE})`, bgSize: '120% auto', playLink: '/games/tictactoe' },
  connectfour: { bgImg: `url(${CONNECTFOUR})`, bgSize: '100% auto', playLink: '/games/connectfour' },
  gomoku: { bgImg: `url(${GOMOKU})`, bgSize: '80% auto', playLink: '/games/gomoku' },
  mahjong: { bgImg: `url(${MAHJONG})`, bgSize: '100% auto', playLink: '/games/mahjong' },
  poker: { bgImg: `url(${POKER})`, bgSize: '100% auto', playLink: '/games/poker' },
  chess: { bgImg: `url(${CHESS})`, bgSize: '80% auto', playLink: '/games/chess' },
}


export function FeaturedSection() {
  const [popularGames, setPopularGames] = useState<PopularGame[]>(defaultPopularGames)

  useEffect(() => {
    let isMounted = true

    async function loadPopularGames() {
      try {
        const res = await fetch(`${import.meta.env.VITE_SERVER_ROOT}/games/popular`, {
          credentials: 'include',
        })

        if (!res.ok) return

        const data = await res.json()
        if (!Array.isArray(data) || data.length === 0) return

        const mapped = (data as PopularGameApiItem[])
          .map((item) => {
            const slug = String(item.slug || item.key || item.game || '').toLowerCase()
            const fallback = imageBySlug[slug]

            if (!fallback) return null

            return {
              title: item.title || item.name || slug,
              description: item.description || 'Popular right now',
              bgImg: fallback.bgImg,
              bgSize: fallback.bgSize,
              playLink: item.playLink || item.path || fallback.playLink,
            } as PopularGame
          })
          .filter(Boolean) as PopularGame[]

        if (isMounted && mapped.length > 0) {
          setPopularGames(mapped)
        }
      } catch {
        // Keep defaults when endpoint is unavailable.
      }
    }

    loadPopularGames()

    return () => {
      isMounted = false
    }
  }, [])

  return (

  
    <section id="games" className="py-10 bg-zinc-950">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <h2 className="text-4xl font-extrabold text-white text-center mb-4">
          Most Popular Games
        </h2>
        <p className="text-xl text-gray-300 text-center mb-16">
          Live defaults today, ready for real-time ranking from your server.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {popularGames.map((game) => (
            <GameCard
              key={game.title}
              title={game.title}
              description={game.description}
              bgImg={game.bgImg}
              bgSize={game.bgSize}
              playLink={game.playLink}
            />
          ))}
        </div>
      </div>
    </section>
  )
}