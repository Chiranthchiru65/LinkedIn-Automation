package auth

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"linkedin-automation/internal/stealth"

	"github.com/go-rod/rod"
)

// Login performs the authentication flow with stealth and error handling
func Login(page *rod.Page, username, password string) error {
	fmt.Println("🔑 Auth: Checking login status...")
	
	// --- PHASE 1: SMART SESSION CHECK ---
	// 1. Navigate to Home (not login) to test cookies
	page.MustNavigate("https://www.linkedin.com/")
	
	// 2. Wait briefly for potential redirect (Feed vs Login)
	time.Sleep(3 * time.Second)

	currentURL := page.MustInfo().URL
	if strings.Contains(currentURL, "feed") || strings.Contains(currentURL, "miniprofile") {
		fmt.Println("✅ Auth: Already logged in! (Cookies detected).")
		return nil
	}

	// 3. If we are not on the feed, ensure we are on the login page
	if !strings.Contains(currentURL, "login") && !strings.Contains(currentURL, "ual") {
		fmt.Println("   - Cookies expired or missing. Navigating to Login...")
		page.MustNavigate("https://www.linkedin.com/login")
		page.MustWaitLoad()
	}

	// --- PHASE 2: HUMAN INTERACTION ---
	
	// 4. Reaction Time (Crucial for Stealth)
	fmt.Println("   - Page loaded. Human reaction pause...")
	reactionTime := time.Duration(rand.Intn(2000)+1500) * time.Millisecond
	time.Sleep(reactionTime)

	// 5. Username (Autofocus usually active, but we type slowly)
	fmt.Println("   - Typing username...")
	userField := page.MustElement("#username")
	stealth.HumanTyping(userField, username)

	// Pause between fields
	time.Sleep(time.Millisecond * 700)

	// 6. Password (Curve Mouse -> Click -> Type)
	fmt.Println("   - Moving mouse to password field...")
	passField := page.MustElement("#password")
	
	stealth.MoveTo(page, passField) 
	stealth.ClickWithRandomDelay(page)
	
	stealth.HumanTyping(passField, password)

	// Pause before submitting
	time.Sleep(time.Millisecond * 500)

	// 7. Submit (Curve Mouse -> Click)
	fmt.Println("   - Moving mouse to login button...")
	loginBtn := page.MustElement("button[type='submit']")
	
	stealth.MoveTo(page, loginBtn)
	stealth.ClickWithRandomDelay(page)

	// --- PHASE 3: VERIFICATION & SECURITY CHECK ---
	
	fmt.Println("⏳ Auth: Credentials submitted. Verifying state...")
	
	// We verify the result for up to 15 seconds
	for i := 0; i < 15; i++ {
		time.Sleep(1 * time.Second)
		url := page.MustInfo().URL
		
		// A. SUCCESS: Reached Feed
		if strings.Contains(url, "feed") || strings.Contains(url, "miniprofile") {
			fmt.Println("✅ Auth: Login Successful! We are on the feed.")
			return nil
		}

		// B. CHALLENGE: 2FA or Captcha
		if strings.Contains(url, "challenge") || strings.Contains(url, "checkpoint") {
			fmt.Println("⚠️  Auth: Security Checkpoint Detected! (2FA/Captcha)")
			fmt.Println("🛑  ACTION REQUIRED: Please solve it manually in the browser window.")
			fmt.Println("⏳ Pausing bot until you solve it...")
			
			// Infinite loop: Wait until user solves it and reaches feed
			for {
				time.Sleep(2 * time.Second)
				if strings.Contains(page.MustInfo().URL, "feed") {
					fmt.Println("✅ Checkpoint passed! Resuming...")
					return nil
				}
			}
		}

		// C. FAILURE: Wrong Password Alert
		// We use Has() to check without crashing if it's missing
		if has, _,_ := page.Has("#error-for-password"); has {
			return fmt.Errorf("❌ Auth Failed: Invalid Credentials detected")
		}
	}

	return fmt.Errorf("❌ Auth Timeout: Login took too long or unknown state")
}