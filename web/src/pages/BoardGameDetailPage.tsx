import { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import ImageGallery from '../components/ImageGallery';
import { useAuth } from '../context/AuthContext';

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
  const { user, accessToken } = useAuth();
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
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${accessToken}` },
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
    return <div className="text-white text-lg md:text-xl">Loading game...</div>;
  }

  if (error || !game) {
    return (
      <div className="p-4">
        <h1 className="text-xl md:text-2xl mb-4">Game not found</h1>
        <Link to="/" className="text-[#4a9eff] no-underline text-base md:text-lg">← Back to games</Link>
      </div>
    );
  }

  const coverImage = game.boardGameImages?.find(img => img.imageType === 'cover');

  return (
    <div className="pb-20 md:pb-0">
      <div className="flex justify-between items-center mb-4 md:mb-5">
        <Link to="/" className="text-[#4a9eff] no-underline text-base md:text-lg pl-14 md:pl-0">
          ← Back to games
        </Link>

        {user && (
          <div className="flex items-center gap-2">
            <Link
              to={`/boardgame/${id}/edit`}
              className="text-white no-underline text-sm md:text-base font-semibold px-4 py-2 rounded-lg bg-[#4a9eff] hover:bg-[#3a8eef] transition-colors duration-200 min-h-[44px] flex items-center"
            >
              Edit
            </Link>
            <button
              onClick={handleDelete}
              disabled={deleting}
              className={`border-none text-white text-xl md:text-2xl font-bold cursor-pointer px-3 py-2 rounded-lg transition-all duration-200 flex items-center justify-center min-w-[44px] min-h-[44px] ${
                deleting
                  ? 'bg-[#444] cursor-not-allowed'
                  : 'bg-[#ff4444] hover:scale-110 active:scale-95'
              }`}
              title="Delete game"
            >
              ✕
            </button>
          </div>
        )}
      </div>

      {/* Main content - stack on mobile, side-by-side on desktop */}
      <div className="flex flex-col md:flex-row gap-6 md:gap-8 lg:gap-10 mt-4 md:mt-5">
        {/* Cover image */}
        <div className="w-full md:w-[350px] lg:w-[400px] md:flex-shrink-0 relative">
          <div className="w-full aspect-square relative">
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
              <div className="w-full h-full bg-[#444] rounded-lg flex items-center justify-center text-7xl md:text-8xl lg:text-[120px]">
                🎲
              </div>
            )}
          </div>
        </div>

        {/* Game details */}
        <div className="flex-1">
          <h1 className="text-white mb-4 md:mb-5 text-2xl md:text-3xl lg:text-4xl font-bold">
            {game.name}
          </h1>

          <div className="text-[#ccc] text-sm md:text-base leading-relaxed">
            <p className="mb-2 md:mb-3">
              <strong className="text-white">Players:</strong> {game.min_players}-{game.max_players || game.min_players}
            </p>

            <p className="mb-2 md:mb-3">
              <strong className="text-white">Play Time:</strong> {game.play_time} minutes
            </p>

            <p className="mb-2 md:mb-3">
              <strong className="text-white">Minimum Age:</strong> {game.min_age}+
            </p>

            <p className="mb-3 md:mb-4">
              <strong className="text-white">Description:</strong>
            </p>
            <p className="text-[#aaa] leading-relaxed mb-4 md:mb-5 text-sm md:text-base">
              {game.description}
            </p>

            <p className="text-xs md:text-sm text-[#666]">
              <strong>Added:</strong> {new Date(game.created_at).toLocaleDateString()}
            </p>
            <p className="text-xs md:text-sm text-[#666]">
              <strong>Last Updated:</strong> {new Date(game.updated_at).toLocaleDateString()}
            </p>
          </div>
        </div>
      </div>

      {game && (
        <div className="mt-8 md:mt-10">
          <h2 className="text-white mb-4 md:mb-5 text-xl md:text-2xl font-bold">
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
