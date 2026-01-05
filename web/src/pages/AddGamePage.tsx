import { useState } from 'react';
   import { Link, useNavigate } from 'react-router-dom';

   export default function AddGamePage() {
     const navigate = useNavigate();
     const [formData, setFormData] = useState({
       name: '',
       min_players: 1,
       max_players: 1,
       play_time: 30,
       min_age: 8,
       description: ''
     });
     const [submitting, setSubmitting] = useState(false);
     const [error, setError] = useState('');

     const handleSubmit = async (e: React.FormEvent) => {
       e.preventDefault();
       setSubmitting(true);
       setError('');

       try {
         const response = await fetch('/api/boardgame', {
           method: 'POST',
           headers: {
             'Content-Type': 'application/json',
           },
           body: JSON.stringify(formData)
         });

         if (!response.ok) {
           throw new Error('Failed to create game');
         }

         // Success! Navigate back to home
         navigate('/');
       } catch (err) {
         setError('Failed to add game. Please try again.');
         setSubmitting(false);
       }
     };

     const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
       const { name, value } = e.target;
       setFormData(prev => ({
         ...prev,
         [name]: name === 'name' || name === 'description' ? value : Number(value)
       }));
     };

     return (
       <div>
         <Link to="/" style={{ color: '#4a9eff', textDecoration: 'none', marginBottom: '20px', display: 'inline-block' }}>
           ← Back to games
         </Link>

         <h1 style={{ color: 'white', marginBottom: '30px' }}>Add New Game</h1>

         {error && (
           <div style={{ 
             backgroundColor: '#ff4444', 
             color: 'white', 
             padding: '12px', 
             borderRadius: '6px',
             marginBottom: '20px'
           }}>
             {error}
           </div>
         )}

         <form onSubmit={handleSubmit} style={{ maxWidth: '600px' }}>
           {/* Game Name */}
           <div style={{ marginBottom: '20px' }}>
             <label style={{ display: 'block', color: 'white', marginBottom: '8px' }}>
               Game Name *
             </label>
             <input
               type="text"
               name="name"
               value={formData.name}
               onChange={handleChange}
               required
               style={{
                 width: '100%',
                 padding: '10px',
                 borderRadius: '6px',
                 border: '1px solid #444',
                 backgroundColor: '#2d2d2d',
                 color: 'white',
                 fontSize: '16px'
               }}
             />
           </div>

           {/* Min Players */}
           <div style={{ marginBottom: '20px' }}>
             <label style={{ display: 'block', color: 'white', marginBottom: '8px' }}>
               Minimum Players *
             </label>
             <input
               type="number"
               name="min_players"
               value={formData.min_players}
               onChange={handleChange}
               min="1"
               required
               style={{
                 width: '100%',
                 padding: '10px',
                 borderRadius: '6px',
                 border: '1px solid #444',
                 backgroundColor: '#2d2d2d',
                 color: 'white',
                 fontSize: '16px'
               }}
             />
           </div>

           {/* Max Players */}
           <div style={{ marginBottom: '20px' }}>
             <label style={{ display: 'block', color: 'white', marginBottom: '8px' }}>
               Maximum Players *
             </label>
             <input
               type="number"
               name="max_players"
               value={formData.max_players}
               onChange={handleChange}
               min="1"
               required
               style={{
                 width: '100%',
                 padding: '10px',
                 borderRadius: '6px',
                 border: '1px solid #444',
                 backgroundColor: '#2d2d2d',
                 color: 'white',
                 fontSize: '16px'
               }}
             />
           </div>

           {/* Play Time */}
           <div style={{ marginBottom: '20px' }}>
             <label style={{ display: 'block', color: 'white', marginBottom: '8px' }}>
               Play Time (minutes) *
             </label>
             <input
               type="number"
               name="play_time"
               value={formData.play_time}
               onChange={handleChange}
               min="1"
               required
               style={{
                 width: '100%',
                 padding: '10px',
                 borderRadius: '6px',
                 border: '1px solid #444',
                 backgroundColor: '#2d2d2d',
                 color: 'white',
                 fontSize: '16px'
               }}
             />
           </div>

           {/* Min Age */}
           <div style={{ marginBottom: '20px' }}>
             <label style={{ display: 'block', color: 'white', marginBottom: '8px' }}>
               Minimum Age *
             </label>
             <input
               type="number"
               name="min_age"
               value={formData.min_age}
               onChange={handleChange}
               min="1"
               required
               style={{
                 width: '100%',
                 padding: '10px',
                 borderRadius: '6px',
                 border: '1px solid #444',
                 backgroundColor: '#2d2d2d',
                 color: 'white',
                 fontSize: '16px'
               }}
             />
           </div>

           {/* Description */}
           <div style={{ marginBottom: '20px' }}>
             <label style={{ display: 'block', color: 'white', marginBottom: '8px' }}>
               Description *
             </label>
             <textarea
               name="description"
               value={formData.description}
               onChange={handleChange}
               required
               rows={4}
               style={{
                 width: '100%',
                 padding: '10px',
                 borderRadius: '6px',
                 border: '1px solid #444',
                 backgroundColor: '#2d2d2d',
                 color: 'white',
                 fontSize: '16px',
                 fontFamily: 'inherit',
                 resize: 'vertical'
               }}
             />
           </div>

           {/* Submit Button */}
           <button
             type="submit"
             disabled={submitting}
             style={{
               padding: '12px 24px',
               backgroundColor: submitting ? '#666' : '#4a9eff',
               color: 'white',
               border: 'none',
               borderRadius: '6px',
               fontSize: '16px',
               cursor: submitting ? 'not-allowed' : 'pointer',
               transition: 'background-color 0.2s'
             }}
           >
             {submitting ? 'Adding...' : 'Add Game'}
           </button>
         </form>
       </div>
     );
   }