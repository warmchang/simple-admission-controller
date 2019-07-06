package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	port = ":8080"
)

var (
	tlscert, tlskey string
)

func main() {
	flag.StringVar(&tlscert, "tlsCertFile", "/etc/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&tlskey, "tlsKeyFile", "/etc/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")
	flag.Parse()

	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		var body []byte
		if r.Body != nil {
			if data, err := ioutil.ReadAll(r.Body); err == nil {
				body = data
			}
		}

		if len(body) == 0 {
			log.Println("no body recieved in request")
			http.Error(w, "empty body", http.StatusBadRequest)
			return
		}

		arRequest := v1beta1.AdmissionReview{}
		if err := json.Unmarshal(body, &arRequest); err != nil {
			log.Println("invalid request type recieved")
			http.Error(w, "incorrect body", http.StatusBadRequest)
		}

		raw := arRequest.Request.Object.Raw
		pod := v1.Pod{}
		if err := json.Unmarshal(raw, &pod); err != nil {
			log.Println("error unmarshalling pod")
			return
		}

		if pod.Namespace != "default" {
			log.Printf("recieved pod event in namespace %s, allowing", pod.Namespace)
			return
		}

		arResponse := v1beta1.AdmissionReview{
			Response: &v1beta1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Message: "Please use another namespace for pod, the default namespace is restricted!",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(arResponse)
	})

	log.Fatal(http.ListenAndServeTLS(port, tlscert, tlskey, nil))
}
