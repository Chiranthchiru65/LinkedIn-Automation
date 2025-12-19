package main

import (
	"fmt"
	"log"

	"linkedin-automation/internal/auth" // <--- Import the Auth module
	"linkedin-automation/internal/config"
	"linkedin-automation/internal/core"
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

	// 4. Block Forever (Success State)
	fmt.Println("🎉 Login Sequence Finished!")
	fmt.Println("🛑 Press Ctrl+C in this terminal to stop the bot.")
	select {} 
}