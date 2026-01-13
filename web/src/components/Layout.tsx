import { Outlet, NavLink } from 'react-router-dom';

export default function Layout() {
  return (
    <div style={{ display: 'flex', height: '100vh', width: '100vw' }}>
      {/* Sidebar */}
      <aside style={{ 
        width: '250px', 
        backgroundColor: '#2d2d2d', 
        color: 'white',
        padding: '20px',
        flexShrink: 0,
        height: '100vh',
        overflowY: 'auto'
      }}>
        <h2 style={{ marginBottom: '30px' }}>🎲 My Game Shelf</h2>
        
        <nav>
          <ul style={{ listStyle: 'none', padding: 0 }}>
            <li style={{ marginBottom: '8px' }}>
              <NavLink
                to="/"
                style={({ isActive }) => ({
                  color: 'white',
                  textDecoration: 'none',
                  display: 'block',
                  padding: '10px 12px',
                  borderRadius: '6px',
                  backgroundColor: isActive ? '#444' : 'transparent',
                  transition: 'background-color 0.2s'
                })}
              >
                📚 Board Games
              </NavLink>
            </li>
            <li style={{ marginBottom: '8px' }}>
                <NavLink
                to="/add"
                style={({ isActive }) => ({
                    color: 'white',
                    textDecoration: 'none',
                    display: 'block',
                    padding: '10px 12px',
                    borderRadius: '6px',
                    backgroundColor: isActive ? '#444' : 'transparent',
                    transition: 'background-color 0.2s'
                })}
                >
                ➕ Add Game
                </NavLink>
             </li>
          </ul>
        </nav>
      </aside>

      {/* Main content area */}
      <main style={{ 
        flex: 1,
        height: '100vh',
        padding: '40px',
        overflowY: 'auto',
        backgroundColor: '#1a1a1a'
      }}>
        <Outlet />
      </main>
    </div>
  );
}