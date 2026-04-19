import { useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const from = (location.state as { from?: { pathname: string } })?.from?.pathname ?? '/';

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await login(username, password);
      navigate(from, { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  const inputClass = "w-full p-3 rounded-lg border border-[#444] bg-[#1a1a1a] text-white text-sm focus:outline-none focus:ring-2 focus:ring-[#4a9eff] focus:border-[#4a9eff] transition-colors duration-200 placeholder:text-[#666]";

  return (
    <div className="flex items-center justify-center min-h-screen px-4">
      <div className="w-full max-w-sm bg-[#2d2d2d] rounded-xl p-6 shadow-xl">
        <h1 className="text-white text-2xl font-bold mb-2 text-center">Sign In</h1>
        <p className="text-[#999] text-sm text-center mb-6">My Game Shelf</p>

        {error && (
          <div className="bg-[#ff4444]/20 border border-[#ff4444] text-[#ff4444] p-3 rounded-lg mb-4 text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-[#ccc] mb-1 text-sm font-medium">Username</label>
            <input
              type="text"
              value={username}
              onChange={e => setUsername(e.target.value)}
              required
              autoComplete="username"
              placeholder="your username"
              className={inputClass}
            />
          </div>

          <div>
            <label className="block text-[#ccc] mb-1 text-sm font-medium">Password</label>
            <input
              type="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              required
              autoComplete="current-password"
              placeholder="••••••••"
              className={inputClass}
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className={`w-full py-3 rounded-lg font-semibold text-white border-none cursor-pointer transition-colors duration-200 ${
              loading
                ? 'bg-[#444] text-[#999] cursor-not-allowed'
                : 'bg-[#4a9eff] hover:bg-[#3a8eef]'
            }`}
          >
            {loading ? 'Signing in...' : 'Sign In'}
          </button>
        </form>

        <p className="text-[#666] text-sm text-center mt-6">
          Have an invite?{' '}
          <Link to="/join" className="text-[#4a9eff] no-underline hover:underline">
            Create account
          </Link>
        </p>
      </div>
    </div>
  );
}
