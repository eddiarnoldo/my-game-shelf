import { useState, useEffect } from 'react';

interface BoardGame {
  id: number;
  name: string;
  min_players: number;
  max_players: number;
}

function App() {
  const [games, setGames] = useState<BoardGame[]>([]);

  useEffect(() => {
    fetch('http://localhost:8080/api/boardgames')
      .then(res => res.json())
      .then(data => setGames(data || []))
      .catch(err => console.error(err));
  }, []);

  return (
    <div style={{ padding: '20px' }}>
      <h1>🎲 My Game Shelf</h1>
      <ul>
        {games.map(game => (
          <li key={game.id}>
            {game.name} ({game.min_players}-{game.max_players} players)
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;