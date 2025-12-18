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
func NewBrowser(headless bool) (*Browser, error) {
        fmt.Println("🚀 Core: Initializing Stealth Browser...")
        l := launcher.New().
                Headless(headless).
                Leakless(false). // set true if you want Rod to guard against orphaned Chrome
                Delete("enable-automation").
                Set("disable-blink-features", "AutomationControlled").
                Set("window-size", "1920,1080")

        url, err := l.Launch()
        if err != nil {
                return nil, fmt.Errorf("failed to launch browser: %w", err)
        }

        // 2. Connect to the browser
        browser := rod.New().ControlURL(url).MustConnect()

        // 3. Create a new tab and inject stealth scripts
        page := stealth.MustPage(browser)

        fmt.Println("🎭 Core: Stealth scripts injected.")

        return &Browser{
                RodBrowser: browser,
                Page:       page,
        }, nil
  }

// Close cleans up resources
func (b *Browser) Close() {
	if b.RodBrowser != nil {
		b.RodBrowser.MustClose()
	}
}