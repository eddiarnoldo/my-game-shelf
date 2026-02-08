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
      style={{
        position: 'fixed',
        inset: 0,
        backgroundColor: 'rgba(0, 0, 0, 0.95)',
        zIndex: 1000,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        padding: '20px',
      }}
      onClick={onClose}
    >
      {/* Close button */}
      <button
        onClick={onClose}
        style={{
          position: 'absolute',
          top: '20px',
          right: '20px',
          backgroundColor: '#ff4444',
          border: 'none',
          color: 'white',
          fontSize: '24px',
          cursor: 'pointer',
          padding: '8px 12px',
          borderRadius: '8px',
          zIndex: 1001,
        }}
      >
        ✕
      </button>

      {/* Previous button */}
      {currentIndex > 0 && (
        <button
          onClick={(e) => {
            e.stopPropagation();
            setCurrentIndex(prev => prev - 1);
            setImageLoaded(false);
          }}
          style={{
            position: 'absolute',
            left: '20px',
            backgroundColor: '#4a9eff',
            border: 'none',
            color: 'white',
            fontSize: '24px',
            cursor: 'pointer',
            padding: '12px 16px',
            borderRadius: '8px',
            zIndex: 1001,
          }}
        >
          ←
        </button>
      )}

      {/* Image container */}
      <div onClick={(e) => e.stopPropagation()} style={{ position: 'relative', maxWidth: '90vw', maxHeight: '90vh' }}>
        {!imageLoaded && (
          <div style={{
            width: '400px',
            height: '400px',
            backgroundColor: '#333',
            borderRadius: '8px',
            animation: 'pulse 1.5s ease-in-out infinite',
          }} />
        )}
        <img
          src={currentImage.imageUrl}
          alt={`Image ${currentIndex + 1}`}
          onLoad={() => setImageLoaded(true)}
          style={{
            maxWidth: '90vw',
            maxHeight: '90vh',
            objectFit: 'contain',
            borderRadius: '8px',
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
            setCurrentIndex(prev => prev + 1);
            setImageLoaded(false);
          }}
          style={{
            position: 'absolute',
            right: '20px',
            backgroundColor: '#4a9eff',
            border: 'none',
            color: 'white',
            fontSize: '24px',
            cursor: 'pointer',
            padding: '12px 16px',
            borderRadius: '8px',
            zIndex: 1001,
          }}
        >
          →
        </button>
      )}
    </div>
  );
}
