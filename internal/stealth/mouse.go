package stealth

import (
	"math"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto" // Standard proto import
)

// Global tracker for mouse position
var currentX, currentY float64 = 0, 0

// MoveTo mimics human cursor movement
func MoveTo(page *rod.Page, element *rod.Element) {
	box := element.MustShape().Box()
	destX := box.X + box.Width/2
	destY := box.Y + box.Height/2

	// Add random jitter
	jitter := 3.0
	destX += (rand.Float64() * jitter * 2) - jitter
	destY += (rand.Float64() * jitter * 2) - jitter

	moveCurve(page, destX, destY)
}

func moveCurve(page *rod.Page, targetX, targetY float64) {
	startX := currentX
	startY := currentY
	steps := 80

	c1x := startX + (targetX-startX)*rand.Float64()
	c1y := startY + (targetY-startY)*rand.Float64()
	c2x := startX + (targetX-startX)*rand.Float64()
	c2y := startY + (targetY-startY)*rand.Float64()

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)

		x := math.Pow(1-t, 3)*startX + 
			 3*math.Pow(1-t, 2)*t*c1x + 
			 3*(1-t)*math.Pow(t, 2)*c2x + 
			 math.Pow(t, 3)*targetX

		y := math.Pow(1-t, 3)*startY + 
			 3*math.Pow(1-t, 2)*t*c1y + 
			 3*(1-t)*math.Pow(t, 2)*c2y + 
			 math.Pow(t, 3)*targetY

		// --- THE NUCLEAR FIX ---
		// Instead of page.Mouse.Move(x,y), we send the raw protocol command.
		// This bypasses any API mismatch in the helper library.
		err := proto.InputDispatchMouseEvent{
			Type: proto.InputDispatchMouseEventTypeMouseMoved,
			X:    x,
			Y:    y,
		}.Call(page)
		
		if err != nil {
			// If this fails, Chrome itself is broken/disconnected
			continue 
		}

		time.Sleep(time.Millisecond * 3)
	}

	currentX = targetX
	currentY = targetY
}

// ClickWithRandomDelay clicks with a human-like hold duration
func ClickWithRandomDelay(page *rod.Page) {
	// 1. Mouse Down
	proto.InputDispatchMouseEvent{
		Type:       proto.InputDispatchMouseEventTypeMousePressed,
		Button:     proto.InputMouseButtonLeft,
		ClickCount: 1,
		X:          currentX, // Use last known position
		Y:          currentY,
	}.Call(page)

	// 2. Hold
	holdTime := time.Duration(rand.Intn(90)+60) * time.Millisecond
	time.Sleep(holdTime)

	// 3. Mouse Up
	proto.InputDispatchMouseEvent{
		Type:       proto.InputDispatchMouseEventTypeMouseReleased,
		Button:     proto.InputMouseButtonLeft,
		ClickCount: 1,
		X:          currentX,
		Y:          currentY,
	}.Call(page)
}