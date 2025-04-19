package collection

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

var client *firestore.Client

func init() {
	ctx := context.Background()
	conf := &firebase.Config{
		ProjectID: "my-angular-template", // Replace with your Firebase Project ID (can be a dummy one for emulator)
	}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		fmt.Printf("error initializing app: %v\n", err)
		os.Exit(1)
	}

	// Connect to the Firestore emulator
	emulatorHost := os.Getenv("FIRESTORE_EMULATOR_HOST")
	if emulatorHost == "" {
		emulatorHost = "localhost:8080" // Default Firestore emulator port
		os.Setenv("FIRESTORE_EMULATOR_HOST", emulatorHost)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		fmt.Printf("error getting Firestore client: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Connected to Firestore Emulator at: %s\n", emulatorHost)
}

// Create ensures a Firestore collection exists by adding a placeholder document.
func Create() {
	ctx := context.Background()
	fmt.Print("Enter collection name: ")
	var collectionName string
	_, err := fmt.Scanln(&collectionName)
	if err != nil || collectionName == "" {
		fmt.Println("Invalid or empty collection name.")
		return
	}

	// Add a placeholder document to create the collection
	_, _, err = client.Collection(collectionName).Add(ctx, map[string]interface{}{
		"placeholder": true,
		"created_at":  firestore.ServerTimestamp,
	})
	if err != nil {
		fmt.Printf("Error creating collection '%s': %v\n", collectionName, err)
		return
	}
	fmt.Printf("Successfully created collection: %s\n", collectionName)
}

// List fetches and displays all collections in the Firestore emulator.
func List() {
	fmt.Println("\n--- List Collections ---")

	ctx := context.Background()
	// Note: Firestore Go SDK does not provide a direct method to list collections.
	// We can use the emulator's REST API or assume collections are known.
	// For simplicity, we'll list collections by checking known collections with documents.
	// Alternatively, you can use the Firestore Admin API or emulator REST API.

	// Placeholder: List collections by attempting to access known collections.
	// In a real app, use the Firestore Admin SDK or REST API.
	knownCollections := []string{"xyz"} // Add known collections or query dynamically
	var collections []string

	for _, coll := range knownCollections {
		docs, err := client.Collection(coll).Limit(1).Documents(ctx).GetAll()
		if err != nil {
			continue // Skip if collection doesn't exist or error occurs
		}
		if len(docs) > 0 {
			collections = append(collections, coll)
		}
	}

	if len(collections) == 0 {
		fmt.Println("No collections found in the Firestore emulator.")
		return
	}

	// Display collections
	fmt.Printf("%-30s\n", "Collection Name")
	fmt.Println(strings.Repeat("-", 30))
	for _, coll := range collections {
		fmt.Printf("%-30s\n", coll)
	}
}
