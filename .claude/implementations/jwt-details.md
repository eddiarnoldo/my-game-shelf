Implementation complete. Here's a summary of everything built:
                                                                                          
  ---
  Backend                                                                                                                                                
                                                                                                                                                         
  - 10 SQL migration files (users, sessions, refresh_tokens, invites tables + admin seed)                                                                
  - src/internal/models/auth.go — User, Session, RefreshToken, Invite structs                                                                            
  - 4 new repositories — UserRepo, SessionRepo, RefreshTokenRepo, InviteRepo                                                                             
  - src/internal/services/email.go — SMTP email service (interface-based)                                                                                
  - src/api/middleware/auth.go — AuthRequired + AdminRequired Gin middleware                                                                             
  - src/api/handlers/auth.go — register, login, refresh, logout, invite handlers with token theft detection                                              
  - Updated router with public vs protected route groups                                                                                                 
  - Updated api.go and main.go to wire everything                                                                                                        
  - Updated CORS middleware to allow Authorization header + handle OPTIONS preflight                                                                     
  - Added golang-jwt/jwt/v5 dependency                                                                                                                   
                                                                                                                                                         
  Frontend                                                                                                                                               
                                                                                                                                                         
  - web/src/context/AuthContext.tsx — login/logout, token storage in sessionStorage, auto-refresh before expiry                                          
  - web/src/components/ProtectedRoute.tsx — redirects to /login if not auth'd, / if not admin
  - 3 new pages: LoginPage, RegisterPage (pre-fills invite code from ?code query param), InvitePage                                                      
  - Updated App.tsx with new routes + AuthProvider wrapper                                                                                               
  - Updated Layout.tsx — conditional nav (Add Game/Invite admin-only), user avatar + logout at sidebar bottom                                            
  - Updated AddGamePage, EditGamePage, BoardGameDetailPage with Bearer auth headers + conditional Edit/Delete UI                                         
                                                                                                                                                         
  To start: Make sure JWT_SECRET in src/.env is set to a real secret (32+ chars) before running. SMTP vars are optional — if left empty, invite creation 
  will fail with an email error but the invite code is still saved to DB (you can retrieve it directly: SELECT code FROM invites).  