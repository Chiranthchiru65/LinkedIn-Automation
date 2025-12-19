package stealth

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
)

// HumanTyping simulates typing into an element with random delays
func HumanTyping(element *rod.Element, text string) {
	// Ensure random seed is seeded (usually done once in main, but safe here)
	// Note: In newer Go versions, rand is auto-seeded, but this doesn't hurt.
	
	for _, char := range text {
		// 1. Type the character
		element.MustInput(string(char))

		// 2. Random delay (between 50ms and 200ms)
		// This simulates the varying speed of human fingers
		delay := time.Duration(rand.Intn(150)+50) * time.Millisecond
		time.Sleep(delay)
	}
}