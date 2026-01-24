import { Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import HomePage from './pages/HomePage';
import BoardGameDetailPage from './pages/BoardGameDetailPage';
import AddGamePage from './pages/AddGamePage'; 
import NotFoundPage from './pages/NotFoundPage';

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<HomePage />} />
        <Route path="boardgame/:id" element={<BoardGameDetailPage />} />
        <Route path="add" element={<AddGamePage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Route>
    </Routes>
  );
}

export default App;