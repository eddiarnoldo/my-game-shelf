import { useRef, useState } from 'react';

interface ImageUploadButtonProps {
  onImageSelected: (file: File) => void;
  disabled?: boolean;
}

export default function ImageUploadButton({ onImageSelected, disabled = false }: ImageUploadButtonProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [error, setError] = useState('');

  const handleClick = () => {
    if (!disabled && inputRef.current) {
      inputRef.current.click();
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
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

      setError('');
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
        setError('');
        onImageSelected(file);
      } else if (file) {
        setError('Please select an image file');
      }
    }
  };

  return (
    <div>
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
      {error && (
        <div className="bg-[#ff4444] text-white p-2 rounded-md mt-2 text-sm">
          {error}
        </div>
      )}
    </div>
  );
}
