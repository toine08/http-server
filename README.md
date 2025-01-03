# http-server

boot.dev course

---

## Assignment 1.1

I needed to build a simple server that only sends `404 not found` as a response.

This server can handle multiple requests concurrently, utilizing the full capacity of the CPU. For example, Python with Django or Flask does not natively support multi-threading for handling requests.

Here is the main function:

```go
func main() {
    // Initialize the multi-threaded server
	mux := http.NewServeMux()
    // Create the server and specify the address and handler. Here it's port 8080 and the handler is the mux variable
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
    // Add a handler function to respond with a 404 error for all requests
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.NotFound(w, req)
	})

    // Start the server and listen on the specified port
	server.ListenAndServe()
}
```

## Assignment 1.2

For this assignment, I needed to serve an HTML file called `index.html` when it's requested at the root.

Here is the main function:

```go
func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
    // Use Handle to serve the HTML file at the root path
	mux.Handle("/", http.FileServer(http.Dir(".")))
	server.ListenAndServe()
}
```

#### Note: 
Before, I used `mux.HandleFunc()` and now only `.Handle()` because maybe it was overkill. In Assignment 1.2, `mux.Handle()` is used because `http.FileServer` returns an `http.Handler`, which is directly compatible with `mux.Handle()`.

## Assignment 1.3

For this one, I needed to return an image from `/assets`.

I added the logo to a folder called `/assets` and then added this line:
```go
mux.Handle("/assets/", http.FileServer(http.Dir(".")))
```
#### Note:
I didn't understand at first. I tried to use `http.FileServer()` without moving `logo.png` to the folder but didn't find any solution. Maybe in the future.

## Assignment 1.4

For this one, I had to add a readiness endpoint accessible from `/healthz` to check if our server is ready to receive some requests. I also had to update the file server path to avoid potential conflict with the file server handler. Now instead of using `/`, we go through `/app/`.

Here is the code for this version:
```go
func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
    // Serve files from the current directory under the /app/ path, stripping the /app/ prefix
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
    // Serve the logo.png file at the /assets path
	mux.Handle("/assets/", http.FileServer(http.Dir(".")))
    // Add a readiness endpoint at /healthz to check if the server is ready to receive requests
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	server.ListenAndServe()
}
```

#### Note: 
I didn't understand clearly what I was doing, but with the help of the docs and the chatbot on boot.dev, it was easier to understand. 


## Assignment 2.1

Well, this one was hard. A lot of things I never used and it feels a bit like walking in fog and you can't see 1 meter in front of you. I tried to not use AI but in the end, I was lost and didn't clearly understand some things I had to do. 

### Here is the assignment (copied from the website because it's a long one):
---
Create a struct in main.go that will hold any stateful, in-memory data we'll need to keep track of. In our case, we just need to keep track of the number of requests we've received.

```go
type apiConfig struct {
	fileserverHits atomic.Int32
}
```

The atomic.Int32 type is a really cool standard-library type that allows us to safely increment and read an integer value across multiple goroutines (HTTP requests).

Next, write a new middleware method on a *apiConfig that increments the fileserverHits counter every time it's called. Here's the method signature I used:
```go
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	// ...
}
```

The atomic.Int32 type has an .Add() method, use it to safely increment the number of fileserverHits.

Wrap the http.FileServer handler with the middleware method we just wrote. For example:
```go
mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
```

Create a new handler that writes the number of requests that have been counted as plain text in this format to the HTTP response:

```
Hits: x
```

Where x is the number of requests that have been processed. This handler should be a method on the *apiConfig struct so that it can access the fileserverHits data.

Register that handler with the serve mux on the /metrics path.
Finally, create and register a handler on the /reset path that, when hit, will reset your fileserverHits back to 0.

It should follow the same design as the previous handlers.

Remember, similar to the metrics endpoint, /reset will need to be a method on the *apiConfig struct so that it can also access the fileserverHits

---

Here is the code after finishing the assignment:

```go
// Create a new structure to hold the atomic.Int32
type apiConfig struct {
	fileserverHits atomic.Int32
}

// This function is the middleware, it's called every time the /app/ is reached
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

// This function is used to print the value of the fileserverHits when it's called on the /metrics page
func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, req *http.Request) {
	hits := cfg.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
}

// This function resets the hits to 0, it's called on the /reset page
func (cfg *apiConfig) handleReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits reset to 0")))
}

func main() {
	// ...existing code...

	// Instantiate the config here
	cfg := &apiConfig{}
	// Update the handle for /app/ to add the middleware
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	// ...existing code...

	// Add a readiness endpoint at /healthz to check if the server is ready to receive requests
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	// Here I call the handleMetrics to add 1 to the fileServerHits variable and the reset.
	mux.HandleFunc("/metrics", cfg.handleMetrics)
	mux.HandleFunc("/reset", cfg.handleReset)
	// Start the server and listen on the specified port
	server.ListenAndServe()
}
```
#### Note:
 This assignment was quite challenging. It was rated 9/10 in difficulty on the course, and I felt a bit lost at times. I tried to avoid using AI assistance, but eventually, I needed help to understand some concepts. I am starting to grasp why certain functions are created, but it still feels somewhat unclear.

## Assignment 2.2

 For this one, I needed to specify which method to use for which route to avoid any issues.

 This was easy in Golang because you can specify the method, the route, and even the port. So, I only had to add GET/POST to the different routes where it was required.

 ```go
 func main() {
	//...existing code...
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	//added the method GET for /metrics and /healthz
	mux.HandleFunc("GET /metrics", cfg.handleMetrics)
	//added the method POST for the reset route
	mux.HandleFunc("POST /reset", cfg.handleReset)
	//...existing code
}
 ```

#### Note:
This was not really hard and clear to understand. At first, I thought I needed to handle the case where the method was not the correct one to return the desired HTTP code.

## Assignment 2.3

This one was also an easy one. I had to add the /api/ route before healthz, metrics, and reset.

 ```go
  func main() {
	//...existing code...
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	//added the method GET for /metrics and /healthz
	mux.HandleFunc("GET /api/metrics", cfg.handleMetrics)
	//added the method POST for the reset route
	mux.HandleFunc("POST /api/reset", cfg.handleReset)
	//...existing code
}
 ```

#### Note:

 This one was really not hard to understand. Nice to have an easy one haha.

## Assignment 3.4

This one was a bit difficult. I didn't know how to use the HTML template and thought maybe I needed to create a new HTML file. I had to use AI to inform myself.

### Assignment:

Swap out the GET /api/metrics endpoint, which just returns plain text, for a GET /admin/metrics that returns HTML to be rendered in the browser.

Update the POST /api/reset to POST /admin/reset. Its functionality should not change.

```go
//...New version of function handleMetrics:
func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, req *http.Request) {

	hits := cfg.fileserverHits.Load()
	html := fmt.Sprintf(`<html>
  <body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func main(){
	//...existing code...
	mux.HandleFunc("GET /admin/metrics", cfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handleReset)
	//...existing code...
}
```

#### Note:
I thought it would be hard, but it wasn't the case. It was just difficult to know how to integrate the HTML in the function. Sadly, I had to use AI...

## Assignment 4.2

Well, this one was not an easy one but definitely interesting. I had to use the AI chat from boot.dev to get some help.

### Assignment:
Add a new endpoint to the Chirpy API that accepts a POST request at /api/validate_chirp. It should expect a JSON body of this shape:

```go

//I never know how to name my function...
func handleValidation(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return

	}
	if len(params.Body) > 140 {
		w.WriteHeader(400)
		w.Write([]byte(`{"error": "Chirp is too long"}`))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(`{"valid":true}`))

}


//add this line in the main function
mux.HandleFunc("POST /api/validate_chirp", handleValidation)

```
#### Note:

This one I thought it was going to be harder than that. Of course, I used AI but I was more using it like a counselor than a teacher. I don't like writing some code if it's useless so I ask for its opinion before writing my code. But it also helped me with some debugging because at first, I thought the function was not working correctly but I was just using the wrong tab on Thunder Client...

## Assignment 4.6

This one was nice, I was able to resolve it almost alone. Just needed to google 2-3 things to find out how to loop over a string.

### Assignment:
We need to update the /api/validate_chirp endpoint to replace all "profane" words with 4 asterisks: `****`.

Assuming the length validation passed, replace any of the following words in the Chirp with the static 4-character string `****`:

	kerfuffle
	sharbert
	fornax

Be sure to match against uppercase versions of the words as well, but not punctuation. "Sharbert!" does not need to be replaced, we'll consider it a different word due to the exclamation point. Finally, instead of the valid boolean, your handler should return the cleaned version of the text in a JSON response.

```go
func handleValidation(w http.ResponseWriter, req *http.Request) {
//...existing code...
validateString := params.Body
	words := strings.Split(validateString, " ")
	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			words[i] = "****"
		}

	}
	validateString = strings.Join(words, " ")
	returnedString := fmt.Sprintf(`{"cleaned_body": "%s"}`, validateString)
	w.WriteHeader(200)
	w.Write([]byte(returnedString))
//...existing code...
}
```
#### Note:
Good to feel confident on an assignment haha.

## Assignment 5.3

Wow this one was a lot installation and setting things up. I had trouble to install postgresql because I wasn't paying attention to the logs from the brew install (like exporting to the .zshrc the config to use the command line). 

### Assignment:
Well I had to install postgresql, goose (actually it was the precedent assignment), sqlc and add some sql file and folder and finally add some imports and code in the main function and the type apiConfig struct

```go
//... existing code
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/toine08/http-server/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}
//...existing code...
func main() {
	//load env key
	godotenv.Load()
	cfg := &apiConfig{}
	//get the wanted env variable
	dbURL := os.Getenv("DB_URL")
	//open the connection to the db with the env variable
	db, err := sql.Open("postgres", dbURL)
	//stop server if there is an error
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	cfg.dbQueries = database.New(db)
	//...existing code...
}

```

#### Note:

Well this was a lot of installation, at some point I was lost in the indication and things I have to do. But I did it, I had to use AI but I had trouble installing postgresql. 

## Assignment 5.5

Check Assignment 5.6...

## Assignment 5.6

Well, this one was hard. Too hard, I had to look at the answers. But in this bad news, I discovered that I can work in other files...
I mean, I know it's possible, but I didn't know this for Go. I haven't written anything for `Assignment 5.5` because it was about creating a user, and I wasn't able to resolve 5.6, certainly because my function wasn't correct, sadly.

Also, I won't show the code for this one because it would be too long. Sorry.

## Assignment 5.9

Ok, back to coding this time (not that I haven't coded for 5.5 or 5.6, but I didn't count since I had to use the provided answers...)

### Assignment:
Add a new query that retrieves all chirps in ascending order by created_at.
Add a GET /api/chirps endpoint that returns all chirps in the database. It should return them in the same structure as the POST /api/chirps endpoint, but as an array. Use a 200 status code for success. Order them by created_at in ascending order.

```sql
-- name: AllChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;
```

```go
//file handle_chirps_allchirps.go
package main

import "net/http"

func (cfg *apiConfig) handleAllChirps(w http.ResponseWriter, r *http.Request) {
	var chirps []Chirp
	rows, err := cfg.dbQueries.AllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error retrieving data", err)
		return
	}
	for _, row := range rows {
		var chirp Chirp
		chirp.ID = row.ID
		chirp.CreatedAt = row.CreatedAt
		chirp.UpdatedAt = row.UpdatedAt
		chirp.Body = row.Body
		chirp.UserID = row.UserID
		chirps = append(chirps, chirp)
	}

	respondWithJSON(w, 200, chirps)
}

//file main.go

//existing code...
func main(){
	//existing code...
	mux.HandleFunc("GET /api/chirps", cfg.handleAllChirps)
	//existing code...
}
```

#### Note:
With this code, I was able to retrieve all the chirps created in the DB. Which is nice, hehe. I had to use AI to help me because the course doesn't explain how to manage data received from the DB. Maybe some hints would have been nice. But once I figured that out, it was pretty easy.

## Assignment 5.10

A lot of things to do in this chapter !

### Assignment:
Basically adding the route to return a chirp by an ID

```go
//handle_chirps_chirpsbyid.go
func (cfg *apiConfig) handleChirpsById(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")
	uuid, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, 404, "Error while getting the value", err)
		return
	}
	row, err := cfg.dbQueries.ChirpByID(r.Context(), uuid)
	if err != nil {
		respondWithError(w, 500, "Error there is no ID who match the chirp", err)
	}
	respondWithJSON(w, 200, Chirp{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		Body:      row.Body,
		UserID:    row.UserID,
	})

}


//main.go
//main function as always:
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handleChirpsById)

```
```sql
-- /chirp.sql
-- name: ChirpByID :one
SELECT id, created_at, updated_at, body, user_id FROM chirps WHERE id = $1;
```

#### Note:
Well this one was nice to do. I had to use AI to get some info about the function `http.Request.PathValue` and how to use it because the doc wasn't giving any example. I also use the AI to write the query for myself...(Probably the next course will be SQL haha).

It feels great to code and to feel better at this. 


## Assignment 6.1

Soooo, this one was hard. At first, I thought it would be easy to do, but after a few tries, I changed my mind. I got some issues with the `/api/login`, but the issue was coming from `handle_users_create.go`...

### Assignment:
Update the POST /api/users endpoint. The body parameters should now require a new password field. Add a POST /api/login endpoint. This endpoint should allow a user to log in. In a future exercise, this endpoint will be used to give the user a token that they can use to make authenticated requests. For now, let's just make sure password validation is working. It should accept this body:
```json
{
	"password": "04234",
	"email": "lane@example.com"
}
```
I have updated the table users, see `003_add_hashed_password.sql`. I have created the query to getUserByEmail:
```sql
-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password FROM users WHERE email =$1;
```

I have created two files in `internal/auth`:
- hashPassword.go
- checkHashPassword.go

```go
//both files in here:
func HashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if (err != nil) {
		return password, errors.New("Error while hashing password")
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
```

I have updated the `handle_user_create.go`:
```go
//had to update the User struct...
type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	//...existing code
	hashedPassword, err := auth.HashPassword(params.Password)
	
	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	//...existing code
}
```
```go
func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		fmt.Printf("Error finding user: %v\n", err) // Debug print
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		fmt.Printf("Error comparing passwords: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
```

#### Note:
This one was hard. I had to use AI to get some help for debugging; otherwise, I would have gone crazy before finding the error...

## Assignment 6.6

The assignment was quite extensive, so here is a summary:

Add JWT functionalities and update login and chirps to work with JWT.

Here are the files created:

- Added files to the auth package:
	- `makeJWT.go`
	- `validateJWT.go`
	- `getBearerToken.go`
- Added `handle_users_authentication.go`
- Added a new function to the queries to get user by ID.

#### Note:
This one was really challenging. There were a lot of things to do, and I felt lost at times. I am not happy with my performance. Even though I knew where to add certain things, I didn't know how to use some functions. I had to rely heavily on AI. Maybe going all the way through the HTTP server wasn't the best idea...



## Assignment 6.11
This one was also a big one, here is the assignment and what I did.

### Assignment:
- Create a new database table called `refresh_tokens`.
- Add a `func MakeRefreshToken() (string, error)` function to your `internal/auth` package. It should use the following to generate a random 256-bit (32-byte) hex-encoded string.
- Update the `POST /api/login` endpoint to return a refresh token, as well as an access token:
	```json
	{
		"id": "5a47789c-a617-444a-8a80-b50359247804",
		"created_at": "2021-07-01T00:00:00Z",
		"updated_at": "2021-07-01T00:00:00Z",
		"email": "lane@example.com",
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		"refresh_token": "56aa826d22baab4b5ec2cea41a59ecbba03e542aedbb31d9b80326ac8ffcfa2a"
	}
	```
- Create a `POST /api/refresh` endpoint. This new endpoint does not accept a request body but requires a refresh token to be present in the headers, in the same `Authorization: Bearer <token>` format.
- Create a new `POST /api/revoke` endpoint. This new endpoint does not accept a request body but requires a refresh token to be present in the headers, in the same `Authorization: Bearer <token>` format.

For the code, sadly there are too many things I have updated, but here is what I have updated and created:

- Created new schema `004_add_refresh_tokens.sql` and `refresh_token.sql` for the queries.
- Created new endpoint `handle_users_refresh.go` and `handle_users_revoke.go` and updated the `main.go` to add the routes.
- Created new function `makeRefreshToken.go`.

#### Note:
This one was hard, but clearer than the previous one. I still had to use AI, but this was to insert debug and help to know how to use a function. I really enjoy coding in Go, but I know I have some improvements to make.

## Assignment 7.1

### Assignment:
Add a PUT /api/users endpoint so that users can update their own (but not others') email and password. It requires:

Here is my code:
```go
func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type userResponse struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No valid bearer token provided", err)
		return
	}

	claims, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)

	newUser, err := cfg.dbQueries.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:             claims,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error while updating the password", err)
		return
	}

	response := userResponse{
		ID:        newUser.ID,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	respondWithJSON(w, http.StatusOK, response)

}
```

```sql
-- name: UpdateUser :one
UPDATE users
SET email = $2, hashed_password = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;
```

I have also updated main.go to add the route

#### Note:
This one was moderately challenging. I've learned that I should start by writing down my tasks before diving into coding. When I first start, I usually have good ideas and feel inspired, but as I work on the implementation, I sometimes lose direction. Despite this, it was a satisfying chapter to complete.

## Assignment 7.4

I'm pleased with how this turned out!

### Assignment:
- Add a new DELETE /api/chirps/{chirpID} route to your server that deletes a chirp from the database by its id.

Here is my code:
```go
func (cfg *apiConfig) handleDeleteChirpsById(w http.ResponseWriter, req *http.Request) {
	//verify identity
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No valid bearer token provided", err)
		return
	}

	claims, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token", err)
		return
	}
	//get the ChirpID
	chirpId := req.PathValue("chirpID")
	uuid, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, 404, "Error while getting the value", err)
		return
	}
	if _, err := cfg.dbQueries.DeleteChirpByID(req.Context(), database.DeleteChirpByIDParams{
		ID:     uuid,
		UserID: claims,
	}); err != nil {
		respondWithError(w, http.StatusForbidden, "Error while deleting the chirp", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
```
```sql
-- name: DeleteChirpByID :one
DELETE FROM chirps WHERE id = $1 AND user_id = $2
RETURNING *;
```

#### Note:

Well this one was also not that hard but not easy. I have avoided to use AI but It's pretty useful for writing correct sql queries hehe. I am pretty proud of what I do and I haven't felt lost so that's cool. 
