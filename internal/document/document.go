package document

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

// Create adds a document to the specified collection using the content of a JSON file from the assets folder.
func Create(collectionName, documentName, fileName string) {
	if collectionName == "" || fileName == "" {
		fmt.Println("Collection name and file name cannot be empty.")
		return
	}

	// Read JSON file from assets folder
	filePath := filepath.Join("assets", fileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading JSON file '%s': %v\n", filePath, err)
		return
	}

	// Parse JSON content
	var content map[string]interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		fmt.Printf("Error parsing JSON content from '%s': %v\n", filePath, err)
		return
	}

	// Add document to the specified collection with auto-generated ID
	ctx := context.Background()
	_, err = client.Collection(collectionName).Doc(documentName).Set(ctx, content)
	if err != nil {
		fmt.Printf("Error creating document in collection '%s': %v\n", collectionName, err)
		return
	}
	fmt.Printf("Successfully created document with ID: %s in collection: %s\n", documentName, collectionName)
}

// List fetches and displays all documents in the specified collection from the Firestore emulator.
func List(collectionName string) {
	if collectionName == "" {
		collectionName = "testCollection" // Fallback for consistency
	}

	fmt.Printf("\n--- List Documents in %s ---\n", collectionName)

	ctx := context.Background()
	// Fetch all documents in the collection
	docs, err := client.Collection(collectionName).Documents(ctx).GetAll()
	if err != nil {
		fmt.Printf("Error listing documents: %v\n", err)
		return
	}

	if len(docs) == 0 {
		fmt.Println("No documents found in the collection.")
		return
	}

	// Display documents in a formatted table
	fmt.Printf("%-40s %-30s\n", "Document ID", "Content (Preview)")
	fmt.Println(strings.Repeat("-", 70))
	for _, doc := range docs {
		data := doc.Data()
		// Create a short preview of the content (e.g., first key-value pair)
		preview := "N/A"
		for k, v := range data {
			preview = fmt.Sprintf("%s: %v", k, v)
			break // Show only the first key-value pair
		}
		fmt.Printf("%-40s %-30s\n", doc.Ref.ID, preview)
	}
}
