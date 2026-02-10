import { Link } from 'react-router-dom';

interface BoardGameCardProps {
  id: number;
  name: string;
  minPlayers: number;
  maxPlayers: number;
  coverImageUrl?: string;
}

export default function BoardGameCard({ id, name, minPlayers, maxPlayers, coverImageUrl }: BoardGameCardProps) {
  return (
    <Link to={`/boardgame/${id}`} className="no-underline block">
      <div className="bg-[#2d2d2d] rounded-lg overflow-hidden cursor-pointer transition-transform duration-200 w-full hover:scale-[1.02] active:scale-95">
        <div className="w-full aspect-square bg-[#444] flex items-center justify-center text-5xl md:text-6xl overflow-hidden relative">
          {coverImageUrl ? (
            <img
              src={coverImageUrl}
              alt={`${name} cover`}
              className="w-full h-full object-cover object-center"
              onError={(e) => {
                const target = e.currentTarget;
                target.style.display = 'none';
                if (target.parentElement) {
                  target.parentElement.innerHTML = '🎲';
                }
              }}
            />
          ) : (
            '🎲'
          )}
        </div>

        <div className="p-3 md:p-4">
          <h3 className="text-white mb-1 text-sm md:text-base font-medium overflow-hidden text-ellipsis whitespace-nowrap">
            {name}
          </h3>
          <p className="text-[#999] text-xs md:text-sm">
            {minPlayers}-{maxPlayers} players
          </p>
        </div>
      </div>
    </Link>
  );
}
