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

			fmt.Print("Enter JSON file name (e.g., data.json): ")
			fileName, _ := reader.ReadString('\n')
			fileName = strings.TrimSpace(fileName)

			document.Create(collectionName, fileName)
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
	fmt.Println(indent + "Main Menu:") // You might want to make this dynamic for submenus

	for _, item := range items {
		fmt.Printf("%s%d. %s\n", indent, item.ID, item.Label)
	}
	fmt.Printf(indent + "Enter your choice: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	selectedID := 0
	_, err := fmt.Sscan(input, &selectedID)
	if err != nil {
		fmt.Println(indent + "Invalid input. Please enter the number of your choice.")
		displayMenu(items, level, actions, reader) // Recursive call to redisplay
		return
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
				displayMenu(item.SubMenu, level+1, actions, reader) // Recursive call for submenu
			}
			break
		}
	}

	if !found {
		fmt.Println(indent + "Invalid choice. Please try again.")
		displayMenu(items, level, actions, reader) // Recursive call to redisplay
	}
}
