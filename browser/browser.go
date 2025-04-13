package browser

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// Browser represents a web browser instance
type Browser struct {
	browser *rod.Browser
	page    *rod.Page
}

// NewBrowser creates a new browser instance
func NewBrowser() (*Browser, error) {
	LogInfo("Creating new browser instance")

	LogDebug("Configuring browser launcher")
	u := launcher.New().
		Headless(false).
		Set("window-size", "1280,800").
		Delete("enable-automation").
		MustLaunch()

	LogDebug("Connecting to browser")
	browser := rod.New().ControlURL(u).MustConnect()

	LogDebug("Creating new page")
	page, err := browser.Page(proto.TargetCreateTarget{URL: "about:blank"})
	if err != nil {
		LogError("Failed to create page: %v", err)
		return nil, fmt.Errorf("failed to create page: %v", err)
	}

	page = page.Timeout(120 * time.Second)

	LogInfo("Browser instance created successfully")
	return &Browser{
		browser: browser,
		page:    page,
	}, nil
}

// Search performs a web search and returns the results
func (b *Browser) Search(query string) (string, []string, error) {
	LogInfo("Starting web search for query: %s", query)

	b.page = b.page.Timeout(180 * time.Second)

	navigationRetries := 5
	var err error
	for i := 0; i < navigationRetries; i++ {
		LogInfo("Navigation attempt %d/%d to Google", i+1, navigationRetries)
		err = b.page.Navigate("https://www.google.com")
		if err == nil {
			err = b.page.WaitLoad()
			if err == nil {
				break
			}
		}
		LogError("Navigation attempt %d failed: %v", i+1, err)
		time.Sleep(time.Duration(3*i+1) * time.Second)
	}
	if err != nil {
		return "", nil, fmt.Errorf("failed to navigate to search engine after %d attempts: %v", navigationRetries, err)
	}

	searchBox, err := b.page.Timeout(30 * time.Second).Element(`input[name="q"]`)
	if err != nil {
		LogError("Failed to find search box: %v", err)
		return "", nil, fmt.Errorf("failed to find search box: %v", err)
	}

	err = searchBox.Focus()
	if err != nil {
		LogError("Failed to focus search box: %v", err)
		return "", nil, fmt.Errorf("failed to focus search box: %v", err)
	}

	visible, err := searchBox.Visible()
	if err != nil {
		LogError("Failed to check search box visibility: %v", err)
		return "", nil, fmt.Errorf("failed to check search box visibility: %v", err)
	}
	if !visible {
		LogError("Search box is not visible")
		return "", nil, fmt.Errorf("search box is not visible")
	}

	inputRetries := 5
	for i := 0; i < inputRetries; i++ {
		err = searchBox.Input(query)
		if err == nil {
			break
		}
		LogError("Input attempt %d failed: %v", i+1, err)
		time.Sleep(time.Second)
	}
	if err != nil {
		return "", nil, fmt.Errorf("failed to type search query after %d attempts: %v", inputRetries, err)
	}

	time.Sleep(2 * time.Second)
	err = b.page.Keyboard.Press('\n')
	if err != nil {
		LogError("Failed to press Enter to submit search: %v", err)
		return "", nil, fmt.Errorf("failed to submit search: %v", err)
	}

	_, err = b.page.Timeout(30 * time.Second).Element(`#search`)
	if err != nil {
		LogError("Failed to load search results: %v", err)
		return "", nil, fmt.Errorf("failed to load search results: %v", err)
	}

	results, err := b.page.Eval(`() => {
		const results = [];
		document.querySelectorAll('#search .g').forEach((result, index) => {
			if (index < 5) {
				const titleElement = result.querySelector('h3');
				const linkElement = result.querySelector('a');
				const snippetElement = result.querySelector('.VwiC3b');

				if (titleElement && linkElement) {
					const title = titleElement.textContent.trim();
					const link = linkElement.href;
					const snippet = snippetElement ? snippetElement.textContent.trim() : '';
					results.push({ title, link, snippet });
				}
			}
		});
		return results;
	}`)
	if err != nil {
		LogError("Failed to extract search results: %v", err)
		return "", nil, fmt.Errorf("failed to extract search results: %v", err)
	}

	var sources []string
	var summary strings.Builder
	summary.WriteString("Search Results:\n\n")

	var resultsArray []map[string]interface{}
	err = results.Value.Unmarshal(&resultsArray)
	if err != nil {
		LogError("Failed to parse search results: %v", err)
		return "", nil, fmt.Errorf("failed to parse search results: %v", err)
	}

	if len(resultsArray) == 0 {
		LogError("No search results found")
		return "No search results found. The search engine may have changed its structure or blocked automated access.", nil, nil
	}

	for i, result := range resultsArray {
		title, titleExists := result["title"]
		link, linkExists := result["link"]
		snippet, snippetExists := result["snippet"]
		if !titleExists || !linkExists || !snippetExists {
			continue
		}

		titleStr, ok := title.(string)
		if !ok {
			continue
		}
		linkStr, ok := link.(string)
		if !ok {
			continue
		}
		snippetStr, ok := snippet.(string)
		if !ok {
			continue
		}

		sources = append(sources, linkStr)
		summary.WriteString(fmt.Sprintf("%d. %s\n", i+1, titleStr))
		summary.WriteString(fmt.Sprintf("   %s\n", snippetStr))
		summary.WriteString(fmt.Sprintf("   Source: %s\n\n", linkStr))
	}

	return summary.String(), sources, nil
}

// VisitPage visits a specific URL and extracts its content
func (b *Browser) VisitPage(url string) (string, error) {
	LogInfo("Visiting page: %s", url)

	err := b.page.Navigate(url)
	if err != nil {
		LogError("Failed to navigate to page: %v", err)
		return "", fmt.Errorf("failed to navigate to page: %v", err)
	}

	err = b.page.WaitLoad()
	if err != nil {
		LogError("Failed to load page: %v", err)
		return "", fmt.Errorf("failed to load page: %v", err)
	}

	content, err := b.page.Eval(`() => {
		document.querySelectorAll('script, style').forEach(e => e.remove());
		const mainContent = document.querySelector('main, article, .content, #content, .main, #main') || document.body;
		return mainContent.innerText;
	}`)
	if err != nil {
		LogError("Failed to extract page content: %v", err)
		return "", fmt.Errorf("failed to extract page content: %v", err)
	}

	contentText := content.Value.String()
	LogInfo("Successfully extracted %d characters from page", len(contentText))
	return contentText, nil
}

// Close closes the browser and cleans up resources
func (b *Browser) Close() error {
	LogInfo("Closing browser")

	if err := b.page.Close(); err != nil {
		LogError("Failed to close page: %v", err)
		return fmt.Errorf("failed to close page: %v", err)
	}

	if err := b.browser.Close(); err != nil {
		LogError("Failed to close browser: %v", err)
		return fmt.Errorf("failed to close browser: %v", err)
	}

	LogInfo("Browser closed successfully")
	return nil
}
