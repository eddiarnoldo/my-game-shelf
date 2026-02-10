import { Outlet, NavLink } from 'react-router-dom';
import { useState } from 'react';

export default function Layout() {
  const [sidebarOpen, setSidebarOpen] = useState(false);

  return (
    <div className="flex h-screen w-screen">
      {/* Mobile hamburger button */}
      <button
        onClick={() => setSidebarOpen(true)}
        className="md:hidden fixed top-4 left-4 z-50 bg-[#2d2d2d] text-white p-3 rounded-lg shadow-lg"
        aria-label="Open menu"
      >
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>

      {/* Mobile sidebar overlay */}
      {sidebarOpen && (
        <div
          className="md:hidden fixed inset-0 bg-black/50 z-40 transition-opacity duration-300"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar - slide out on mobile, always visible on desktop */}
      <aside
        className={`fixed md:static inset-y-0 left-0 z-50 w-[250px] bg-[#2d2d2d] text-white p-5 flex-shrink-0 h-screen overflow-y-auto transform transition-transform duration-300 ease-in-out ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'
        }`}
      >
        <div className="flex justify-between items-center mb-8">
          <h2 className="text-lg md:text-xl font-semibold">🎲 My Game Shelf</h2>
          <button
            onClick={() => setSidebarOpen(false)}
            className="md:hidden text-white p-2"
            aria-label="Close menu"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <nav>
          <ul className="list-none p-0 space-y-2">
            <li>
              <NavLink
                to="/"
                onClick={() => setSidebarOpen(false)}
                className={({ isActive }) =>
                  `text-white no-underline block px-3 py-3 rounded-md transition-colors duration-200 hover:bg-[#3d3d3d] text-base md:text-lg ${isActive ? 'bg-[#444]' : 'bg-transparent'}`
                }
              >
                📚 Board Games
              </NavLink>
            </li>
            <li>
              <NavLink
                to="/add"
                onClick={() => setSidebarOpen(false)}
                className={({ isActive }) =>
                  `text-white no-underline block px-3 py-3 rounded-md transition-colors duration-200 hover:bg-[#3d3d3d] text-base md:text-lg ${isActive ? 'bg-[#444]' : 'bg-transparent'}`
                }
              >
                ➕ Add Game
              </NavLink>
            </li>
          </ul>
        </nav>
      </aside>

      {/* Main content */}
      <main className="flex-1 h-screen p-4 md:p-6 lg:p-10 overflow-y-auto bg-[#1a1a1a] w-full">
        <Outlet />
      </main>
    </div>
  );
}
