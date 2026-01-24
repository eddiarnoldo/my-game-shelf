import { Link } from "react-router-dom";

export default function NotFoundPage() {
  return (
    <div style={{ padding: '2rem' }}>
      <h2>404 - Page Not Found</h2>
      <p>The page you’re looking for doesn’t exist.</p>
      <Link to="/">Go back home</Link>
      <img 
        src="/not-found.png"
        alt="404 Not Found" 
        style={{ marginTop: '20px', maxWidth: '100%', height: 'auto' }} 
      />
    </div>
  );
}