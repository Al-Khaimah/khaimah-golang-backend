# API Endpoints

## 1. Users Module
Handles user authentication, profile management, and preferences.

- **POST /auth/signup**  
  Create a new user account.
  
- **POST /auth/login**  
  Authenticate user and return a JWT or session.
  
- **POST /auth/logout**  
  Logout the user (invalidate the JWT or session).
  
- **POST /auth/refresh-token**  
  Refresh the authentication token.
  
- **GET /user/profile**  
  Fetch the current user's profile details.
  
- **PUT /user/profile**  
  Update user's profile details.
  
- **PATCH /user/profile/password**  
  Allow users to change their password.
  
- **PATCH /user/profile/preferences**  
  Update a user's preferences (e.g., change categories they are interested in).

---

## 2. Categories Module
Handles category-related operations.

- **GET /categories**  
  Fetch all categories.

- **GET /user/categories/{id}/podcasts**  
  Fetch podcasts by the categories the user follows.

- **GET /trending/categories**  
  List the most popular categories based on user interactions.

---

## 3. Podcasts Module
Handles podcast-related operations such as listing, liking, playing, downloading podcasts, etc.

- **GET /podcasts**  
  List all podcasts filtered by category, popularity, etc.

- **GET /podcasts/{id}**  
  Fetch podcast details by ID.

- **POST /podcasts/{id}/like**  
  Like a podcast (increments like count).

- **POST /podcasts/{id}/play**  
  Track podcast play count or add it to the user's history.

- **POST /podcasts/{id}/download**  
  Allow users to download a podcast.

- **GET /user/downloads**  
  List all podcasts the user has downloaded.

- **POST /feedback**  
  Allow users to submit feedback for the app or podcasts.

- **GET /trending/podcasts**  
  List podcasts based on most likes, most played, etc.

---

## 4. Bookmarks Module
Handles user bookmarks for podcasts.

- **GET /user/bookmarks**  
  List all podcasts bookmarked by the user.

- **POST /user/bookmarks/{podcast_id}**  
  Toggle podcast remove and add bookmarks.

---

## 5. Notifications Module
Handles notification-related functionality for the user.

- **GET /notifications**  
  List notifications for the user (e.g., new podcast recommendations).

---

# Summary of Endpoints

- **Users Module**: 8 endpoints
- **Categories Module**: 3 endpoints
- **Podcasts Module**: 8 endpoints
- **Bookmarks Module**: 2 endpoints
- **Notifications Module**: 1 endpoint

**Total**: 21 Endpoints
