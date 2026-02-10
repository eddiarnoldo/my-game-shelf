import { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import ImageGallery from '../components/ImageGallery';

interface BoardGameImageDto {
  id: number;
  imageUrl: string;
  thumbnailUrl: string;
  imageType: string;
  displayOrder: number;
}

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
  boardGameImages: BoardGameImageDto[];
}

export default function BoardGameDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [game, setGame] = useState<BoardGame | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [imageLoaded, setImageLoaded] = useState(false);
  const [gameplayImages, setGameplayImages] = useState<BoardGameImageDto[]>([]);

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
        const gameplay = data.boardGameImages
          ?.filter((img: BoardGameImageDto) => img.imageType === 'gameplay')
          .sort((a: BoardGameImageDto, b: BoardGameImageDto) => a.displayOrder - b.displayOrder) || [];
        setGameplayImages(gameplay);
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

      navigate('/');
    } catch (err) {
      console.error(err);
      alert('Failed to delete game. Please try again.');
      setDeleting(false);
    }
  };

  if (loading) {
    return <div className="text-white">Loading game...</div>;
  }

  if (error || !game) {
    return (
      <div>
        <h1>Game not found</h1>
        <Link to="/" className="text-[#4a9eff] no-underline">← Back to games</Link>
      </div>
    );
  }

  const coverImage = game.boardGameImages?.find(img => img.imageType === 'cover');

  return (
    <div className="h-screen overflow-y-auto">
      <div className="flex justify-between items-center mb-5">
        <Link to="/" className="text-[#4a9eff] no-underline">
          ← Back to games
        </Link>

        <button
          onClick={handleDelete}
          disabled={deleting}
          className={`border-none text-white text-2xl font-bold cursor-pointer px-3 py-2 rounded-lg transition-all duration-200 flex items-center justify-center w-10 h-10 ${
            deleting 
              ? 'bg-[#444] cursor-not-allowed' 
              : 'bg-[#ff4444] hover:scale-110'
          }`}
          title="Delete game"
        >
          ✕
        </button>
      </div>

      <div className="flex gap-10 mt-5">
        <div className="w-[400px] h-[400px] flex-shrink-0 relative">
          {coverImage ? (
            <>
              {!imageLoaded && (
                <div className="absolute inset-0 rounded-lg animate-pulse bg-[#333]" />
              )}
              <img
                src={coverImage.imageUrl}
                alt={`${game.name} cover`}
                onLoad={() => setImageLoaded(true)}
                className="w-full h-full object-cover rounded-lg"
                style={{
                  visibility: imageLoaded ? 'visible' : 'hidden',
                }}
              />
            </>
          ) : (
            <div className="w-full h-full bg-[#444] rounded-lg flex items-center justify-center text-[120px]">
              🎲
            </div>
          )}
        </div>

        <div className="flex-1">
          <h1 className="text-white mb-5 text-4xl">
            {game.name}
          </h1>
          
          <div className="text-[#ccc] text-base leading-[1.8]">
            <p className="mb-3">
              <strong className="text-white">Players:</strong> {game.min_players}-{game.max_players || game.min_players}
            </p>
            
            <p className="mb-3">
              <strong className="text-white">Play Time:</strong> {game.play_time} minutes
            </p>
            
            <p className="mb-3">
              <strong className="text-white">Minimum Age:</strong> {game.min_age}+
            </p>
            
            <p className="mb-5">
              <strong className="text-white">Description:</strong>
            </p>
            <p className="text-[#aaa] leading-[1.6] mb-5">
              {game.description}
            </p>

            <p className="text-sm text-[#666]">
              <strong>Added:</strong> {new Date(game.created_at).toLocaleDateString()}
            </p>
            <p className="text-sm text-[#666]">
              <strong>Last Updated:</strong> {new Date(game.updated_at).toLocaleDateString()}
            </p>
          </div>
        </div>
      </div>

      {game && (
        <div className="mt-10">
          <h2 className="text-white mb-5 text-2xl font-bold">
            Gameplay Images
          </h2>
          <ImageGallery
            boardGameId={game.id}
            images={gameplayImages}
            onImageUploaded={(newImage) => {
              setGameplayImages(prev => [...prev, newImage]);
            }}
          />
        </div>
      )}
    </div>
  );
}
