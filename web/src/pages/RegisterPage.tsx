import { useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';

export default function RegisterPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const [inviteCode, setInviteCode] = useState(searchParams.get('code') ?? '');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [passwordConfirmation, setPasswordConfirmation] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (password !== passwordConfirmation) {
      setError('Passwords do not match');
      return;
    }

    setLoading(true);

    try {
      const res = await fetch('/api/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          invite_code: inviteCode,
          username,
          password,
          password_confirmation: passwordConfirmation,
        }),
      });

      const data = await res.json().catch(() => ({}));

      if (!res.ok) {
        throw new Error(data.error ?? 'Registration failed');
      }

      navigate('/login', { state: { message: 'Account created! You can now sign in.' } });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  const inputClass = "w-full p-3 rounded-lg border border-[#444] bg-[#1a1a1a] text-white text-sm focus:outline-none focus:ring-2 focus:ring-[#4a9eff] focus:border-[#4a9eff] transition-colors duration-200 placeholder:text-[#666]";

  return (
    <div className="flex items-center justify-center min-h-screen px-4">
      <div className="w-full max-w-sm bg-[#2d2d2d] rounded-xl p-6 shadow-xl">
        <h1 className="text-white text-2xl font-bold mb-2 text-center">Create Account</h1>
        <p className="text-[#999] text-sm text-center mb-6">My Game Shelf</p>

        {error && (
          <div className="bg-[#ff4444]/20 border border-[#ff4444] text-[#ff4444] p-3 rounded-lg mb-4 text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-[#ccc] mb-1 text-sm font-medium">Invite Code</label>
            <input
              type="text"
              value={inviteCode}
              onChange={e => setInviteCode(e.target.value)}
              required
              placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
              className={inputClass}
            />
          </div>

          <div>
            <label className="block text-[#ccc] mb-1 text-sm font-medium">Username</label>
            <input
              type="text"
              value={username}
              onChange={e => setUsername(e.target.value)}
              required
              autoComplete="username"
              placeholder="choose a username"
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
              autoComplete="new-password"
              placeholder="••••••••"
              className={inputClass}
            />
          </div>

          <div>
            <label className="block text-[#ccc] mb-1 text-sm font-medium">Confirm Password</label>
            <input
              type="password"
              value={passwordConfirmation}
              onChange={e => setPasswordConfirmation(e.target.value)}
              required
              autoComplete="new-password"
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
            {loading ? 'Creating account...' : 'Create Account'}
          </button>
        </form>

        <p className="text-[#666] text-sm text-center mt-6">
          Already have an account?{' '}
          <Link to="/login" className="text-[#4a9eff] no-underline hover:underline">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  );
}
