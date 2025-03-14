# Rootext Company Tack


## Task Discription

🔹 Task Objective: Implement a simple API to manage posts and user ratings (similar to the Upvote/Downvote system on Reddit). This project includes CRUD operations, sorting by score, and caching with Redis.
⏳ Estimated time: About 10 hours (can be divided into several days) until Friday night at the latest.

1️⃣ User registration and login (with JWT)  
    - Users should be able to register, log in, and log out.  
    - After logging in, they receive a JWT token.  

2️⃣ Create, delete, and retrieve posts  
    - Each user can create a post (including title and text).  
    - Users should be able to manage their own posts (edit, delete).  

3️⃣ Rating and interacting with posts  
    - Users can give positive (+1) or negative (-1) ratings to posts.  
    - The final score of each post = total votes.  
    - Each user can only give one vote to a post (not twice).  

4️⃣ Sorting posts by score and date  
    - Users can see the top posts of the day, week, and month.  

5️⃣ Caching the list of popular posts in Redis  
    - The list of posts with scores (e.g., top 5) should be stored in Redis to improve performance.  

6️⃣ API documentation with Swagger or Postman Collection  
    - APIs should be fully documented on how they work.  


### How run it ?

First, install the Makefile and run `make prod`. This command will run Docker Compose. Then, import the Postman JSON file from the `doc` folder into Postman and send requests to the API.
