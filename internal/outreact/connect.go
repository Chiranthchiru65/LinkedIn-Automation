package outreach

import (
	"fmt"
	"strings"
	"time"

	"linkedin-automation/internal/models"
	"linkedin-automation/internal/stealth"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

// Connect visits a profile and sends a connection request
func Connect(page *rod.Page, lead models.Lead, note string) error {
	fmt.Printf("🏃 Outreach: Visiting %s...\n", lead.Name)
	
	page.MustNavigate(lead.ProfileURL)
	page.MustWaitLoad()
	time.Sleep(3 * time.Second)

	// Scope to Main Content (ignore Navbar)
	mainContent, err := page.Element("main")
	if err != nil {
		return fmt.Errorf("failed to find main content")
	}

	var connectBtn *rod.Element

	// --- STRATEGY 1: CHECK PRIMARY BUTTON ("Connect" or "+ Add") ---
	fmt.Println("   - Checking for primary 'Connect/Add' button...")
	
	btn, err := mainContent.Element("button[aria-label*='Invite'][aria-label*='connect']")
	if err == nil {
		fmt.Println("   - Found Primary Connect Button!")
		connectBtn = btn
	} else {
		buttons, _ := mainContent.Elements("button")
		for _, b := range buttons {
			if strings.TrimSpace(b.MustText()) == "Connect" {
				connectBtn = b
				break
			}
		}
	}

	// --- STRATEGY 2: CHECK "MORE" MENU (If primary not found) ---
	if connectBtn == nil {
		fmt.Println("   - Primary button not found. Checking 'More' menu...")
		
		moreBtn, err := mainContent.Element("button[aria-label='More actions']")
		if err != nil {
			moreBtn, err = mainContent.ElementR("button", "More")
		}

		if moreBtn != nil {
			stealth.MoveTo(page, moreBtn)
			stealth.ClickWithRandomDelay(page)
			time.Sleep(1 * time.Second)

			items, _ := page.Elements(".artdeco-dropdown__item")
			for _, item := range items {
				text := strings.ToLower(item.MustText())
				if (strings.Contains(text, "connect") || strings.Contains(text, "invite")) && !strings.Contains(text, "remove") {
					connectBtn = item
					break
				}
			}
			
			if connectBtn == nil {
				page.Keyboard.Press(input.Escape)
			}
		}
	}

	// --- EXECUTE CONNECTION ---
	if connectBtn == nil {
		return fmt.Errorf("skipped: No Connect/Add button found")
	}

	fmt.Println("   - Clicking Connect/Add...")
	stealth.MoveTo(page, connectBtn)
	stealth.ClickWithRandomDelay(page)
	time.Sleep(2 * time.Second)

	// --- HANDLE "ADD A NOTE" MODAL ---
	addNoteBtn, err := page.Element("button[aria-label='Add a note']")
	if err == nil {
		fmt.Println("   - Modal 1: Clicking 'Add a note'...")
		stealth.MoveTo(page, addNoteBtn)
		stealth.ClickWithRandomDelay(page)
		time.Sleep(1 * time.Second)

		fmt.Println("   - Modal 2: Typing message...")
		textArea, err := page.Element("textarea[name='message']")
		if err == nil {
			stealth.HumanTyping(textArea, note)
			time.Sleep(1 * time.Second)
			
			// --- STEP 7: SEND (ACTIVATED) ---
			// We try finding the button by Text ("Send") or Label ("Send now")
			// Your screenshot shows the button simply says "Send"
			sendBtn, err := page.ElementR("button", "Send")
			if err != nil {
				sendBtn, _ = page.Element("button[aria-label='Send now']")
			}

			if sendBtn != nil {
				fmt.Println("   - Found Send button. Sending request...")
				stealth.MoveTo(page, sendBtn)
				
				// 🔥🔥 THE SAFETY LOCK IS OFF! THIS WILL REALLY SEND! 🔥🔥
				stealth.ClickWithRandomDelay(page) 
				
				fmt.Println("   ✅ Connection Request SENT!")
				time.Sleep(2 * time.Second) // Wait for modal to close
			} else {
				fmt.Println("   ⚠️ Error: Could not find 'Send' button!")
			}
		}
	} else {
		// Fallback if the "Add a note" button was skipped
		if textArea, err := page.Element("textarea[name='message']"); err == nil {
             stealth.HumanTyping(textArea, note)
             sendBtn, _ := page.ElementR("button", "Send")
             if sendBtn != nil {
                 stealth.MoveTo(page, sendBtn)
                 stealth.ClickWithRandomDelay(page)
                 fmt.Println("   ✅ Connection Request SENT!")
             }
        } else {
		     fmt.Println("   - No 'Add Note' flow detected. Request might be sent.")
        }
	}

	// Close any lingering modals
	page.Keyboard.Press(input.Escape)
	return nil
}