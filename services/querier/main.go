package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davidswisa/multiple-containers-in-pod/pkg/orm"
	"github.com/rs/cors"
)

var (
	index = 0
)

func producerHandler() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {

			// os.Setenv("API_HOST", "20.10.1.8")
			client := orm.NewORMClient()

			b, res, err := client.Get("reservations", http.Header{})
			if err != nil {
				return
			}
			fmt.Println(b)
			fmt.Println(res)
			if b == "null" {
				b = "{}"
			}

			wrt.Write([]byte(b))
		}
	})
}

func main() {

	// Add handle func for producer.
	mux := http.NewServeMux()

	mux.HandleFunc("/reservations", producerHandler())

	// Run the web server.
	handler := cors.Default().Handler(mux)
	fmt.Println("starting producer-api...")
	log.Fatal(http.ListenAndServe(":8081", handler))
}
