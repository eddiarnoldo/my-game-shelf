# Responsive Design Refactoring Plan

## Overview
Refactor the My Game Shelf project from desktop-only to mobile-first responsive design using Tailwind CSS v4.

## User Requirements
1. **Sidebar**: Slide-out drawer on mobile
2. **Navigation**: Side drawer (not full screen)
3. **HomePage Grid**: 2 columns on tablet (md breakpoint)
4. **Lightbox**: Swipeable images on mobile
5. **Transitions**: Smooth animations for sidebar open/close

## Implementation Tasks

### Phase 1: Tailwind Configuration
- Update `index.css` with custom scrollbar styles, smooth scrolling, and mobile optimizations
- Keep default Tailwind v4 breakpoints: sm:640px, md:768px, lg:1024px, xl:1280px

### Phase 2: Core Layout
**File: `web/src/components/Layout.tsx`**
- Add mobile state management for sidebar visibility
- Implement hamburger menu button (visible only on mobile: `md:hidden`)
- Convert sidebar to fixed position slide-out drawer on mobile
- Add backdrop overlay when sidebar is open on mobile
- Smooth transition animations for sidebar
- Main content area: full width on mobile, adjusted for sidebar on desktop

### Phase 3: Page Components

**File: `web/src/pages/HomePage.tsx`**
- Responsive grid: 1 col mobile → 2 col md → 3 col lg
- Replace auto-fill grid with explicit responsive columns
- Responsive padding: `p-4 md:p-6 lg:p-10`
- Empty state with responsive emoji and text sizes

**File: `web/src/pages/BoardGameDetailPage.tsx`**
- **Critical**: Stack layout vertically on mobile (`flex-col`)
- Desktop: Side-by-side layout (`md:flex-row`)
- Cover image: full width mobile, fixed 400px desktop
- Responsive typography scale
- Delete button repositioned for mobile
- Gameplay gallery with responsive grid

**File: `web/src/pages/AddGamePage.tsx`**
- Form container: full width mobile, max-w-2xl desktop
- Larger input touch targets
- Full-width submit button on mobile
- Responsive image preview sizing

### Phase 4: Supporting Components

**File: `web/src/components/BoardGameCard.tsx`**
- Fluid width: `w-full` instead of fixed max-width
- Maintain aspect ratio across all sizes
- Remove hover scale on touch devices

**File: `web/src/components/ImageGallery.tsx`**
- Grid: 1 col mobile → 2 col sm → 3 col md
- Touch-friendly upload button (min 44px)

**File: `web/src/components/Lightbox.tsx`**
- Add swipe gesture support for mobile
- Larger navigation buttons on mobile (min 44px touch target)
- Position buttons for easy thumb reach
- Close button with larger touch area

**File: `web/src/components/ImageUploadButton.tsx`**
- Ensure minimum 44px touch target
- Responsive sizing

**File: `web/src/pages/NotFoundPage.tsx`**
- Responsive image sizing
- Adjusted text and spacing for mobile

## Key Mobile-First Principles
1. Base styles for mobile (smallest screens)
2. Use `min-width` breakpoints (`md:`, `lg:`) to scale up
3. Stack vertically by default, side-by-side on larger screens
4. Full width by default, constrained on larger screens
5. Minimum 44x44px touch targets
6. Simplify UI on mobile - show only essential elements

## Color Palette (Dark Theme)
- Backgrounds: `#1a1a1a`, `#2d2d2d`, `#3d3d3d`
- Text: `#ffffff`, `#cccccc`, `#aaaaaa`, `#999999`
- Accents: `#4a9eff` (blue), `#ff4444` (red)
- Borders: `#444444`, `#666666`

---

## Implementation Code

### 1. Update `web/src/index.css`

```css
@import "tailwindcss";

:root {
  font-family: system-ui, Avenir, Helvetica, Arial, sans-serif;
}

body {
  margin: 0;
  min-width: 320px;
  min-height: 100vh;
}

/* Custom scrollbar for dark theme */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #1a1a1a;
}

::-webkit-scrollbar-thumb {
  background: #444;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #555;
}

/* Hide scrollbar on mobile for cleaner look */
@media (max-width: 768px) {
  ::-webkit-scrollbar {
    width: 4px;
    height: 4px;
  }
}

/* Smooth scrolling */
html {
  scroll-behavior: smooth;
}

/* Remove tap highlight on mobile */
* {
  -webkit-tap-highlight-color: transparent;
}

/* Improve touch action on mobile */
button, a, input, textarea, select {
  touch-action: manipulation;
}
```

### 2. Update `web/src/components/Layout.tsx`

```tsx
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
```

### 3. Update `web/src/pages/HomePage.tsx`

```tsx
import { useState, useEffect } from 'react';
import BoardGameCard from '../components/BoardGameCard';

interface BoardGame {
  id: number;
  name: string;
  min_players: number;
  max_players: number;
  coverImageUrl?: string;
}

export default function HomePage() {
  const [games, setGames] = useState<BoardGame[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('/api/boardgames')
      .then(res => res.json())
      .then(data => {
        setGames(data || []);
        setLoading(false);
      })
      .catch(err => {
        console.error(err);
        setLoading(false);
      });
  }, []);

  if (loading) {
    return <div className="text-white text-lg md:text-xl">Loading games...</div>;
  }

  return (
    <div className="h-screen overflow-y-auto">
      <h1 className="mb-4 md:mb-6 lg:mb-8 text-white text-2xl md:text-3xl lg:text-4xl font-bold">
        Board Games
      </h1>

      {games.length === 0 ? (
        <div className="text-center py-10 md:py-15 px-4 text-[#999]">
          <div className="text-6xl md:text-7xl lg:text-8xl mb-4 md:mb-5">🎲</div>
          <h2 className="text-[#ccc] mb-2 text-lg md:text-xl lg:text-2xl">No games yet</h2>
          <p className="text-sm md:text-base">Start building your collection by adding your first board game!</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 md:gap-5 lg:gap-6">
          {games.map(game => (
            <BoardGameCard
              key={game.id}
              id={game.id}
              name={game.name}
              minPlayers={game.min_players}
              maxPlayers={game.max_players}
              coverImageUrl={game.coverImageUrl}
            />
          ))}
        </div>
      )}
    </div>
  );
}
```

### 4. Update `web/src/pages/BoardGameDetailPage.tsx`

```tsx
import { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import ImageGallery from '../components/ImageGallery';

interface BoardGameImageDto {
  id: number;
  imageUrl: string;
  thumbnailUrl: string;
  imageType: string;
  displayOrder: number;
}

interface BoardGame {
  id: number;
  name: string;
  min_players: number;
  max_players: number;
  play_time: number;
  min_age: number;
  description: string;
  created_at: string;
  updated_at: string;
  boardGameImages: BoardGameImageDto[];
}

export default function BoardGameDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [game, setGame] = useState<BoardGame | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [imageLoaded, setImageLoaded] = useState(false);
  const [gameplayImages, setGameplayImages] = useState<BoardGameImageDto[]>([]);

  useEffect(() => {
    fetch(`/api/boardgames/${id}`)
      .then(res => {
        if (!res.ok) {
          throw new Error('Game not found');
        }
        return res.json();
      })
      .then(data => {
        setGame(data);
        const gameplay = data.boardGameImages
          ?.filter((img: BoardGameImageDto) => img.imageType === 'gameplay')
          .sort((a: BoardGameImageDto, b: BoardGameImageDto) => a.displayOrder - b.displayOrder) || [];
        setGameplayImages(gameplay);
        setLoading(false);
      })
      .catch(err => {
        console.error(err);
        setError(true);
        setLoading(false);
      });
  }, [id]);

  const handleDelete = async () => {
    if (!window.confirm(`Are you sure you want to delete "${game?.name}"?`)) {
      return;
    }

    setDeleting(true);

    try {
      const response = await fetch(`/api/boardgames/${id}`, {
        method: 'DELETE'
      });

      if (!response.ok) {
        throw new Error('Failed to delete game');
      }

      navigate('/');
    } catch (err) {
      console.error(err);
      alert('Failed to delete game. Please try again.');
      setDeleting(false);
    }
  };

  if (loading) {
    return <div className="text-white text-lg md:text-xl">Loading game...</div>;
  }

  if (error || !game) {
    return (
      <div className="p-4">
        <h1 className="text-xl md:text-2xl mb-4">Game not found</h1>
        <Link to="/" className="text-[#4a9eff] no-underline text-base md:text-lg">← Back to games</Link>
      </div>
    );
  }

  const coverImage = game.boardGameImages?.find(img => img.imageType === 'cover');

  return (
    <div className="h-screen overflow-y-auto pb-20 md:pb-0">
      <div className="flex justify-between items-center mb-4 md:mb-5">
        <Link to="/" className="text-[#4a9eff] no-underline text-base md:text-lg">
          ← Back to games
        </Link>

        <button
          onClick={handleDelete}
          disabled={deleting}
          className={`border-none text-white text-xl md:text-2xl font-bold cursor-pointer px-3 py-2 rounded-lg transition-all duration-200 flex items-center justify-center min-w-[44px] min-h-[44px] ${
            deleting
              ? 'bg-[#444] cursor-not-allowed'
              : 'bg-[#ff4444] hover:scale-110 active:scale-95'
          }`}
          title="Delete game"
        >
          ✕
        </button>
      </div>

      {/* Main content - stack on mobile, side-by-side on desktop */}
      <div className="flex flex-col md:flex-row gap-6 md:gap-8 lg:gap-10 mt-4 md:mt-5">
        {/* Cover image */}
        <div className="w-full md:w-[350px] lg:w-[400px] md:flex-shrink-0 relative">
          <div className="w-full aspect-square relative">
            {coverImage ? (
              <>
                {!imageLoaded && (
                  <div className="absolute inset-0 rounded-lg animate-pulse bg-[#333]" />
                )}
                <img
                  src={coverImage.imageUrl}
                  alt={`${game.name} cover`}
                  onLoad={() => setImageLoaded(true)}
                  className="w-full h-full object-cover rounded-lg"
                  style={{
                    visibility: imageLoaded ? 'visible' : 'hidden',
                  }}
                />
              </>
            ) : (
              <div className="w-full h-full bg-[#444] rounded-lg flex items-center justify-center text-7xl md:text-8xl lg:text-[120px]">
                🎲
              </div>
            )}
          </div>
        </div>

        {/* Game details */}
        <div className="flex-1">
          <h1 className="text-white mb-4 md:mb-5 text-2xl md:text-3xl lg:text-4xl font-bold">
            {game.name}
          </h1>

          <div className="text-[#ccc] text-sm md:text-base leading-relaxed">
            <p className="mb-2 md:mb-3">
              <strong className="text-white">Players:</strong> {game.min_players}-{game.max_players || game.min_players}
            </p>

            <p className="mb-2 md:mb-3">
              <strong className="text-white">Play Time:</strong> {game.play_time} minutes
            </p>

            <p className="mb-2 md:mb-3">
              <strong className="text-white">Minimum Age:</strong> {game.min_age}+
            </p>

            <p className="mb-3 md:mb-4">
              <strong className="text-white">Description:</strong>
            </p>
            <p className="text-[#aaa] leading-relaxed mb-4 md:mb-5 text-sm md:text-base">
              {game.description}
            </p>

            <p className="text-xs md:text-sm text-[#666]">
              <strong>Added:</strong> {new Date(game.created_at).toLocaleDateString()}
            </p>
            <p className="text-xs md:text-sm text-[#666]">
              <strong>Last Updated:</strong> {new Date(game.updated_at).toLocaleDateString()}
            </p>
          </div>
        </div>
      </div>

      {game && (
        <div className="mt-8 md:mt-10">
          <h2 className="text-white mb-4 md:mb-5 text-xl md:text-2xl font-bold">
            Gameplay Images
          </h2>
          <ImageGallery
            boardGameId={game.id}
            images={gameplayImages}
            onImageUploaded={(newImage) => {
              setGameplayImages(prev => [...prev, newImage]);
            }}
          />
        </div>
      )}
    </div>
  );
}
```

### 5. Update `web/src/pages/AddGamePage.tsx`

```tsx
import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

export default function AddGamePage() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: '',
    min_players: 1,
    max_players: 1,
    play_time: 30,
    min_age: 8,
    description: ''
  });
  const [coverImage, setCoverImage] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      if (!file.type.startsWith('image/')) {
        setError('Please select an image file');
        return;
      }

      if (file.size > 10 * 1024 * 1024) {
        setError('Image must be less than 10MB');
        return;
      }

      setCoverImage(file);

      const reader = new FileReader();
      reader.onloadend = () => {
        setImagePreview(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError('');

    try {
      const gameResponse = await fetch('/api/boardgame', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData)
      });

      if (!gameResponse.ok) {
        throw new Error('Failed to create game');
      }

      const newGame = await gameResponse.json();
      const gameId = newGame.id;

      if (coverImage) {
        const formData = new FormData();
        formData.append('image', coverImage);
        formData.append('imageType', 'cover');

        const imageResponse = await fetch(`/api/boardgame/${gameId}/images`, {
          method: 'POST',
          body: formData
        });

        if (!imageResponse.ok) {
          console.error('Failed to upload image');
        }
      }

      navigate('/');
    } catch {
      setError('Failed to add game. Please try again.');
      setSubmitting(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'name' || name === 'description' ? value : Number(value)
    }));
  };

  return (
    <div className="h-screen overflow-y-auto p-4 md:p-6 pb-20 md:pb-6">
      <Link to="/" className="text-[#4a9eff] no-underline mb-4 md:mb-5 inline-block text-base md:text-lg">
        ← Back to games
      </Link>

      <h1 className="text-white mb-6 md:mb-8 text-2xl md:text-3xl font-bold">Add New Game</h1>

      {error && (
        <div className="bg-[#ff4444] text-white p-3 md:p-4 rounded-md mb-4 md:mb-5 text-sm md:text-base">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="w-full max-w-full md:max-w-2xl lg:max-w-3xl">
        <div className="mb-5 md:mb-6">
          <label className="block text-white mb-2 text-sm md:text-base font-medium">
            Cover Image
          </label>
          <input
            type="file"
            accept="image/*"
            onChange={handleImageChange}
            className="w-full p-3 rounded-md border border-[#444] bg-[#2d2d2d] text-white text-sm md:text-base min-h-[48px]"
          />
          {imagePreview && (
            <div className="mt-3">
              <img
                src={imagePreview}
                alt="Preview"
                className="w-full max-w-[200px] md:max-w-[250px] rounded-md border-2 border-[#444]"
              />
            </div>
          )}
        </div>

        <div className="mb-5 md:mb-6">
          <label className="block text-white mb-2 text-sm md:text-base font-medium">
            Game Name *
          </label>
          <input
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            required
            className="w-full p-3 rounded-md border border-[#444] bg-[#2d2d2d] text-white text-sm md:text-base min-h-[48px]"
          />
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 md:gap-6">
          <div className="mb-5 md:mb-6">
            <label className="block text-white mb-2 text-sm md:text-base font-medium">
              Minimum Players *
            </label>
            <input
              type="number"
              name="min_players"
              value={formData.min_players}
              onChange={handleChange}
              min="1"
              required
              className="w-full p-3 rounded-md border border-[#444] bg-[#2d2d2d] text-white text-sm md:text-base min-h-[48px]"
            />
          </div>

          <div className="mb-5 md:mb-6">
            <label className="block text-white mb-2 text-sm md:text-base font-medium">
              Maximum Players *
            </label>
            <input
              type="number"
              name="max_players"
              value={formData.max_players}
              onChange={handleChange}
              min="1"
              required
              className="w-full p-3 rounded-md border border-[#444] bg-[#2d2d2d] text-white text-sm md:text-base min-h-[48px]"
            />
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 md:gap-6">
          <div className="mb-5 md:mb-6">
            <label className="block text-white mb-2 text-sm md:text-base font-medium">
              Play Time (minutes) *
            </label>
            <input
              type="number"
              name="play_time"
              value={formData.play_time}
              onChange={handleChange}
              min="1"
              required
              className="w-full p-3 rounded-md border border-[#444] bg-[#2d2d2d] text-white text-sm md:text-base min-h-[48px]"
            />
          </div>

          <div className="mb-5 md:mb-6">
            <label className="block text-white mb-2 text-sm md:text-base font-medium">
              Minimum Age *
            </label>
            <input
              type="number"
              name="min_age"
              value={formData.min_age}
              onChange={handleChange}
              min="1"
              required
              className="w-full p-3 rounded-md border border-[#444] bg-[#2d2d2d] text-white text-sm md:text-base min-h-[48px]"
            />
          </div>
        </div>

        <div className="mb-6 md:mb-8">
          <label className="block text-white mb-2 text-sm md:text-base font-medium">
            Description *
          </label>
          <textarea
            name="description"
            value={formData.description}
            onChange={handleChange}
            required
            rows={4}
            className="w-full p-3 rounded-md border border-[#444] bg-[#2d2d2d] text-white text-sm md:text-base resize-y min-h-[120px]"
          />
        </div>

        <button
          type="submit"
          disabled={submitting}
          className={`w-full md:w-auto px-6 md:px-8 py-3 md:py-4 text-white border-none rounded-md text-base md:text-lg cursor-pointer transition-colors duration-200 min-h-[48px] font-medium ${
            submitting
              ? 'bg-[#666] cursor-not-allowed'
              : 'bg-[#4a9eff] hover:bg-[#3a8eef] active:scale-95'
          }`}
        >
          {submitting ? 'Adding...' : 'Add Game'}
        </button>
      </form>
    </div>
  );
}
```

### 6. Update `web/src/components/BoardGameCard.tsx`

```tsx
import { Link } from 'react-router-dom';

interface BoardGameCardProps {
  id: number;
  name: string;
  minPlayers: number;
  maxPlayers: number;
  coverImageUrl?: string;
}

export default function BoardGameCard({ id, name, minPlayers, maxPlayers, coverImageUrl }: BoardGameCardProps) {
  return (
    <Link to={`/boardgame/${id}`} className="no-underline block">
      <div className="bg-[#2d2d2d] rounded-lg overflow-hidden cursor-pointer transition-transform duration-200 w-full hover:scale-[1.02] active:scale-95">
        <div className="w-full aspect-square bg-[#444] flex items-center justify-center text-5xl md:text-6xl overflow-hidden relative">
          {coverImageUrl ? (
            <img
              src={coverImageUrl}
              alt={`${name} cover`}
              className="w-full h-full object-cover object-center"
              onError={(e) => {
                const target = e.currentTarget;
                target.style.display = 'none';
                if (target.parentElement) {
                  target.parentElement.innerHTML = '🎲';
                }
              }}
            />
          ) : (
            '🎲'
          )}
        </div>

        <div className="p-3 md:p-4">
          <h3 className="text-white mb-1 text-sm md:text-base font-medium overflow-hidden text-ellipsis whitespace-nowrap">
            {name}
          </h3>
          <p className="text-[#999] text-xs md:text-sm">
            {minPlayers}-{maxPlayers} players
          </p>
        </div>
      </div>
    </Link>
  );
}
```

### 7. Update `web/src/components/ImageGallery.tsx`

```tsx
import { useState } from 'react';
import ImageUploadButton from './ImageUploadButton';
import Lightbox from './Lightbox';

interface BoardGameImageDto {
  id: number;
  imageUrl: string;
  thumbnailUrl: string;
  imageType: string;
  displayOrder: number;
}

interface ImageGalleryProps {
  boardGameId: number;
  images: BoardGameImageDto[];
  onImageUploaded: (newImage: BoardGameImageDto) => void;
}

export default function ImageGallery({ boardGameId, images, onImageUploaded }: ImageGalleryProps) {
  const [lightboxOpen, setLightboxOpen] = useState(false);
  const [lightboxIndex, setLightboxIndex] = useState(0);
  const [uploading, setUploading] = useState(false);
  const [uploadError, setUploadError] = useState('');
  const [imageLoadStates, setImageLoadStates] = useState<Record<number, boolean>>({});

  const handleImageSelected = async (file: File) => {
    setUploading(true);
    setUploadError('');

    try {
      const formData = new FormData();
      formData.append('image', file);
      formData.append('imageType', 'gameplay');

      const response = await fetch(`/api/boardgame/${boardGameId}/images`, {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Failed to upload image');
      }

      const result = await response.json();

      const newImage: BoardGameImageDto = {
        id: result.imageId,
        imageUrl: `/api/boardgame/${boardGameId}/image/${result.imageId}`,
        thumbnailUrl: `/api/boardgame/${boardGameId}/image/${result.imageId}/thumbnail`,
        imageType: 'gameplay',
        displayOrder: images.length,
      };

      onImageUploaded(newImage);
    } catch {
      setUploadError('Failed to upload image. Please try again.');
    } finally {
      setUploading(false);
    }
  };

  const openLightbox = (index: number) => {
    setLightboxIndex(index);
    setLightboxOpen(true);
  };

  return (
    <div>
      {uploadError && (
        <div className="bg-[#ff4444] text-white p-3 md:p-4 rounded-md mb-4 md:mb-5 text-sm md:text-base">
          {uploadError}
        </div>
      )}

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-5">
        {images.map((image, index) => (
          <div
            key={image.id}
            onClick={() => openLightbox(index)}
            className="relative w-full aspect-square cursor-pointer rounded-lg overflow-hidden transition-transform duration-200 hover:scale-105 active:scale-95"
          >
            {!imageLoadStates[image.id] && (
              <div className="absolute inset-0 rounded-lg animate-pulse bg-[#333]" />
            )}
            <img
              src={image.thumbnailUrl}
              alt={`Gameplay image ${index + 1}`}
              onLoad={() => setImageLoadStates(prev => ({ ...prev, [image.id]: true }))}
              className="w-full h-full object-cover"
              style={{
                visibility: imageLoadStates[image.id] ? 'visible' : 'hidden',
              }}
            />
          </div>
        ))}

        <ImageUploadButton
          onImageSelected={handleImageSelected}
          disabled={uploading}
        />
      </div>

      {lightboxOpen && (
        <Lightbox
          images={images}
          initialIndex={lightboxIndex}
          onClose={() => setLightboxOpen(false)}
        />
      )}
    </div>
  );
}
```

### 8. Update `web/src/components/Lightbox.tsx`

```tsx
import { useState, useEffect, useCallback } from 'react';

interface LightboxProps {
  images: { id: number; imageUrl: string; }[];
  initialIndex: number;
  onClose: () => void;
}

export default function Lightbox({ images, initialIndex, onClose }: LightboxProps) {
  const [currentIndex, setCurrentIndex] = useState(initialIndex);
  const [imageLoaded, setImageLoaded] = useState(false);
  const [touchStart, setTouchStart] = useState<number | null>(null);
  const [touchEnd, setTouchEnd] = useState<number | null>(null);

  const minSwipeDistance = 50;

  const handlePrev = useCallback(() => {
    if (currentIndex > 0) {
      setCurrentIndex(prev => prev - 1);
      setImageLoaded(false);
    }
  }, [currentIndex]);

  const handleNext = useCallback(() => {
    if (currentIndex < images.length - 1) {
      setCurrentIndex(prev => prev + 1);
      setImageLoaded(false);
    }
  }, [currentIndex, images.length]);

  // Keyboard navigation
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose();
      } else if (e.key === 'ArrowLeft') {
        handlePrev();
      } else if (e.key === 'ArrowRight') {
        handleNext();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [currentIndex, images.length, onClose, handlePrev, handleNext]);

  // Touch swipe handlers
  const onTouchStart = (e: React.TouchEvent) => {
    setTouchEnd(null);
    setTouchStart(e.targetTouches[0].clientX);
  };

  const onTouchMove = (e: React.TouchEvent) => {
    setTouchEnd(e.targetTouches[0].clientX);
  };

  const onTouchEnd = () => {
    if (!touchStart || !touchEnd) return;
    const distance = touchStart - touchEnd;
    const isLeftSwipe = distance > minSwipeDistance;
    const isRightSwipe = distance < -minSwipeDistance;

    if (isLeftSwipe) {
      handleNext();
    } else if (isRightSwipe) {
      handlePrev();
    }
  };

  const currentImage = images[currentIndex];

  return (
    <div
      className="fixed inset-0 bg-black/95 z-[1000] flex items-center justify-center p-4 md:p-5"
      onClick={onClose}
      onTouchStart={onTouchStart}
      onTouchMove={onTouchMove}
      onTouchEnd={onTouchEnd}
    >
      {/* Close button */}
      <button
        onClick={onClose}
        className="absolute top-4 right-4 bg-[#ff4444] border-none text-white text-xl md:text-2xl cursor-pointer p-3 rounded-lg z-[1001] min-w-[48px] min-h-[48px] flex items-center justify-center hover:scale-110 active:scale-95 transition-transform"
        aria-label="Close"
      >
        ✕
      </button>

      {/* Previous button */}
      {currentIndex > 0 && (
        <button
          onClick={(e) => {
            e.stopPropagation();
            handlePrev();
          }}
          className="absolute left-2 md:left-5 bg-[#4a9eff] border-none text-white text-xl md:text-2xl cursor-pointer p-3 md:px-4 md:py-3 rounded-lg z-[1001] min-w-[48px] min-h-[48px] flex items-center justify-center hover:scale-110 active:scale-95 transition-transform"
          aria-label="Previous image"
        >
          ←
        </button>
      )}

      {/* Image container */}
      <div
        onClick={(e) => e.stopPropagation()}
        className="relative max-w-[95vw] max-h-[85vh] md:max-w-[90vw] md:max-h-[90vh]"
      >
        {!imageLoaded && (
          <div className="w-[300px] h-[300px] md:w-[400px] md:h-[400px] bg-[#333] rounded-lg animate-pulse" />
        )}
        <img
          src={currentImage.imageUrl}
          alt={`Image ${currentIndex + 1} of ${images.length}`}
          onLoad={() => setImageLoaded(true)}
          className="max-w-[95vw] max-h-[85vh] md:max-w-[90vw] md:max-h-[90vh] object-contain rounded-lg"
          style={{
            visibility: imageLoaded ? 'visible' : 'hidden',
            position: imageLoaded ? 'relative' : 'absolute',
          }}
        />
      </div>

      {/* Next button */}
      {currentIndex < images.length - 1 && (
        <button
          onClick={(e) => {
            e.stopPropagation();
            handleNext();
          }}
          className="absolute right-2 md:right-5 bg-[#4a9eff] border-none text-white text-xl md:text-2xl cursor-pointer p-3 md:px-4 md:py-3 rounded-lg z-[1001] min-w-[48px] min-h-[48px] flex items-center justify-center hover:scale-110 active:scale-95 transition-transform"
          aria-label="Next image"
        >
          →
        </button>
      )}

      {/* Image counter */}
      <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2 text-white text-sm md:text-base bg-black/50 px-3 py-1 rounded-full">
        {currentIndex + 1} / {images.length}
      </div>
    </div>
  );
}
```

### 9. Update `web/src/components/ImageUploadButton.tsx`

```tsx
import { useRef, useState } from 'react';

interface ImageUploadButtonProps {
  onImageSelected: (file: File) => void;
  disabled?: boolean;
}

export default function ImageUploadButton({ onImageSelected, disabled = false }: ImageUploadButtonProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [isDragging, setIsDragging] = useState(false);

  const handleClick = () => {
    if (!disabled && inputRef.current) {
      inputRef.current.click();
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      onImageSelected(file);
    }
    // Reset input value to allow selecting the same file again
    if (inputRef.current) {
      inputRef.current.value = '';
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    if (!disabled) {
      setIsDragging(true);
    }
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    if (!disabled) {
      const file = e.dataTransfer.files?.[0];
      if (file && file.type.startsWith('image/')) {
        onImageSelected(file);
      }
    }
  };

  return (
    <div
      onClick={handleClick}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      className={`
        relative w-full aspect-square
        border-2 border-dashed rounded-lg
        flex flex-col items-center justify-center
        cursor-pointer
        transition-all duration-200
        min-h-[150px]
        ${disabled
          ? 'border-[#444] bg-[#2d2d2d]/50 cursor-not-allowed'
          : isDragging
            ? 'border-[#4a9eff] bg-[#4a9eff]/10'
            : 'border-[#444] bg-[#2d2d2d] hover:border-[#666] hover:bg-[#3d3d3d]'
        }
      `}
    >
      <input
        ref={inputRef}
        type="file"
        accept="image/*"
        onChange={handleFileChange}
        className="hidden"
        disabled={disabled}
      />

      <svg
        className="w-10 h-10 md:w-12 md:h-12 mb-2 text-[#666]"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M12 4v16m8-8H4"
        />
      </svg>

      <span className="text-[#999] text-sm md:text-base text-center px-2">
        {disabled ? 'Uploading...' : 'Add Image'}
      </span>
      <span className="text-[#666] text-xs mt-1 text-center px-2">
        Click or drag & drop
      </span>
    </div>
  );
}
```

### 10. Update `web/src/pages/NotFoundPage.tsx`

```tsx
import { Link } from 'react-router-dom';
import notFoundImage from '../assets/not-found.svg';

export default function NotFoundPage() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] p-4 text-center">
      <img
        src={notFoundImage}
        alt="404 Not Found"
        className="w-full max-w-[250px] md:max-w-[350px] lg:max-w-[400px] h-auto mb-6 md:mb-8"
      />
      <h1 className="text-white text-3xl md:text-4xl lg:text-5xl font-bold mb-3 md:mb-4">
        Page Not Found
      </h1>
      <p className="text-[#999] text-base md:text-lg mb-6 md:mb-8 max-w-md">
        Sorry, the page you're looking for doesn't exist.
      </p>
      <Link
        to="/"
        className="inline-block px-6 md:px-8 py-3 md:py-4 bg-[#4a9eff] text-white rounded-lg text-base md:text-lg font-medium transition-colors duration-200 hover:bg-[#3a8eef] active:scale-95 min-h-[48px]"
      >
        ← Back to Home
      </Link>
    </div>
  );
}
```

---

## Summary of Changes

### Files Modified:
1. `web/src/index.css` - Added mobile optimizations, smooth scrolling, custom scrollbar
2. `web/src/components/Layout.tsx` - Slide-out mobile sidebar with hamburger menu
3. `web/src/pages/HomePage.tsx` - Responsive grid (1→2→3→4 columns)
4. `web/src/pages/BoardGameDetailPage.tsx` - Vertical mobile layout, swipeable content
5. `web/src/pages/AddGamePage.tsx` - Full-width mobile form, touch-friendly inputs
6. `web/src/components/BoardGameCard.tsx` - Fluid width, responsive typography
7. `web/src/components/ImageGallery.tsx` - Responsive grid layout
8. `web/src/components/Lightbox.tsx` - Swipe gestures, larger touch targets, image counter
9. `web/src/components/ImageUploadButton.tsx` - Touch-friendly, drag & drop support
10. `web/src/pages/NotFoundPage.tsx` - Responsive image and typography

### Key Features Implemented:
- ✅ Slide-out sidebar drawer on mobile
- ✅ Smooth 300ms transitions for sidebar
- ✅ 2-column grid on tablet (md: 768px+)
- ✅ Swipe gestures in lightbox (left/right to navigate)
- ✅ Minimum 44px touch targets throughout
- ✅ Responsive typography scaling
- ✅ Mobile-first breakpoints (base → sm → md → lg)
- ✅ Full keyboard navigation support maintained
- ✅ Touch-friendly form inputs with larger padding

### Testing Checklist:
- [ ] Test sidebar open/close on mobile (< 768px)
- [ ] Test sidebar always visible on desktop (≥ 768px)
- [ ] Test grid responsiveness at different widths
- [ ] Test swipe gestures in lightbox on mobile
- [ ] Test keyboard navigation (arrows, escape) in lightbox
- [ ] Test form inputs are easy to tap on mobile
- [ ] Test all buttons have adequate touch targets
- [ ] Test page transitions and animations are smooth

