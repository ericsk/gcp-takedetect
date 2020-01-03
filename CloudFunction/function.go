package p

import (
	"context"
	"log"

	storage "cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
)

// GCSEvent is the payload of a Google Cloud Storage event.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// DetectGCSImage prints a message when a file is changed in a Cloud Storage bucket.
func DetectGCSImage(ctx context.Context, e GCSEvent) error {
	log.Printf("Processing file: %s; bucket: %s", e.Name, e.Bucket)

	// Creates a Storage client.
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	defer storageClient.Close()

	// Create the object reader
	rc, err := storageClient.Bucket(e.Bucket).Object(e.Name).NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to create the object reader: %v", err)
	}
	defer rc.Close()

	// Read object content and convert it to vision.Image type
	data, err := vision.NewImageFromReader(rc)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	// Creates a Cloud Vision API client.
	visionClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer visionClient.Close()

	// Using Cloud Vision API to detect labels in the image
	labels, err := visionClient.DetectLabels(ctx, data, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}

	log.Println("Labels:")
	for _, label := range labels {
		log.Println(label.Description)
	}

	return nil
}
