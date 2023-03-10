package web

import (
	"context"
	"log"
	"net/http"
)

func RunStaticServer(ctx context.Context, addr string) error {
	staticHandler := http.FileServer(http.Dir("./static"))
	mux := http.NewServeMux()
	mux.Handle("/", CORSMiddleware(staticHandler))
	staticServer := http.Server{
		Handler: mux,
	}
	go func() {
		<-ctx.Done()
		log.Println("Closing static http server...")
		if err := staticServer.Close(); err != nil {
			log.Println(err)
		}
	}()
	log.Printf("Static server stated listen at %v\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		return err
	}
	return nil
}

func CORSMiddleware(handler http.Handler) http.Handler {
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(wrapper)
}
