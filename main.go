package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/atotto/clipboard"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

func main() {
	bucket := flag.String("b", "simonso-access-control", "Bucket name")
	path := flag.String("p", "sg-7d482705", "Path to object")
	region := flag.String("r", "us-west-2", "AWS region")
	expiry := flag.String("x", "120s", "Expiry")

	// Build the S3 request URL
	URL := "https://s3-" + *region + ".amazonaws.com/" + *bucket + "/" + *path
	fmt.Printf("S3 URL: %v\n", URL)

	// Grab AWS credentials from environment variables
	creds := credentials.NewEnvCredentials()

	// Create a new v4 signer
	s := v4.NewSigner(creds)

	// Create a new HTTP request object
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatalf("Error creating HTTP request\n")
	}

	// Parse the expiry or duration string (time that the signed URL will be valid for)
	d, err := time.ParseDuration(*expiry)
	if err != nil {
		log.Fatalf("Failed to parse duration %s\n", d)
	}

	// Create the signed URL
	// s.Presign() returns a header but we don't care
	// The signed URL will be in req.URL
	_, err = s.Presign(req, nil, "s3", *region, d, time.Now())
	if err != nil {
		log.Fatalf("Error signing request\n")
	}

	// Print signed URL
	fmt.Printf("Signed URL: %+v\n", req.URL)

	// Copy URL to clipboard
	err = clipboard.WriteAll(req.URL.String())
	if err != nil {
		log.Printf("Error copying URL to clipboard: %+v\n", err)
	}
}
