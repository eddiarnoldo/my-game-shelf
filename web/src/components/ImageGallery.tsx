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
        <div className="bg-[#ff4444] text-white p-3 rounded-md mb-5">
          {uploadError}
        </div>
      )}

      <div className="grid grid-cols-[repeat(auto-fill,minmax(200px,1fr))] gap-5">
        {images.map((image, index) => (
          <div
            key={image.id}
            onClick={() => openLightbox(index)}
            className="relative w-full aspect-square cursor-pointer rounded-lg overflow-hidden transition-transform duration-200 hover:scale-105"
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
