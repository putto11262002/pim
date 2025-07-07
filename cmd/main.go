package main

import (
	"fmt"
	"github/putto11262002/pim"
	"os"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

func main() {

	e := pim.NewEditor()
	renderer := pim.NewRenderer(os.Stdout)

	if err := renderer.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing renderer: %v\n", err)
		os.Exit(1)
	}
	defer renderer.Cleanup()

	if err := renderer.Render(e); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering editor: %v\n", err)
		os.Exit(1)
	}

	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.CtrlC {
			return true, nil // Stop listening for keyboard input
		}

		if err := e.Execute(key); err != nil {
			return false, fmt.Errorf("error executing key: %w", err)
		}

		// Render the editor state after each key press
		if err := renderer.Render(e); err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering editor: %v\n", err)
			os.Exit(1)
		}

		return false, nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}
