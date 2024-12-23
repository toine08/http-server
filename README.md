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

