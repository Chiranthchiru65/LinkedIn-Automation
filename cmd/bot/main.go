package main

import (
	"fmt"
	"log"
	
	// Import the config package we just made. 
	// The path depends on your module name (go.mod). 
	// If you named it "linkedin-automation", it looks like this:
	"linkedin-automation/internal/config"
)

func main() {
	// Try to load the config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Print it out to prove it works
	fmt.Println("✅ Config Loaded Successfully!")
	fmt.Printf("   Target: %s in %s\n", cfg.Target.Keywords, cfg.Target.Location)
	fmt.Printf("   User:   %s\n", cfg.Credentials.Username)
	fmt.Printf("   Mode:   Headless=%v\n", cfg.Settings.Headless)
}