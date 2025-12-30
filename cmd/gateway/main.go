package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	pb "github.com/JoYBoy7214/distributed_shortener/api/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type ShortUrlRequest struct {
	OriginalUrl string `json:"OriginalUrl"`
	UserId      string `json:"UserId"`
}
type ShortUrlResponse struct {
	ShortUrl string `json:"ShortUrl"`
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		postSubmitHandler(w, r)
	}
}
func postSubmitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var msg ShortUrlRequest
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "error in decoding the incoming request", http.StatusBadRequest)
	}
	shortRequest := pb.CreateShortUrlRequest{
		OriginalUrl: msg.OriginalUrl,
		UserId:      msg.UserId,
	}
	res, err := client.CreateShortUrl(r.Context(), &shortRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var response ShortUrlResponse
	response.ShortUrl = res.ShortUrl
	json.NewEncoder(w).Encode(response)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {

	shortUrl := r.PathValue("shortId")
	//fmt.Println(shortUrl)
	res, err := client.GetOriginalUrl(r.Context(), &pb.GetOriginalUrlRequest{
		ShortUrl: shortUrl,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Println(res.OriginalUrl)
	http.Redirect(w, r, res.OriginalUrl, http.StatusSeeOther)

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
	http.HandleFunc("/{shortId}", redirectHandler)
	fmt.Println("http server starts at 8080")
	http.ListenAndServe(":8080", nil)

}
