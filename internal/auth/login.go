package auth

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"linkedin-automation/internal/stealth"

	"github.com/go-rod/rod"
)

func Login(page *rod.Page, username, password string) error {
	fmt.Println("🔑 Auth: Checking login status...")
	
	// 1. Navigate to Home
	// If cookies are valid, this redirects straight to /feed
	page.MustNavigate("https://www.linkedin.com/")
	page.MustWaitLoad()
	
	// 2. CHECK: Are we already logged in?
	// We check if the URL contains "feed" OR if the main profile icon exists
	// We wait briefly (2s) to see where the redirect lands us
	time.Sleep(2 * time.Second)

	currentURL := page.MustInfo().URL
	if strings.Contains(currentURL, "feed") || strings.Contains(currentURL, "miniprofile") {
		fmt.Println("✅ Auth: Already logged in! Skipping credentials.")
		return nil
	}

	// 3. If not on feed, we must be on login (or forced to login page)
	// Ensure we are explicitly on the login page now
	if !strings.Contains(currentURL, "login") {
		page.MustNavigate("https://www.linkedin.com/login")
		page.MustWaitLoad()
	}
	
	// --- START EXISTING LOGIN LOGIC ---
	
	fmt.Println("   - Page loaded. Human reaction pause...")
	reactionTime := time.Duration(rand.Intn(2000)+2000) * time.Millisecond
	time.Sleep(reactionTime)

	fmt.Println("   - Typing username...")
	userField := page.MustElement("#username")
	stealth.HumanTyping(userField, username)

	time.Sleep(time.Millisecond * 700)

	fmt.Println("   - Moving mouse to password field...")
	passField := page.MustElement("#password")
	stealth.MoveTo(page, passField) 
	stealth.ClickWithRandomDelay(page)
	stealth.HumanTyping(passField, password)

	time.Sleep(time.Millisecond * 500)

	fmt.Println("   - Moving mouse to login button...")
	loginBtn := page.MustElement("button[type='submit']")
	stealth.MoveTo(page, loginBtn)
	stealth.ClickWithRandomDelay(page)

	fmt.Println("✅ Auth: Credentials submitted. Waiting for navigation...")
	time.Sleep(5 * time.Second)
	
	return nil
}