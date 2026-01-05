import { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';

interface BoardGame {
  id: number;
  name: string;
  min_players: number;
  max_players: number;
  play_time: number;
  min_age: number;
  description: string;
  created_at: string;
  updated_at: string;
}

export default function GameDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [game, setGame] = useState<BoardGame | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    fetch(`/api/boardgames/${id}`)
      .then(res => {
        if (!res.ok) {
          throw new Error('Game not found');
        }
        return res.json();
      })
      .then(data => {
        setGame(data);
        setLoading(false);
      })
      .catch(err => {
        console.error(err);
        setError(true);
        setLoading(false);
      });
  }, [id]);

  const handleDelete = async () => {
    if (!window.confirm(`Are you sure you want to delete "${game?.name}"?`)) {
      return;
    }

    setDeleting(true);

    try {
      const response = await fetch(`/api/boardgames/${id}`, {
        method: 'DELETE'
      });

      if (!response.ok) {
        throw new Error('Failed to delete game');
      }

      // Success! Navigate back to home
      navigate('/');
    } catch (err) {
      console.error(err);
      alert('Failed to delete game. Please try again.');
      setDeleting(false);
    }
  };

  if (loading) {
    return <div style={{ color: 'white' }}>Loading game...</div>;
  }

  if (error || !game) {
    return (
      <div style={{ color: 'white' }}>
        <h1>Game not found</h1>
        <Link to="/" style={{ color: '#4a9eff' }}>← Back to games</Link>
      </div>
    );
  }

  return (
    <div>
      {/* Header with back button and delete */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <Link to="/" style={{ color: '#4a9eff', textDecoration: 'none' }}>
          ← Back to games
        </Link>

    <button
    onClick={handleDelete}
    disabled={deleting}
    style={{
        backgroundColor: deleting ? '#444' : '#ff4444',
        border: 'none',
        color: 'white',
        fontSize: '24px',
        fontWeight: 'bold',
        cursor: deleting ? 'not-allowed' : 'pointer',
        padding: '8px 12px',
        borderRadius: '8px',
        transition: 'transform 0.2s, background-color 0.2s',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        width: '40px',
        height: '40px'
    }}
    onMouseEnter={(e) => !deleting && (e.currentTarget.style.transform = 'scale(1.1)')}
    onMouseLeave={(e) => e.currentTarget.style.transform = 'scale(1)'}
    title="Delete game"
    >
    ✕
    </button>
      </div>

      <div style={{ display: 'flex', gap: '40px', marginTop: '20px' }}>
        {/* Game image */}
        <div style={{
          width: '400px',
          height: '400px',
          backgroundColor: '#444',
          borderRadius: '8px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: '120px',
          flexShrink: 0
        }}>
          🎲
        </div>

        {/* Game details */}
        <div style={{ flex: 1 }}>
          <h1 style={{ color: 'white', marginBottom: '20px', fontSize: '36px' }}>
            {game.name}
          </h1>
          
          <div style={{ color: '#ccc', fontSize: '16px', lineHeight: '1.8' }}>
            <p style={{ marginBottom: '12px' }}>
              <strong style={{ color: 'white' }}>Players:</strong> {game.min_players}-{game.max_players || game.min_players}
            </p>
            
            <p style={{ marginBottom: '12px' }}>
              <strong style={{ color: 'white' }}>Play Time:</strong> {game.play_time} minutes
            </p>
            
            <p style={{ marginBottom: '12px' }}>
              <strong style={{ color: 'white' }}>Minimum Age:</strong> {game.min_age}+
            </p>
            
            <p style={{ marginBottom: '20px' }}>
              <strong style={{ color: 'white' }}>Description:</strong>
            </p>
            <p style={{ color: '#aaa', lineHeight: '1.6', marginBottom: '20px' }}>
              {game.description}
            </p>

            <p style={{ fontSize: '14px', color: '#666' }}>
              <strong>Added:</strong> {new Date(game.created_at).toLocaleDateString()}
            </p>
            <p style={{ fontSize: '14px', color: '#666' }}>
              <strong>Last Updated:</strong> {new Date(game.updated_at).toLocaleDateString()}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}