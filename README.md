# http-server
bootdev course


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

## Assignement 1.2

