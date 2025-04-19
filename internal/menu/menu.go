package menu

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/swap1210/local-firebase-emulator-initializer/internal/collection"
	"github.com/swap1210/local-firebase-emulator-initializer/internal/document"
	"github.com/swap1210/local-firebase-emulator-initializer/internal/user"
)

// MenuEntry represents a menu item from the JSON file.
type MenuEntry struct {
	ID      int         `json:"id"`
	Label   string      `json:"label"`
	Action  string      `json:"action"`
	SubMenu []MenuEntry `json:"subMenu,omitempty"` // For nested menus
}

// WelcomeScreen displays the initial welcome message.
func WelcomeScreen() {
	fmt.Println("************************************")
	fmt.Println("* Welcome to the Simple Project! *")
	fmt.Println("************************************")
	fmt.Println()
}

// LoadMenuFromJSON reads the menu configuration from a JSON file.
func LoadMenuFromJSON(filename string) []MenuEntry {
	data, err := os.ReadFile("assets/" + filename)
	if err != nil {
		fmt.Printf("Error reading JSON file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	var menu struct {
		Menu []MenuEntry `json:"menu"`
	}
	if err := json.Unmarshal(data, &menu); err != nil {
		fmt.Printf("Error unmarshalling JSON from '%s': %v\n", filename, err)
		os.Exit(1)
	}
	return menu.Menu
}

// MainMenu handles the display and interaction of the main menu.
func MainMenu(menuItems []MenuEntry) {
	reader := bufio.NewReader(os.Stdin)

	actionMap := map[string]func(){
		"user.Create":       user.Create,
		"user.List":         user.List,
		"collection.Create": collection.Create,
		"collection.List":   collection.List,
		"document.Create": func() {
			fmt.Print("Enter collection name: ")
			collectionName, _ := reader.ReadString('\n')
			collectionName = strings.TrimSpace(collectionName)

			fmt.Print("Enter document name (auto-generated if empty): ")
			documentName, _ := reader.ReadString('\n')
			documentName = strings.TrimSpace(documentName)

			fmt.Print("Enter JSON file name (e.g., data.json): ")
			dataFileName, _ := reader.ReadString('\n')
			dataFileName = strings.TrimSpace(dataFileName)

			document.Create(collectionName, documentName, dataFileName)
		},
		"document.List": func() {
			fmt.Print("Enter collection name: ")
			collectionName, _ := reader.ReadString('\n')
			collectionName = strings.TrimSpace(collectionName)
			document.List(collectionName)
		},
		"exit": func() {
			fmt.Println("Exiting the application. Goodbye!")
			os.Exit(0)
		},
	}

	displayMenu(menuItems, 0, actionMap, reader)
}

func displayMenu(items []MenuEntry, level int, actions map[string]func(), reader *bufio.Reader) {
	indent := strings.Repeat("  ", level)

	for {
		// Display menu
		fmt.Println()
		fmt.Println(indent + "Menu:")
		for _, item := range items {
			fmt.Printf("%s%d. %s\n", indent, item.ID, item.Label)
		}
		if level > 0 {
			fmt.Printf("%s0. Go Back\n", indent) // <-- Add a "Go Back" option for submenus
		}
		fmt.Printf(indent + "Enter your choice: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		selectedID := 0
		_, err := fmt.Sscan(input, &selectedID)
		if err != nil {
			fmt.Println(indent + "Invalid input. Please enter the number of your choice.")
			continue
		}

		// Handle 'Go Back'
		if level > 0 && selectedID == 0 {
			return // <-- Just return to go back to parent menu
		}

		found := false
		for _, item := range items {
			if item.ID == selectedID {
				found = true
				if item.Action != "" {
					if action, ok := actions[item.Action]; ok {
						action()
					} else {
						fmt.Printf(indent+"Error: Action '%s' not implemented.\n", item.Action)
					}
				}
				if len(item.SubMenu) > 0 {
					displayMenu(item.SubMenu, level+1, actions, reader) // Enter submenu
				}
				break
			}
		}

		if !found {
			fmt.Println(indent + "Invalid choice. Please try again.")
		}
	}
}
