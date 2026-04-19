import { useState } from 'react';
import { useAuth } from '../context/AuthContext';

export default function InvitePage() {
  const { accessToken } = useAuth();
  const [email, setEmail] = useState('');
  const [status, setStatus] = useState<'idle' | 'loading' | 'success' | 'error'>('idle');
  const [message, setMessage] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setStatus('loading');
    setMessage('');

    try {
      const res = await fetch('/api/invites', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${accessToken}`,
        },
        body: JSON.stringify({ email }),
      });

      const data = await res.json().catch(() => ({}));

      if (!res.ok) {
        throw new Error(data.error ?? 'Failed to send invite');
      }

      setStatus('success');
      setMessage(`Invite sent to ${email}`);
      setEmail('');
    } catch (err) {
      setStatus('error');
      setMessage(err instanceof Error ? err.message : 'Failed to send invite');
    }
  };

  return (
    <div className="flex items-center justify-center min-h-[60vh] px-4">
      <div className="w-full max-w-sm bg-[#2d2d2d] rounded-xl p-6 shadow-xl">
        <h1 className="text-white text-2xl font-bold mb-2">Send Invite</h1>
        <p className="text-[#999] text-sm mb-6">
          Enter an email address to send an account invite.
        </p>

        {status === 'success' && (
          <div className="bg-green-900/30 border border-green-500 text-green-400 p-3 rounded-lg mb-4 text-sm">
            {message}
          </div>
        )}

        {status === 'error' && (
          <div className="bg-[#ff4444]/20 border border-[#ff4444] text-[#ff4444] p-3 rounded-lg mb-4 text-sm">
            {message}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-[#ccc] mb-1 text-sm font-medium">Email Address</label>
            <input
              type="email"
              value={email}
              onChange={e => setEmail(e.target.value)}
              required
              placeholder="user@example.com"
              className="w-full p-3 rounded-lg border border-[#444] bg-[#1a1a1a] text-white text-sm focus:outline-none focus:ring-2 focus:ring-[#4a9eff] focus:border-[#4a9eff] transition-colors duration-200 placeholder:text-[#666]"
            />
          </div>

          <button
            type="submit"
            disabled={status === 'loading'}
            className={`w-full py-3 rounded-lg font-semibold text-white border-none cursor-pointer transition-colors duration-200 ${
              status === 'loading'
                ? 'bg-[#444] text-[#999] cursor-not-allowed'
                : 'bg-[#4a9eff] hover:bg-[#3a8eef]'
            }`}
          >
            {status === 'loading' ? 'Sending...' : 'Send Invite'}
          </button>
        </form>
      </div>
    </div>
  );
}
