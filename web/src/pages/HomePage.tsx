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
    return <div className="text-white">Loading games...</div>;
  }

  return (
    <div className="h-screen overflow-y-auto">
      <h1 className="mb-[30px] text-white">Board Games</h1>
      
      {games.length === 0 ? (
        <div className="text-center py-15 px-5 text-[#999]">
          <div className="text-[80px] mb-5">🎲</div>
          <h2 className="text-[#ccc] mb-2.5">No games yet</h2>
          <p>Start building your collection by adding your first board game!</p>
        </div>
      ) : (
        <div className="grid grid-cols-[repeat(auto-fill,minmax(220px,250px))] gap-5 justify-start">
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
