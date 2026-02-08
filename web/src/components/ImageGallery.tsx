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

      // Construct new image DTO with returned imageId
      const newImage: BoardGameImageDto = {
        id: result.imageId,
        imageUrl: `/api/boardgame/${boardGameId}/image/${result.imageId}`,
        thumbnailUrl: `/api/boardgame/${boardGameId}/image/${result.imageId}/thumbnail`,
        imageType: 'gameplay',
        displayOrder: images.length,
      };

      onImageUploaded(newImage);
    } catch (err) {
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
        <div style={{
          backgroundColor: '#ff4444',
          color: 'white',
          padding: '12px',
          borderRadius: '6px',
          marginBottom: '20px'
        }}>
          {uploadError}
        </div>
      )}

      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))',
        gap: '20px',
      }}>
        {images.map((image, index) => (
          <div
            key={image.id}
            onClick={() => openLightbox(index)}
            style={{
              position: 'relative',
              width: '100%',
              aspectRatio: '1',
              cursor: 'pointer',
              borderRadius: '8px',
              overflow: 'hidden',
              transition: 'transform 0.2s',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.transform = 'scale(1.05)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.transform = 'scale(1)';
            }}
          >
            {!imageLoadStates[image.id] && (
              <div style={{
                position: 'absolute',
                inset: 0,
                borderRadius: '8px',
                animation: 'pulse 1.5s ease-in-out infinite',
                backgroundColor: '#333',
              }} />
            )}
            <img
              src={image.thumbnailUrl}
              alt={`Gameplay image ${index + 1}`}
              onLoad={() => setImageLoadStates(prev => ({ ...prev, [image.id]: true }))}
              style={{
                width: '100%',
                height: '100%',
                objectFit: 'cover',
                visibility: imageLoadStates[image.id] ? 'visible' : 'hidden',
              }}
            />
          </div>
        ))}

        {/* Upload button */}
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
