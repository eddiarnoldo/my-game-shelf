import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useState } from 'react';
import { useAuth } from '../context/AuthContext';

export default function Layout() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/');
  };

  const navLinkClass = ({ isActive }: { isActive: boolean }) =>
    `text-white no-underline block px-3 py-3 rounded-md transition-colors duration-200 hover:bg-[#3d3d3d] text-base md:text-lg ${isActive ? 'bg-[#444]' : 'bg-transparent'}`;

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

      {/* Sidebar */}
      <aside
        className={`fixed md:static inset-y-0 left-0 z-50 w-[250px] bg-[#2d2d2d] text-white p-5 flex-shrink-0 h-screen overflow-y-auto flex flex-col transform transition-transform duration-300 ease-in-out ${
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

        <nav className="flex-1">
          <ul className="list-none p-0 space-y-2">
            <li>
              <NavLink to="/" onClick={() => setSidebarOpen(false)} className={navLinkClass}>
                📚 Board Games
              </NavLink>
            </li>

            {user?.role === 'admin' && (
              <li>
                <NavLink to="/add" onClick={() => setSidebarOpen(false)} className={navLinkClass}>
                  ➕ Add Game
                </NavLink>
              </li>
            )}

            {user?.role === 'admin' && (
              <li>
                <NavLink to="/invite" onClick={() => setSidebarOpen(false)} className={navLinkClass}>
                  ✉ Invite
                </NavLink>
              </li>
            )}
          </ul>
        </nav>

        {/* User section at bottom of sidebar */}
        <div className="mt-auto pt-4 border-t border-[#444]">
          {user ? (
            <div className="flex items-center gap-3">
              <div className="w-9 h-9 rounded-full bg-[#4a9eff] flex items-center justify-center text-white font-bold text-sm flex-shrink-0">
                {user.username[0].toUpperCase()}
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-white text-sm font-medium truncate">{user.username}</p>
                <p className="text-[#999] text-xs capitalize">{user.role}</p>
              </div>
              <button
                onClick={handleLogout}
                className="text-[#999] hover:text-white text-xs px-2 py-1 rounded hover:bg-[#3d3d3d] transition-colors duration-200 flex-shrink-0"
                aria-label="Log out"
              >
                Out
              </button>
            </div>
          ) : (
            <NavLink to="/login" onClick={() => setSidebarOpen(false)} className={navLinkClass}>
              Sign In
            </NavLink>
          )}
        </div>
      </aside>

      {/* Main content */}
      <main className="flex-1 h-screen p-4 md:p-6 lg:p-10 overflow-y-auto bg-[#1a1a1a] w-full">
        <Outlet />
      </main>
    </div>
  );
}
