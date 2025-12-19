package auth

import (
	"fmt"
	"math/rand"
	"time"

	"linkedin-automation/internal/stealth"

	"github.com/go-rod/rod"
)

func Login(page *rod.Page, username, password string) error {
	fmt.Println("🔑 Auth: Starting login flow...")
	
	// 1. Navigate
	page.MustNavigate("https://www.linkedin.com/login")
	page.MustWaitLoad()
	
	// 2. HUMAN REACTION TIME (Critical for Autofocus)
	// We wait 2-4 seconds. A bot types instantly. A human looks at the screen.
	fmt.Println("   - Page loaded. Human reaction pause...")
	reactionTime := time.Duration(rand.Intn(2000)+2000) * time.Millisecond
	time.Sleep(reactionTime)

	// 3. TYPE USERNAME (No Mouse Move needed!)
	// The cursor is already blinking here.
	fmt.Println("   - Typing username...")
	userField := page.MustElement("#username")
	stealth.HumanTyping(userField, username)

	// Pause: "Thinking" about the password
	time.Sleep(time.Millisecond * 700)

	// 4. PASSWORD (Move -> Click -> Type)
	fmt.Println("   - Moving mouse to password field...")
	passField := page.MustElement("#password")
	
	// This will curve from Top-Left (0,0) to the Password Box
	stealth.MoveTo(page, passField) 
	stealth.ClickWithRandomDelay(page)
	
	stealth.HumanTyping(passField, password)

	// Pause: Visually checking the button
	time.Sleep(time.Millisecond * 500)

	// 5. SUBMIT (Move -> Click)
	fmt.Println("   - Moving mouse to login button...")
	loginBtn := page.MustElement("button[type='submit']")
	
	// This will curve from Password Box to Submit Button
	stealth.MoveTo(page, loginBtn)
	stealth.ClickWithRandomDelay(page)

	// 6. Wait for Login Success
	fmt.Println("✅ Auth: Credentials submitted. Waiting for navigation...")
	time.Sleep(5 * time.Second)
	
	return nil
}