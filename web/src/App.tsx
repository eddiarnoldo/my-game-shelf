import { Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import ProtectedRoute from './components/ProtectedRoute';
import Layout from './components/Layout';
import HomePage from './pages/HomePage';
import BoardGameDetailPage from './pages/BoardGameDetailPage';
import AddGamePage from './pages/AddGamePage';
import EditGamePage from './pages/EditGamePage';
import NotFoundPage from './pages/NotFoundPage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import InvitePage from './pages/InvitePage';

function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<HomePage />} />
          <Route path="boardgame/:id" element={<BoardGameDetailPage />} />
          <Route path="add" element={
            <ProtectedRoute><AddGamePage /></ProtectedRoute>
          } />
          <Route path="boardgame/:id/edit" element={
            <ProtectedRoute><EditGamePage /></ProtectedRoute>
          } />
          <Route path="login" element={<LoginPage />} />
          <Route path="join" element={<RegisterPage />} />
          <Route path="invite" element={
            <ProtectedRoute requireAdmin><InvitePage /></ProtectedRoute>
          } />
          <Route path="*" element={<NotFoundPage />} />
        </Route>
      </Routes>
    </AuthProvider>
  );
}

export default App;
