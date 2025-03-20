# Al-Khaimah Golang Backend


### Starting the Server
To start the server, run:
```bash
bash start-alkhaimah.sh
```
logs can be found in the file: alkhaimah.log.

### Stopping the Server
To stop the server, run:
```bash
sudo pkill -f alkhaimah
```

# API Endpoints

## 1. Users Module
Handles user authentication, profile management, and preferences.

- **POST /auth/signup** ✅  
  Create a new user account.

- **POST /auth/login** ✅  
  Authenticate user and return a JWT or session.

- **POST /auth/logout** ✅  
  Logout the user (invalidate the JWT or session).

- **GET /user/profile** ✅  
  Fetch the current user's profile details.

- **PUT /user/profile** ✅  
  Update user's profile details.

- **PATCH /user/profile/password** ✅  
  Allow users to change their password.

- **PATCH /user/profile/preferences** ✅  
  Update a user's preferences (e.g., change categories they are interested in).

- **GET /user/bookmarks** ✅
  List all podcasts bookmarked by the user.

---

## 2. Categories Module
Handles category-related operations.

- **GET /categories** ✅  
  Fetch all categories.

- **GET /trending/categories** ⏳  
  List the most popular categories based on user interactions.

---

## 3. Podcasts Module
Handles podcast-related operations such as listing, liking, playing, downloading podcasts, etc.

- **GET /podcasts** ✅   
  Get all podcasts paginated, will be used for the searching (top left corner in design).

- **GET /podcasts/recommended** ✅ 
  Fetch the latest 10 podcasts for each of the categories the user follows (on main page).

- **GET /podcasts/category/{category_id}** ✅  
  On the main page when user scrolls to the left for the category-podcasts and click on "view All" it will get ALL podcasts for that category.

- **GET /podcasts/{id}** ✅
  Fetch podcast details by ID.

- **POST /podcasts/{id}/like** ✅ 
  Like a podcast (increments like count).

- **POST /podcasts/{id}/play** ⏳  
  Track podcast play count or add it to the user's history.

- **POST /podcasts/{id}/download** ⏳  
  Allow users to download a podcast.

- **GET /user/downloads** ⏳  
  List all podcasts the user has downloaded.

- **GET /trending/podcasts** ⏳  
  List podcasts based on most likes, most played, etc.

- **GET /user/history** ⏳  
  List all podcasts watched by the user.

- **POST /user/bookmarks/{podcast_id}** ⏳  
  Toggle podcast remove and add bookmarks.

---

## 4. Notifications Module
Handles notification-related functionality for the user.

- **GET /notifications** ⏳  
  List notifications for the user (e.g., new podcast recommendations).
- //ToDo::

---

## 5. Admin Module
Handles admin-related functionalities.

- **POST /admin/categories** ✅  
  Create a new category.

- **PUT /admin/categories/{id}** ✅  
  Edit an existing category.

- **DELETE /admin/categories/{id}** ✅  
  Delete a category.

- **DELETE /admin/users/{id}** ✅  
  Delete a user account.

- **GET /admin/users** ✅  
  Fetch all users.

---
# Summary of Endpoints

- **Users Module**: 8 endpoints
- **Categories Module**: 3 endpoints
- **Podcasts Module**: 11 endpoints
- **Notifications Module**: 1 endpoint
- **Admin Module**: 5 endpoints

**Total**: **28 Endpoints**