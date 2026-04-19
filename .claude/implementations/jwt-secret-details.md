JWT_SECRET is the private key used to sign and verify JWTs.                                                                                            
                                                                                                                                                         
  When a user logs in, the server creates a token payload like:                                                                                          
                                                                                                                                                         
  {
    "user_id": 1,                                                                                                                                        
    "username": "admin",                                                                                                                                 
    "role": "admin",                     
    "exp": 1712345678                                                                                                                                    
  }                                                                                                                                                      
                                         
  That payload is then cryptographically signed using JWT_SECRET with HMAC-SHA256. The result is the third part of the JWT:                              
                                                                                     
  header.payload.SIGNATURE                                                                                                                               
                                                                                     
  The signature is computed as:                                                                                                                          
   
  HMAC-SHA256(base64(header) + "." + base64(payload), JWT_SECRET)                                                                                        
                                                                                                                                                         
  Why this matters:                      
                                                                                                                                                         
  - Anyone can read a JWT — the payload is just base64, not encrypted                                                                                    
  - But only the server (holding JWT_SECRET) can produce a valid signature
  - When a request arrives with a JWT, the server re-computes the signature from the payload and compares it. If they match → the token is genuine and   
  untampered. If someone edits the payload (e.g. changes "role": "user" to "role": "admin"), the signature won't match → rejected                        
                                                                                                                                                         
  The risk of a weak or leaked secret:                                                                                                                   
                                                                                     
  If JWT_SECRET is guessed or stolen, an attacker can forge tokens with any claims they want — including "role": "admin" — and the server will accept    
  them as valid. This is why:                                                        
                                                                                                                                                         
  1. It must be random and long (32+ bytes)                                                                                                              
  2. It must never be committed to git — only live in .env
  3. If leaked, rotate it immediately — all existing tokens become invalid (which is actually a clean revocation mechanism)                              
                                                                                                                                                         
  In the codebase this happens in src/api/middleware/auth.go:31-35:                                                                                      
                                                                                                                                                         
  token, err := jwt.ParseWithClaims(tokenStr, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {                                                   
      if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {                                                                                               
          return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])                                                                       
      }                                                                                                                                                  
      return []byte(jwtSecret), nil  // ← this is where the secret is used to verify                                                                     
  }) 


# How can I generate it?                                                                                                        
                                                                                                                                                         
  openssl (simplest):                                                                                                                                    
  openssl rand -base64 32                                                                                                                                
                                                                                                                                                         
  /dev/urandom via xxd:                                                                                                                                  
  head -c 32 /dev/urandom | xxd -p | tr -d '\n'                                       
                                                                                                                                                         
  Python (built-in):                                                                                                                                     
  python3 -c "import secrets; print(secrets.token_hex(32))"
                                                                                                                                                         
  Run any of these in terminal, copy the output, paste it into src/.env:                                                                                 
                                                                                                                                                         
  JWT_SECRET=your-generated-value-here                                                                                                                   
                                                                                                                                                         
  The openssl one is the most common. Output looks like: K8mP3xQz9vL2nR7wY1cA5jF6hE0tB4sD+uG8iN=  