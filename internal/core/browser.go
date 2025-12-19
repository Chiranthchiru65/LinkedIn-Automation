package core

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

// Browser is our wrapper around the Rod browser
type Browser struct {
	RodBrowser *rod.Browser
	Page       *rod.Page
}

// NewBrowser initializes a chrome instance with stealth settings
// internal/core/browser.go

// internal/core/browser.go

func NewBrowser(headless bool) (*Browser, error) {
        fmt.Println("🚀 Core: Initializing Stealth Browser...")

        l := launcher.New().
                Headless(headless).
                Leakless(false).
                Delete("enable-automation").
                Set("disable-blink-features", "AutomationControlled").
                Set("window-size", "1920,1080").
                Set("user-data-dir", `C:\temp\rod-profile`). // avoid profile lock
                Bin(`C:\Program Files\Google\Chrome\Application\chrome.exe`)

        fmt.Println("🔧 Launching Chrome...")
        url, err := l.Launch()
        if err != nil {
                return nil, fmt.Errorf("failed to launch browser: %w", err)
        }
        fmt.Println("✅ Chrome launched:", url)

        fmt.Println("🔌 Connecting to Chrome...")
        browser := rod.New().ControlURL(url)
        if err := browser.Connect(); err != nil {
                return nil, fmt.Errorf("failed to connect: %w", err)
        }

        page := stealth.MustPage(browser)
        fmt.Println("🎭 Core: Stealth scripts injected.")

        return &Browser{RodBrowser: browser, Page: page}, nil
  }

// Close cleans up resources
func (b *Browser) Close() {
	if b.RodBrowser != nil {
		b.RodBrowser.MustClose()
	}
}