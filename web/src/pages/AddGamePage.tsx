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
    <div className="h-screen overflow-y-auto p-5">
      <Link to="/" className="text-[#4a9eff] no-underline mb-5 inline-block">
        ← Back to games
      </Link>

      <h1 className="text-white mb-[30px]">Add New Game</h1>

      {error && (
        <div className="bg-[#ff4444] text-white p-3 rounded-md mb-5">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="max-w-[600px]">
        <div className="mb-5">
          <label className="block text-white mb-2">
            Cover Image
          </label>
          <input
            type="file"
            accept="image/*"
            onChange={handleImageChange}
            className="w-full p-2.5 rounded-md border border-[#444] bg-[#2d2d2d] text-white text-base"
          />
          {imagePreview && (
            <div className="mt-3">
              <img 
                src={imagePreview} 
                alt="Preview" 
                className="max-w-[200px] max-h-[200px] rounded-md border-2 border-[#444]"
              />
            </div>
          )}
        </div>

        <div className="mb-5">
          <label className="block text-white mb-2">
            Game Name *
          </label>
          <input
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            required
            className="w-full p-2.5 rounded-md border border-[#444] bg-[#2d2d2d] text-white text-base"
          />
        </div>
        
        <div className="mb-5">
          <label className="block text-white mb-2">
            Minimum Players *
          </label>
          <input
            type="number"
            name="min_players"
            value={formData.min_players}
            onChange={handleChange}
            min="1"
            required
            className="w-full p-2.5 rounded-md border border-[#444] bg-[#2d2d2d] text-white text-base"
          />
        </div>

        <div className="mb-5">
          <label className="block text-white mb-2">
            Maximum Players *
          </label>
          <input
            type="number"
            name="max_players"
            value={formData.max_players}
            onChange={handleChange}
            min="1"
            required
            className="w-full p-2.5 rounded-md border border-[#444] bg-[#2d2d2d] text-white text-base"
          />
        </div>

        <div className="mb-5">
          <label className="block text-white mb-2">
            Play Time (minutes) *
          </label>
          <input
            type="number"
            name="play_time"
            value={formData.play_time}
            onChange={handleChange}
            min="1"
            required
            className="w-full p-2.5 rounded-md border border-[#444] bg-[#2d2d2d] text-white text-base"
          />
        </div>

        <div className="mb-5">
          <label className="block text-white mb-2">
            Minimum Age *
          </label>
          <input
            type="number"
            name="min_age"
            value={formData.min_age}
            onChange={handleChange}
            min="1"
            required
            className="w-full p-2.5 rounded-md border border-[#444] bg-[#2d2d2d] text-white text-base"
          />
        </div>

        <div className="mb-5">
          <label className="block text-white mb-2">
            Description *
          </label>
          <textarea
            name="description"
            value={formData.description}
            onChange={handleChange}
            required
            rows={4}
            className="w-full p-2.5 rounded-md border border-[#444] bg-[#2d2d2d] text-white text-base font-[inherit] resize-y"
          />
        </div>

        <button
          type="submit"
          disabled={submitting}
          className={`px-6 py-3 text-white border-none rounded-md text-base cursor-pointer transition-colors duration-200 ${
            submitting 
              ? 'bg-[#666] cursor-not-allowed' 
              : 'bg-[#4a9eff] hover:bg-[#3a8eef]'
          }`}
        >
          {submitting ? 'Adding...' : 'Add Game'}
        </button>
      </form>
    </div>
  );
}
