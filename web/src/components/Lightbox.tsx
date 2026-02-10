import { useState, useEffect } from 'react';

interface LightboxProps {
  images: { id: number; imageUrl: string; }[];
  initialIndex: number;
  onClose: () => void;
}

export default function Lightbox({ images, initialIndex, onClose }: LightboxProps) {
  const [currentIndex, setCurrentIndex] = useState(initialIndex);
  const [imageLoaded, setImageLoaded] = useState(false);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose();
      } else if (e.key === 'ArrowLeft' && currentIndex > 0) {
        setCurrentIndex(prev => prev - 1);
        setImageLoaded(false);
      } else if (e.key === 'ArrowRight' && currentIndex < images.length - 1) {
        setCurrentIndex(prev => prev + 1);
        setImageLoaded(false);
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [currentIndex, images.length, onClose]);

  const currentImage = images[currentIndex];

  return (
    <div
      className="fixed inset-0 bg-black/95 z-[1000] flex items-center justify-center p-5"
      onClick={onClose}
    >
      <button
        onClick={onClose}
        className="absolute top-5 right-5 bg-[#ff4444] border-none text-white text-2xl cursor-pointer px-3 py-2 rounded-lg z-[1001]"
      >
        ✕
      </button>

      {currentIndex > 0 && (
        <button
          onClick={(e) => {
            e.stopPropagation();
            setCurrentIndex(prev => prev - 1);
            setImageLoaded(false);
          }}
          className="absolute left-5 bg-[#4a9eff] border-none text-white text-2xl cursor-pointer px-4 py-3 rounded-lg z-[1001]"
        >
          ←
        </button>
      )}

      <div onClick={(e) => e.stopPropagation()} className="relative max-w-[90vw] max-h-[90vh]">
        {!imageLoaded && (
          <div className="w-[400px] h-[400px] bg-[#333] rounded-lg animate-pulse" />
        )}
        <img
          src={currentImage.imageUrl}
          alt={`Image ${currentIndex + 1}`}
          onLoad={() => setImageLoaded(true)}
          className="max-w-[90vw] max-h-[90vh] object-contain rounded-lg"
          style={{
            visibility: imageLoaded ? 'visible' : 'hidden',
            position: imageLoaded ? 'relative' : 'absolute',
          }}
        />
      </div>

      {currentIndex < images.length - 1 && (
        <button
          onClick={(e) => {
            e.stopPropagation();
            setCurrentIndex(prev => prev + 1);
            setImageLoaded(false);
          }}
          className="absolute right-5 bg-[#4a9eff] border-none text-white text-2xl cursor-pointer px-4 py-3 rounded-lg z-[1001]"
        >
          →
        </button>
      )}
    </div>
  );
}
