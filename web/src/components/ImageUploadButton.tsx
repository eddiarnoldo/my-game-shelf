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

    if (!file.type.startsWith('image/')) {
      setError('Please select an image file');
      return;
    }

    if (file.size > 10 * 1024 * 1024) {
      setError('Image must be less than 10MB');
      return;
    }

    setError('');
    onImageSelected(file);

    e.target.value = '';
  };

  return (
    <div>
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        onChange={handleFileChange}
        className="hidden"
      />
      <button
        onClick={handleClick}
        disabled={disabled}
        className={`w-full aspect-square border-2 border-dashed border-[#666] rounded-lg text-[#999] text-5xl cursor-pointer transition-all duration-200 ${
          disabled 
            ? 'bg-[#333] cursor-not-allowed' 
            : 'bg-[#2d2d2d] hover:bg-[#3d3d3d]'
        }`}
      >
        +
      </button>
      {error && (
        <div className="text-[#ff4444] text-sm mt-2">
          {error}
        </div>
      )}
    </div>
  );
}
