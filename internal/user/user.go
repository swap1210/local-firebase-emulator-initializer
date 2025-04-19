package user

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

var authClient *auth.Client

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

	// Ensure the Auth emulator is used via environment variable
	emulatorHost := os.Getenv("FIREBASE_AUTH_EMULATOR_HOST")
	if emulatorHost == "" {
		emulatorHost = "localhost:9099" // Default Auth emulator port
		os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", emulatorHost)
	}

	authClient, err = app.Auth(ctx)
	if err != nil {
		fmt.Printf("error getting Auth client: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Connected to Firebase Auth Emulator at: %s\n", emulatorHost)
}

// Create handles the user creation functionality by reading from a JSON file.
func Create() {
	// Read users from JSON file
	data, err := os.ReadFile("assets/users.json")
	if err != nil {
		fmt.Printf("error reading users JSON file: %v", err)
	}

	var users struct {
		Users []struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"users"`
	}
	if err := json.Unmarshal(data, &users); err != nil {
		fmt.Printf("error unmarshalling users JSON: %v", err)
	}

	if len(users.Users) == 0 {
		fmt.Printf("no users found in JSON file")
	}

	// Process each user in a loop
	var firstError error
	for _, user := range users.Users {
		if user.Username == "" || user.Password == "" {
			fmt.Printf("Skipping user with empty username or password\n")
			if firstError == nil {
				firstError = fmt.Errorf("one or more users had empty username or password")
			}
			continue
		}

		params := (&auth.UserToCreate{}).
			Email(user.Username).
			Password(user.Password)

		ctx := context.Background()
		u, err := authClient.CreateUser(ctx, params)
		if err != nil {
			fmt.Printf("Error creating user %s: %v\n", user.Username, err)
			if firstError == nil {
				firstError = fmt.Errorf("error creating user %s: %v", user.Username, err)
			}
			continue
		}
		fmt.Printf("Successfully created user: %v\n", u.UID)
	}
}

// List handles the user listing functionality (example for future expansion).
func List() {
	fmt.Println("\n--- List Users ---")

	// Read users from JSON file to get email identifiers
	data, err := os.ReadFile("assets/users.json")
	if err != nil {
		fmt.Printf("Error reading users JSON file: %v\n", err)
		return
	}

	var users struct {
		Users []struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"users"`
	}
	if err := json.Unmarshal(data, &users); err != nil {
		fmt.Printf("Error unmarshalling users JSON: %v\n", err)
		return
	}

	if len(users.Users) == 0 {
		fmt.Println("No users found in JSON file.")
		return
	}

	// Prepare identifiers for GetUsers
	var identifiers []auth.UserIdentifier
	for _, user := range users.Users {
		if user.Username == "" {
			continue
		}
		identifiers = append(identifiers, auth.EmailIdentifier{Email: user.Username})
	}

	if len(identifiers) == 0 {
		fmt.Println("No valid user identifiers found in JSON file.")
		return
	}

	// Fetch users using GetUsers
	ctx := context.Background()
	result, err := authClient.GetUsers(ctx, identifiers)
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
		return
	}

	// Collect user data
	var userList []struct {
		UID   string
		Email string
	}
	for _, u := range result.Users {
		userList = append(userList, struct {
			UID   string
			Email string
		}{
			UID:   u.UID,
			Email: u.Email,
		})
	}

	if len(userList) == 0 {
		fmt.Println("No users found in the Firebase Auth emulator for the provided identifiers.")
		return
	}

	// Display users in a formatted table
	fmt.Printf("%-30s %-40s\n", "Email", "UID")
	fmt.Println(strings.Repeat("-", 70))
	for _, u := range userList {
		fmt.Printf("%-30s %-40s\n", u.Email, u.UID)
	}

	// Report users not found (if any)
	for _, nf := range result.NotFound {
		switch id := nf.(type) {
		case auth.EmailIdentifier:
			fmt.Printf("User not found: %s\n", id.Email)
		}
	}
}
