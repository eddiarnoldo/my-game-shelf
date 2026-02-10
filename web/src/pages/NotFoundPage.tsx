import { Link } from 'react-router-dom';

export default function NotFoundPage() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] p-4 text-center">
      <img
        src="/not-found.png"
        alt="404 Not Found"
        className="w-full max-w-[250px] md:max-w-[350px] lg:max-w-[400px] h-auto mb-6 md:mb-8"
      />
      <h1 className="text-white text-3xl md:text-4xl lg:text-5xl font-bold mb-3 md:mb-4">
        Page Not Found
      </h1>
      <p className="text-[#999] text-base md:text-lg mb-6 md:mb-8 max-w-md">
        Sorry, the page you're looking for doesn't exist.
      </p>
      <Link
        to="/"
        className="inline-block px-6 md:px-8 py-3 md:py-4 bg-[#4a9eff] text-white rounded-lg text-base md:text-lg font-medium transition-colors duration-200 hover:bg-[#3a8eef] active:scale-95 min-h-[48px]"
      >
        ← Back to Home
      </Link>
    </div>
  );
}
