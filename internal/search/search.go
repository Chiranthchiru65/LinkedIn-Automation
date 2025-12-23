package search

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"linkedin-automation/internal/models"

	"github.com/go-rod/rod"
)

type SearchEngine struct {
	Page       *rod.Page
	Keywords   string
	Location   string
	MaxResults int
}

func NewEngine(page *rod.Page, keywords, location string) *SearchEngine {
	return &SearchEngine{
		Page:       page,
		Keywords:   keywords,
		Location:   location,
		MaxResults: 5,
	}
}

func (s *SearchEngine) GenerateSearchURL() string {
	baseURL := "https://www.linkedin.com/search/results/people/?"
	params := url.Values{}
	fullQuery := fmt.Sprintf("%s %s", s.Keywords, s.Location)
	params.Set("keywords", fullQuery)
	params.Add("origin", "GLOBAL_SEARCH_HEADER")
	return baseURL + params.Encode()
}

func (s *SearchEngine) Run() ([]models.Lead, error) {
	var leads []models.Lead

	// 1. Navigate
	targetURL := s.GenerateSearchURL()
	fmt.Printf(" Search: Navigating to %s\n", targetURL)
	s.Page.MustNavigate(targetURL)
	s.Page.MustWaitLoad()

	// 2. WAIT for the Main Container (Based on your screenshot)
	// We use the stable "search-results-container" class
	fmt.Println(" Waiting for search container...")
	
	containerSelector := ".search-results-container"
	
	err := rod.Try(func() {
		s.Page.Timeout(10 * time.Second).MustElement(containerSelector)
	})

	if err != nil {
		fmt.Printf("   Current Title: %s\n", s.Page.MustInfo().Title)
		return nil, fmt.Errorf("could not find container: %s", containerSelector)
	}

	// 3. Select the Container
	container := s.Page.MustElement(containerSelector)
	
	// 4. Get all List Items (li)
	// We assume any 'li' inside this container is a result card
	items, err := container.Elements("li")
	if err != nil {
		return nil, fmt.Errorf("container found, but no 'li' items: %v", err)
	}

	fmt.Printf(" Found %d list items. Parsing structure...\n", len(items))

	// 5. Iterate and Extract using "Link Logic"
	for _, item := range items {
		if len(leads) >= s.MaxResults {
			break
		}

		// STRATEGY: Find ANY link (<a>) that contains "/in/" in the href.
		// This ignores class names and looks for the "Profile Link" pattern.
		linkEl, err := item.Element("a[href*='/in/']")
		if err != nil {
			// This 'li' might be a header, footer, or divider. Skip it.
			continue
		}

		profileURL, _ := linkEl.Attribute("href")
		
		// Clean URL (remove query params)
		if idx := strings.Index(*profileURL, "?"); idx != -1 {
			*profileURL = (*profileURL)[:idx]
		}

		// Extract Name
		// The first link is usually the profile picture or name.
		// We try to get text. If empty, we try to find a second link or specific span.
		cleanName := strings.TrimSpace(linkEl.MustText())
		
		// If the first link was just an image (empty text), try to find the text link
		if cleanName == "" {
			// Look for a link that is NOT just an image (has text length > 3)
			links, _ := item.Elements("a[href*='/in/']")
			for _, l := range links {
				t := strings.TrimSpace(l.MustText())
				if len(t) > 2 { // Filter out empty or "..."
					cleanName = strings.Split(t, "\n")[0]
					break
				}
			}
		}

		// Fallback: If still no name, use a placeholder
		if cleanName == "" {
			cleanName = "LinkedIn Member"
		}

		// Skip "LinkedIn Member" (Generic profiles)
		if strings.Contains(cleanName, "LinkedIn Member") {
			continue
		}

		newLead := models.Lead{
			Name:       cleanName,
			ProfileURL: *profileURL,
			Headline:   "Software Engineer", // Placeholder since subtitle class is scrambled
		}
		leads = append(leads, newLead)
		fmt.Printf("   -> Scraped: %s \n      URL: %s\n", newLead.Name, newLead.ProfileURL)
	}

	if len(leads) == 0 {
		return nil, fmt.Errorf("found list items, but extracted 0 leads (check selectors)")
	}

	return leads, nil
}