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
