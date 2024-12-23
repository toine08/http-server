package main

import(
	"net/http"
)

func main(){
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:":8080",
		Handler:mux,
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request){
		http.NotFound(w,req)
			
	})
	server.ListenAndServe()
}
