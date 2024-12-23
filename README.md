# http-server

bootdev course

---

## Assignment 1.1

I needed to build a simple server that only needs to send `404 not found` as a response.

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

Note: Before I have used `mux.HandleFunc()` and now only `.Handle()` because maybe it was overkill. In Assignment 1.2, `mux.Handle()` is used because `http.FileServer` returns an `http.Handler`, which is directly compatible with `mux.Handle()`.

## Assignment 1.3

For this one, I need to return an image from `/assets`.

I added the logo to a folder called `/assets` and then added this line:
```go
mux.Handle("/assets/", http.FileServer(http.Dir(".")))
```
Personal comment:
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

Personal comment: I didn't understand clearly what I was doing, but with the help of the docs and the chatbot on boot.dev, it was easier to understand. 
