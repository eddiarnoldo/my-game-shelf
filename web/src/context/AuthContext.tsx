import { createContext, useContext, useEffect, useRef, useState } from 'react';

interface AuthUser {
  id: number;
  username: string;
  role: string;
}

interface AuthContextValue {
  user: AuthUser | null;
  accessToken: string | null;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

function decodeJwtPayload(token: string): { user_id: number; username: string; role: string; exp: number } {
  const payload = token.split('.')[1];
  return JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')));
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const refreshTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const scheduleRefresh = (token: string) => {
    if (refreshTimerRef.current) clearTimeout(refreshTimerRef.current);
    try {
      const { exp } = decodeJwtPayload(token);
      const msUntilExpiry = exp * 1000 - Date.now();
      const delay = Math.max(msUntilExpiry - 60_000, 0);
      refreshTimerRef.current = setTimeout(() => doRefresh(), delay);
    } catch {
      // malformed token — ignore
    }
  };

  const doRefresh = async (): Promise<boolean> => {
    const storedRefreshToken = sessionStorage.getItem('refresh_token');
    if (!storedRefreshToken) return false;

    try {
      const res = await fetch('/api/refresh', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: storedRefreshToken }),
      });

      if (!res.ok) {
        clearAuth();
        return false;
      }

      const data = await res.json();
      applyTokens(data.access_token, data.refresh_token);
      return true;
    } catch {
      clearAuth();
      return false;
    }
  };

  const applyTokens = (newAccessToken: string, newRefreshToken: string) => {
    sessionStorage.setItem('access_token', newAccessToken);
    sessionStorage.setItem('refresh_token', newRefreshToken);
    setAccessToken(newAccessToken);

    try {
      const claims = decodeJwtPayload(newAccessToken);
      setUser({ id: claims.user_id, username: claims.username, role: claims.role });
    } catch {
      clearAuth();
      return;
    }

    scheduleRefresh(newAccessToken);
  };

  const clearAuth = () => {
    sessionStorage.removeItem('access_token');
    sessionStorage.removeItem('refresh_token');
    setAccessToken(null);
    setUser(null);
    if (refreshTimerRef.current) clearTimeout(refreshTimerRef.current);
  };

  // Restore session on mount
  useEffect(() => {
    const storedToken = sessionStorage.getItem('access_token');
    if (!storedToken) return;

    try {
      const claims = decodeJwtPayload(storedToken);
      if (claims.exp * 1000 > Date.now()) {
        setAccessToken(storedToken);
        setUser({ id: claims.user_id, username: claims.username, role: claims.role });
        scheduleRefresh(storedToken);
      } else {
        // Expired — try refresh
        doRefresh();
      }
    } catch {
      clearAuth();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const login = async (username: string, password: string) => {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });

    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      throw new Error(data.error ?? 'Login failed');
    }

    const data = await res.json();
    applyTokens(data.access_token, data.refresh_token);
  };

  const logout = async () => {
    const storedRefreshToken = sessionStorage.getItem('refresh_token');
    if (storedRefreshToken) {
      await fetch('/api/logout', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: storedRefreshToken }),
      }).catch(() => {});
    }
    clearAuth();
  };

  return (
    <AuthContext.Provider value={{ user, accessToken, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
