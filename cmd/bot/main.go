package main

import (
	"fmt"
	"log"
	"time"

	"linkedin-automation/internal/auth" // <--- Import the Auth module
	"linkedin-automation/internal/config"
	"linkedin-automation/internal/core"
	outreach "linkedin-automation/internal/outreact"
	"linkedin-automation/internal/search"
)

func main() {
	// 1. Load Config (Getting credentials from .env)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Println("✅ Config Loaded")

	// 2. Initialize Browser
	linkedinBot, err := core.NewBrowser(cfg.Settings.Headless)
	if err != nil {
		log.Fatalf("Failed to initialize browser: %v", err)
	}
	
	// We keep the browser open to watch the magic
	// defer linkedinBot.Close() 

	// 3. EXECUTE LOGIN (The new part)
	// We pass the browser page + username + password to our login function
	fmt.Println("🤖 Bot: initiating login sequence...")
	
	err = auth.Login(linkedinBot.Page, cfg.Credentials.Username, cfg.Credentials.Password)
	if err != nil {
		log.Fatalf("❌ Login failed: %v", err)
	}
	fmt.Println("🚀 Core: Initializing Search Engine...")
	
	// Initialize Search with data from Config
	// (Ensure your config struct has these, or hardcode strings for testing)
	// engine := search.NewEngine(linkedinBot.Page, "Software Engineer", "Bangalore")
	
	// // Start the search
	// engine.Run()
	fmt.Println("🚀 Core: Initializing Search Engine...")
	engine := search.NewEngine(linkedinBot.Page, "Software Engineer", "Bangalore")
	
	// Run the Harvest
	leads, err := engine.Run()
	if err != nil {
		log.Fatalf("❌ Search failed: %v", err)
	}

	fmt.Println("------------------------------------------------")
	fmt.Printf("🎉 HARVEST COMPLETE! Found %d leads:\n", len(leads))
	for i, lead := range leads {
		// Safety Limit: Only try 2 people for this test run
		if i >= 2 { 
			break 
		}

		fmt.Printf("\n--- Processing Lead %d/%d: %s ---\n", i+1, len(leads), lead.Name)

		// Create a personalized message
		// Note: LinkedIn limits notes to 300 chars
		message := fmt.Sprintf("Hi %s, I found your profile while searching for Software Engineers in Bangalore. I am building a Golang automation tool and would love to connect!", lead.Name)

		// Execute the connection logic
		err := outreach.Connect(linkedinBot.Page, lead, message)
		if err != nil {
			fmt.Printf("⚠️ %v\n", err)
		} else {
			fmt.Println("✅ Success: Connection sequence completed.")
		}

		// COOLDOWN: Wait 10 seconds before the next person
		fmt.Println("⏳ Cooling down...")
		time.Sleep(10 * time.Second)
	}

	// 4. Block Forever (Success State)
	fmt.Println("🎉 Login Sequence Finished!")
	fmt.Println("🛑 Press Ctrl+C in this terminal to stop the bot.")
	select {} 
}