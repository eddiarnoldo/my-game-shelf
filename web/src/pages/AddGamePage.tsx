import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import ImageUploadButton from '../components/ImageUploadButton';
import { useAuth } from '../context/AuthContext';

export default function AddGamePage() {
  const navigate = useNavigate();
  const { accessToken } = useAuth();
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

  useEffect(() => {
    if (!coverImage) { setImagePreview(null); return; }
    const reader = new FileReader();
    reader.onloadend = () => setImagePreview(reader.result as string);
    reader.readAsDataURL(coverImage);
  }, [coverImage]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError('');

    try {
      const gameResponse = await fetch('/api/boardgame', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${accessToken}`,
        },
        body: JSON.stringify(formData)
      });

      if (!gameResponse.ok) {
        throw new Error('Failed to create game');
      }

      const newGame = await gameResponse.json();
      const gameId = newGame.id;

      if (coverImage) {
        const imageFormData = new FormData();
        imageFormData.append('image', coverImage);
        imageFormData.append('imageType', 'cover');

        const imageResponse = await fetch(`/api/boardgame/${gameId}/images`, {
          method: 'POST',
          headers: { 'Authorization': `Bearer ${accessToken}` },
          body: imageFormData
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

  const inputClass = "w-full p-3 rounded-lg border border-[#444] bg-[#1a1a1a] text-white text-sm md:text-base min-h-[48px] focus:outline-none focus:ring-2 focus:ring-[#4a9eff] focus:border-[#4a9eff] transition-colors duration-200 placeholder:text-[#666]";
  const labelClass = "block text-[#ccc] mb-2 text-sm font-medium";
  const sectionLabelClass = "text-[#999] text-xs uppercase tracking-widest mb-4 font-semibold";
  const cardClass = "bg-[#2d2d2d] rounded-xl p-5 md:p-6";

  return (
    <div className="h-screen overflow-y-auto p-4 md:p-6 pb-20 md:pb-6">
      <Link to="/" className="text-[#4a9eff] no-underline mb-4 md:mb-5 inline-block text-base md:text-lg pl-14 md:pl-0">
        ← Back to games
      </Link>

      <h1 className="text-white mb-6 md:mb-8 text-2xl md:text-3xl font-bold">Add New Game</h1>

      {error && (
        <div className="bg-[#ff4444] text-white p-3 md:p-4 rounded-md mb-4 md:mb-5 text-sm md:text-base">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="w-full max-w-5xl">
        <div className="flex flex-col lg:flex-row gap-6 lg:gap-8 items-start">

          {/* Left column — Cover image card */}
          <div className="w-full lg:w-[300px] flex-shrink-0 lg:sticky lg:top-0">
            <div className={cardClass}>
              <p className={sectionLabelClass}>🖼 Cover Image</p>
              {imagePreview ? (
                <div>
                  <div className="relative">
                    <img
                      src={imagePreview}
                      alt="Cover preview"
                      className="w-full aspect-square object-cover rounded-lg"
                    />
                    <button
                      type="button"
                      onClick={() => setCoverImage(null)}
                      className="absolute top-2 right-2 w-7 h-7 rounded-full bg-[#1a1a1a]/80 text-white flex items-center justify-center text-sm hover:bg-[#ff4444] transition-colors duration-200"
                      aria-label="Remove image"
                    >
                      ✕
                    </button>
                  </div>
                  {coverImage && (
                    <p className="text-[#666] text-xs mt-2 truncate">{coverImage.name}</p>
                  )}
                </div>
              ) : (
                <ImageUploadButton onImageSelected={(f) => setCoverImage(f)} disabled={submitting} />
              )}
            </div>
          </div>

          {/* Right column — Form section cards */}
          <div className="flex-1 flex flex-col gap-5">

            {/* Card 1 — Game Details */}
            <div className={cardClass}>
              <p className={sectionLabelClass}>🎮 Game Details</p>
              <div>
                <label className={labelClass}>
                  Game Name <span className="text-[#ff4444]">*</span>
                </label>
                <input
                  type="text"
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  required
                  placeholder="e.g. Wingspan"
                  className={inputClass}
                />
              </div>
            </div>

            {/* Card 2 — Players & Time */}
            <div className={cardClass}>
              <p className={sectionLabelClass}>👥 Players &amp; Time</p>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className={labelClass}>
                    Min Players <span className="text-[#ff4444]">*</span>
                  </label>
                  <input
                    type="number"
                    name="min_players"
                    value={formData.min_players}
                    onChange={handleChange}
                    min="1"
                    required
                    className={inputClass}
                  />
                </div>
                <div>
                  <label className={labelClass}>
                    Max Players <span className="text-[#ff4444]">*</span>
                  </label>
                  <input
                    type="number"
                    name="max_players"
                    value={formData.max_players}
                    onChange={handleChange}
                    min="1"
                    required
                    className={inputClass}
                  />
                </div>
                <div>
                  <label className={labelClass}>
                    Play Time (min) <span className="text-[#ff4444]">*</span>
                  </label>
                  <input
                    type="number"
                    name="play_time"
                    value={formData.play_time}
                    onChange={handleChange}
                    min="1"
                    required
                    className={inputClass}
                  />
                </div>
                <div>
                  <label className={labelClass}>
                    Min Age <span className="text-[#ff4444]">*</span>
                  </label>
                  <input
                    type="number"
                    name="min_age"
                    value={formData.min_age}
                    onChange={handleChange}
                    min="1"
                    required
                    className={inputClass}
                  />
                </div>
              </div>
            </div>

            {/* Card 3 — Description */}
            <div className={cardClass}>
              <p className={sectionLabelClass}>📝 Description</p>
              <div>
                <label className={labelClass}>
                  Description <span className="text-[#ff4444]">*</span>
                </label>
                <textarea
                  name="description"
                  value={formData.description}
                  onChange={handleChange}
                  required
                  rows={5}
                  className={`${inputClass} resize-y min-h-[120px]`}
                />
              </div>
            </div>

            {/* Submit button */}
            <button
              type="submit"
              disabled={submitting}
              className={`w-full py-4 rounded-xl font-semibold text-base min-h-[52px] text-white border-none cursor-pointer transition-colors duration-200 active:scale-[0.98] ${
                submitting
                  ? 'bg-[#444] text-[#999] cursor-not-allowed'
                  : 'bg-[#4a9eff] hover:bg-[#3a8eef] shadow-lg shadow-[#4a9eff]/20'
              }`}
            >
              {submitting ? 'Adding...' : 'Add Game to Shelf'}
            </button>

          </div>
        </div>
      </form>
    </div>
  );
}
