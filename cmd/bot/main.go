package main

import (
	"fmt"
	"log"

	"linkedin-automation/internal/config"
	"linkedin-automation/internal/core"
)

func main() {
	// 1. Load Config
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
	// defer linkedinBot.Close()

	// 3. Test Navigation (The Login Page)
	fmt.Println("🤖 Bot: Navigating to LinkedIn Login...")
	linkedinBot.Page.MustNavigate("https://www.linkedin.com/login")

	// 4. Wait so you can see it
	// fmt.Println("⏳ Waiting 10 seconds before closing...")
	// time.Sleep(10 * time.Second)
}

