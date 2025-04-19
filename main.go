package main

import (
	"github.com/swap1210/local-firebase-emulator-initializer/internal/menu"
)

func main() {
	menu.WelcomeScreen()
	menuItems := menu.LoadMenuFromJSON("assets/menu.json")
	menu.MainMenu(menuItems)
}
