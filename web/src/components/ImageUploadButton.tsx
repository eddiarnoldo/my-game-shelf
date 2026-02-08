import { useRef, useState } from 'react';

interface ImageUploadButtonProps {
  onImageSelected: (file: File) => void;
  disabled?: boolean;
}

export default function ImageUploadButton({ onImageSelected, disabled }: ImageUploadButtonProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [error, setError] = useState('');

  const handleClick = () => {
    if (!disabled) {
      fileInputRef.current?.click();
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Validate file type
    if (!file.type.startsWith('image/')) {
      setError('Please select an image file');
      return;
    }

    // Validate file size (max 10MB)
    if (file.size > 10 * 1024 * 1024) {
      setError('Image must be less than 10MB');
      return;
    }

    setError('');
    onImageSelected(file);

    // Reset input so same file can be selected again
    e.target.value = '';
  };

  return (
    <div>
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        onChange={handleFileChange}
        style={{ display: 'none' }}
      />
      <button
        onClick={handleClick}
        disabled={disabled}
        style={{
          width: '100%',
          aspectRatio: '1',
          backgroundColor: disabled ? '#333' : '#2d2d2d',
          border: '2px dashed #666',
          borderRadius: '8px',
          color: '#999',
          fontSize: '48px',
          cursor: disabled ? 'not-allowed' : 'pointer',
          transition: 'all 0.2s',
        }}
        onMouseEnter={(e) => {
          if (!disabled) {
            e.currentTarget.style.backgroundColor = '#3d3d3d';
          }
        }}
        onMouseLeave={(e) => {
          if (!disabled) {
            e.currentTarget.style.backgroundColor = '#2d2d2d';
          }
        }}
      >
        +
      </button>
      {error && (
        <div style={{
          color: '#ff4444',
          fontSize: '14px',
          marginTop: '8px'
        }}>
          {error}
        </div>
      )}
    </div>
  );
}
