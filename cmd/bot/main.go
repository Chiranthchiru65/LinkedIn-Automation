package main

import (
	"fmt"
	"log"
	"time"

	"linkedin-automation/internal/auth"
	"linkedin-automation/internal/config"
	"linkedin-automation/internal/core"
	outreach "linkedin-automation/internal/outreact"
	"linkedin-automation/internal/search"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Println(" Config Loaded")

	linkedinBot, err := core.NewBrowser(cfg.Settings.Headless)
	if err != nil {
		log.Fatalf("Failed to initialize browser: %v", err)
	}

	fmt.Println(" Bot: initiating login sequence...")

	err = auth.Login(linkedinBot.Page, cfg.Credentials.Username, cfg.Credentials.Password)
	if err != nil {
		log.Fatalf(" Login failed: %v", err)
	}

	fmt.Println(" Core: Initializing Search Engine...")
	engine := search.NewEngine(linkedinBot.Page, "Software Engineer", "Bangalore")

	leads, err := engine.Run()
	if err != nil {
		log.Fatalf(" Search failed: %v", err)
	}

	fmt.Println("------------------------------------------------")
	fmt.Printf(" HARVEST COMPLETE! Found %d leads:\n", len(leads))
	for i, lead := range leads {
		if i >= 2 {
			break
		}

		fmt.Printf("\n--- Processing Lead %d/%d: %s ---\n", i+1, len(leads), lead.Name)

		message := fmt.Sprintf("Hi %s, I found your profile while searching for Software Engineers in Bangalore. I am building a Golang automation tool and would love to connect!", lead.Name)

		err := outreach.Connect(linkedinBot.Page, lead, message)
		if err != nil {
			fmt.Printf(" %v\n", err)
		} else {
			fmt.Println("Success: Connection sequence completed.")
		}

		fmt.Println("Cooling down...")
		time.Sleep(10 * time.Second)
	}

	fmt.Println("Login Sequence Finished!")
	fmt.Println(" Press Ctrl+C in this terminal to stop the bot.")
	select {}
}
