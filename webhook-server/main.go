package main

import (
	"encoding/json"
	"log"
	"net/http"
)

/*
A typical kubernetes request comes in as

	{
		apiVersion: "admission.k8s.io/v1",
		kind: "AdmissionReview",
		request: {
			uid: 123,
			allowed: true
		}
	}

And the response we need to send back is of the form

	{
		apiVersion: "admission.k8s.io/v1",
		kind: "AdmissionReview",
		response: {
			uid: 123,
			allowed: true,
		}
	}

This is the reason why we have this struct Admission Review which has the following
*/
type AdmissionReview struct {
	APIVersion string    `json:"apiVersion"`
	Kind       string    `json:"kind"`
	Request    *Request  `json:"request,omitempty"`
	Response   *Response `json:"response,omitempty"`
}

type Request struct {
	UID string `json:"uid"`
}

type Response struct {
	UID     string `json:"uid"`
	Allowed bool   `json:"allowed"`
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("recieved admission request")
	var review AdmissionReview
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		log.Printf("Failed to decode request +%v", err.Error())
		return
	}
	log.Println(review.Request.UID)

	response := AdmissionReview{
		APIVersion: "admission.k8s.io/v1",
		Kind:       "AdmissionReview",
		Response: &Response{
			UID:     review.Request.UID,
			Allowed: true,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	writeErr := json.NewEncoder(w).Encode(response)
	if writeErr != nil {
		log.Printf("Failed to write response +%v", writeErr.Error())
		return
	}
	log.Println("admission response sent")
}
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func main() {
	http.HandleFunc("/validate", validateHandler)
	http.HandleFunc("/healthz", healthHandler)
	log.Println("starting webhook server on :8443")
	err := http.ListenAndServeTLS(":8443", "certs/tls.crt", "certs/tls.key", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
