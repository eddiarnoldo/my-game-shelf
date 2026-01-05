import { Link } from 'react-router-dom';

   interface BoardGameCardProps {
     id: number;
     name: string;
     minPlayers: number;
     maxPlayers: number;
   }

   export default function BoardGameCard({ id, name, minPlayers, maxPlayers }: BoardGameCardProps) {
     return (
       <Link
         to={`/boardgame/${id}`}
         style={{ textDecoration: 'none' }}
       >
         <div
           style={{
             backgroundColor: '#2d2d2d',
             borderRadius: '8px',
             padding: '16px',
             cursor: 'pointer',
             transition: 'transform 0.2s',
           }}
           onMouseEnter={(e) => e.currentTarget.style.transform = 'scale(1.02)'}
           onMouseLeave={(e) => e.currentTarget.style.transform = 'scale(1)'}
         >
           {/* Placeholder image */}
           <div style={{
             width: '100%',
             height: '200px',
             backgroundColor: '#444',
             borderRadius: '4px',
             marginBottom: '12px',
             display: 'flex',
             alignItems: 'center',
             justifyContent: 'center',
             fontSize: '64px'
           }}>
             🎲
           </div>
           
           <h3 style={{ color: 'white', marginBottom: '8px', fontSize: '18px' }}>
             {name}
           </h3>
           <p style={{ color: '#999', fontSize: '14px' }}>
             {minPlayers}-{maxPlayers} players
           </p>
         </div>
       </Link>
     );
   }