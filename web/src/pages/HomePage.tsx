import { useState, useEffect } from 'react';
import BoardGameCard from '../components/BoardGameCard';

interface BoardGame {
  id: number;
  name: string;
  min_players: number;
  max_players: number;
  coverImageUrl?: string;
}

export default function HomePage() {
  const [games, setGames] = useState<BoardGame[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('/api/boardgames')
      .then(res => res.json())
      .then(data => {
        setGames(data || []);
        setLoading(false);
      })
      .catch(err => {
        console.error(err);
        setLoading(false);
      });
  }, []);

  if (loading) {
    return <div className="text-white text-lg md:text-xl">Loading games...</div>;
  }

  return (
    <div className="h-screen overflow-y-auto">
      <h1 className="mb-4 md:mb-6 lg:mb-8 text-white text-2xl md:text-3xl lg:text-4xl font-bold pl-14 md:pl-0">
        Board Games
      </h1>

      {games.length === 0 ? (
        <div className="text-center py-10 md:py-15 px-4 text-[#999]">
          <div className="text-6xl md:text-7xl lg:text-8xl mb-4 md:mb-5">🎲</div>
          <h2 className="text-[#ccc] mb-2 text-lg md:text-xl lg:text-2xl">No games yet</h2>
          <p className="text-sm md:text-base">Start building your collection by adding your first board game!</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 md:gap-5 lg:gap-6">
          {games.map(game => (
            <BoardGameCard
              key={game.id}
              id={game.id}
              name={game.name}
              minPlayers={game.min_players}
              maxPlayers={game.max_players}
              coverImageUrl={game.coverImageUrl}
            />
          ))}
        </div>
      )}
    </div>
  );
}
