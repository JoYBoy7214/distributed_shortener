package main

import (
	"fmt"
	"log"
	"net/http"

	pb "github.com/JoYBoy7214/distributed_shortener/api/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func submitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if r.Method == "POST" {
		postSubmitHandler(w, r)
	}
}
func postSubmitHandler(w http.ResponseWriter, r *http.Request) {

}

var client pb.ShortenerClient

func main() {

	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error in creating a grpc client", err)
	}
	defer conn.Close()
	client = pb.NewShortenerClient(conn)

	http.HandleFunc("/submit", submitHandler)
	http.ListenAndServe(":8080", nil)
	fmt.Println("http server starts at 8080")

}
