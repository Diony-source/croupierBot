// cmd/CroupierBot/main.go
package main

import (
	"fmt"

	"github.com/Diony-source/CroupierBot/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Println("Error running the application:", err)
	}
}
