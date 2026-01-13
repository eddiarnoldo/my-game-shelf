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
    return <div style={{ color: 'white' }}>Loading games...</div>;
  }

  return (
    <div>
      <h1 style={{ marginBottom: '30px', color: 'white' }}>Board Games</h1>
      
      {games.length === 0 ? (
        // Empty state
        <div style={{
          textAlign: 'center',
          padding: '60px 20px',
          color: '#999'
        }}>
          <div style={{ fontSize: '80px', marginBottom: '20px' }}>🎲</div>
          <h2 style={{ color: '#ccc', marginBottom: '10px' }}>No games yet</h2>
          <p>Start building your collection by adding your first board game!</p>
        </div>
      ) : (
        // Game grid
      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fill, minmax(220px, 250px))',  // Smaller cards
        gap: '20px',
        justifyContent: 'start'
      }}>
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