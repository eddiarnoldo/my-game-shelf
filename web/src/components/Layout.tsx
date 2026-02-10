import { Outlet, NavLink } from 'react-router-dom';

export default function Layout() {
  return (
    <div className="flex h-screen w-screen">
      <aside className="w-[250px] bg-[#2d2d2d] text-white p-5 flex-shrink-0 h-screen overflow-y-auto">
        <h2 className="mb-[30px]">🎲 My Game Shelf</h2>
        
        <nav>
          <ul className="list-none p-0">
            <li className="mb-2">
              <NavLink
                to="/"
                className={({ isActive }) => 
                  `text-white no-underline block px-3 py-2.5 rounded-md transition-colors duration-200 hover:bg-[#3d3d3d] ${isActive ? 'bg-[#444]' : 'bg-transparent'}`
                }
              >
                📚 Board Games
              </NavLink>
            </li>
            <li className="mb-2">
              <NavLink
                to="/add"
                className={({ isActive }) => 
                  `text-white no-underline block px-3 py-2.5 rounded-md transition-colors duration-200 hover:bg-[#3d3d3d] ${isActive ? 'bg-[#444]' : 'bg-transparent'}`
                }
              >
                ➕ Add Game
              </NavLink>
            </li>
          </ul>
        </nav>
      </aside>

      <main className="flex-1 h-screen p-10 overflow-y-auto bg-[#1a1a1a]">
        <Outlet />
      </main>
    </div>
  );
}
