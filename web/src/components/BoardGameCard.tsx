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
    <Link
      to={`/boardgame/${id}`}
      style={{ textDecoration: 'none' }}
    >
      <div
        style={{
          backgroundColor: '#2d2d2d',
          borderRadius: '8px',
          overflow: 'hidden',
          cursor: 'pointer',
          transition: 'transform 0.2s',
          maxWidth: '250px',  // Smaller
          width: '100%',
        }}
        onMouseEnter={(e) => e.currentTarget.style.transform = 'scale(1.02)'}
        onMouseLeave={(e) => e.currentTarget.style.transform = 'scale(1)'}
      >
        {/* Cover image or placeholder */}
        <div style={{
          width: '100%',
          aspectRatio: '1',  // Square!
          backgroundColor: '#444',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: '64px',
          overflow: 'hidden',
          position: 'relative'
        }}>
          {coverImageUrl ? (
            <img 
              src={coverImageUrl} 
              alt={`${name} cover`}
              style={{
                width: '100%',
                height: '100%',
                objectFit: 'cover',
                objectPosition: 'center'
              }}
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
        
        {/* Content section */}
        <div style={{ padding: '12px' }}>
          <h3 style={{ 
            color: 'white', 
            marginBottom: '6px', 
            fontSize: '16px',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap'
          }}>
            {name}
          </h3>
          <p style={{ color: '#999', fontSize: '13px' }}>
            {minPlayers}-{maxPlayers} players
          </p>
        </div>
      </div>
    </Link>
  );
}