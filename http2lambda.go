package main

import (
	"flag"
    "fmt"
	"encoding/json"
	"errors"
	"log"
    "net/http"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/gen/lambda"
)

func main() {
	port := flag.Int("p", 80, "port")

	flag.Parse()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func init() {
    http.HandleFunc("/invoke", handler)
}

func getParam(w http.ResponseWriter, r *http.Request, p string) (string, error) {
	v := r.URL.Query().Get(p)

	if len(v) != 0 {
		return v, nil
	} else {
		http.Error(w, fmt.Sprintf("required parameter '%s' is missing", p), http.StatusBadRequest)
		log.Printf("ERROR: required parameter '%s' is missing", p)
		return "", errors.New("required parameter missing")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: request %s", r.URL.Query())

	accessKey, err := getParam(w, r, "access-key")

	if err != nil {
		return
	}

	secretKey, err := getParam(w, r, "secret-key")

	if err != nil {
		return
	}

	token := r.URL.Query().Get("access-token")

	creds := aws.DetectCreds(accessKey, secretKey, token)

	region, err := getParam(w, r, "region")

	if err != nil {
		return
	}

	cli := lambda.New(creds, region, nil)

	fName, err := getParam(w, r, "function")

	if err != nil {
		return
	}

	q := r.URL.Query()

	delete(q, "access-key")
	delete(q, "secret-key")
	delete(q, "region")
	delete(q, "access-token")
	delete(q, "function")

	log.Printf("INFO: going to invoke %s in %s", fName, region)

	json, err := json.Marshal(q)

	if err != nil {
		http.Error(w, fmt.Sprintf("json marshalling error: %s", err), http.StatusInternalServerError)
		log.Printf("ERROR: json marshalling error %s", err)
		return
	}

	invoke := &lambda.InvokeAsyncRequest{
		FunctionName: &fName,
		InvokeArgs: json,
	}

	resp, err := cli.InvokeAsync(invoke)

	if err != nil {
		http.Error(w, fmt.Sprintf("invocation error on %s: %s", fName, err), http.StatusInternalServerError)
		log.Printf("invocation error on %s: %s", fName, err)
		return
	}

	w.WriteHeader(*resp.Status)
}
