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