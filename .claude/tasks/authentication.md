# Authentication for app

Right now the application allows CRUD operations on the boardgames however Edit, Delete and Create should only be allowed for users.

## Goal
We need to implement a login mechanism into our API and also the react SPA application so we can show a logged user and if that logged user is an admin (which will be the case for any logged user the application is capable of showing the EDIT, DELETE and ADD buttons accross the application. 

It also should add some JWT implementation, the idea is that we create an invite system in which we can send an invite to an account email i.e and then that sends the customer back to the application so they can set a password and setup their admin accounts. 

Once they have an account with a username and password our application should have a /login endpoint that will return a JWT with signature and the roles information, for now we can just assign an admin role as the app really only needs that.

If the app detects there is a logged account it should show the buttons but also our backend should validate the JWT in the Go middleware so it can detect it's a valid signature and allow the http request to be resolved only when a JWT valid is present. 

We should also include a refresh token retuned in the login, we want to create a refresh token rotation mechanism however we don't want to allow users to use the refresh token forever, for this we should create a session table that gets related to the user id and the refresh token, we can set a refresh token that expires in 1 day and a session that expires in 3 days, that way at most they will be able to use refresh tokens for 3 days until the session expires.

## New Tables to create
- User (stores information about users i.e email, username, password, created_at) need to store password using bcrypt
- Session (stores a user session and it's expires_at will be used to determine when to force a user to login in case session alread expired)
- RefreshTokens (has a relation to session, stores the refresh tokens hashed that we give each user session and revokes old tokens on /refresh call)
- Invites 
  - invites:
  - id
  - code (unique)
  - created_by (user_id)
  - used (boolean)
  - used_by (user_id, nullable)
  - expires_at (optional but recommended)

## Invite logic
Admin users should be able to create invites, once logged in a new page, send invite should show up, this page should be minimal it will display an email input and send invite, once clicked the app will perform a request to backend and generate an invite using a code and saving it to DB (any new invite for an already existing invite should delete the previous one)

Once the user gets an email it should have an invite code, the app will have a new page `Join` which takes in an invite code, username and password + password confirmation, if the invite is valid it should add the user into the users page.


## New endpoints
- `/invites` Generates a random code, Saves it in DB, sends email to the email send on the request, returns a 200 status code and displays a UI message
- `/register` Find invite by code, Check: exists, not used, not expired, Create user with password provided, Mark invite as used
- `/login` takes username and password and returns access_token and refresh_token for that user , it also creates a session for that user and sets a expire at
- `/refresh` takes in a refresh_token and returns a new access_token and refresh_token as long as the session saved for that refresh token is not expired
- `/logout` takes in a refresh_token and revokes it for that session, it should also delete the session so user is forced to login

## UI/Frontend changes
- Need to create a login page, this page takes in a username and a password once logged in it takes them to the home page.
- Update the header to show a logged username info and placeholder avatar img showing a u
- Update pages to show/not show edit,delete and add new game links depending if user is logged in
- Implement login/logout logic to either request and save the access_token in session or logout an delete access_token, refresh_token from session
- Create invite page only seen by admin, this page takes an email an performs a backend request to send a email with the invite code
- Create a register page, this page takes as input an invite code a username, password & password confirmation, performs a request to backend and upon confirmation creates a user and returns it back to the main page or login (ask just registered user where do they want to go)

## Migrations
- Add new migrations for new tables required, users, sessions and refreshtoken tables
- Add a new migration that adds a default `admin` user with password `mygameshelf12345##@@` this will be useful for new installations so new users can register the first own admin account

## New services
- Need to create a mail service so the app can send an email to a external service i.e the app should be able to send an invite link to a specificed URL


## Update to existing endpoints
Create, Update and Delete endpoints now need authorization to function so we need to modify the Go app middleware to validate the jwt is present and it has and admin role, then we can resolve the request as normal.

