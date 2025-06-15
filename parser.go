package main

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// BookmarkItem represents a single bookmark or folder
type BookmarkItem struct {
	Type         string         `json:"type"` // "bookmark" or "folder"
	Name         string         `json:"name"`
	Href         string         `json:"href,omitempty"`
	Tags         string         `json:"tags,omitempty"`
	ID           string         `json:"id,omitempty"`
	AddDate      string         `json:"add_date,omitempty"`      // Date string, format may vary
	LastModified string         `json:"last_modified,omitempty"` // Date string, format may vary
	Icon         string         `json:"icon,omitempty"`          // Base64 encoded icon data or a data URI
	IconURI      string         `json:"icon_uri,omitempty"`      // Actual URI for the icon
	Children     []BookmarkItem `json:"children,omitempty"`
}

// removeUnwantedElements preprocesses the HTML content string to remove specific
// unwanted tags like <DT> and <p>.
// This simplifies the structure for subsequent parsing.
func removeUnwantedElements(fileContent string) string {
	content := fileContent
	content = strings.ReplaceAll(content, "<p>", "")
	content = strings.ReplaceAll(content, "<P>", "")
	content = strings.ReplaceAll(content, "</p>", "")
	content = strings.ReplaceAll(content, "</P>", "")
	content = strings.ReplaceAll(content, "<DT>", "")
	content = strings.ReplaceAll(content, "<dt>", "")
	// Note: </DT> and </dt> are typically not used in this context, so not explicitly removed.
	return content
}

// getAttribute extracts the value of a specific attribute from an HTML node.
func getAttribute(n *html.Node, key string) string {
	if n == nil {
		return ""
	}
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// getNodeText recursively extracts and concatenates all text data from a node and its children.
func getNodeText(n *html.Node) string {
	if n == nil {
		return ""
	}
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var buf strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		buf.WriteString(getNodeText(c))
	}
	return strings.TrimSpace(buf.String())
}

// parseDLNode recursively parses a <DL> HTML node and its children
// to extract bookmark items and folders.
func parseDLNode(dlNode *html.Node) ([]BookmarkItem, error) {
	items := []BookmarkItem{} // Initialize as empty non-nil slice
	for c := dlNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue // Skip text nodes, comments, etc. at this level
		}

		switch c.Data {
		case "a": // Bookmark link
			item := BookmarkItem{
				Type:         "bookmark",
				Name:         getNodeText(c),
				Href:         getAttribute(c, "href"),
				Tags:         getAttribute(c, "tags"),
				ID:           getAttribute(c, "id"),
				AddDate:      getAttribute(c, "add_date"),
				LastModified: getAttribute(c, "last_modified"),
				Icon:         getAttribute(c, "icon"),
				IconURI:      getAttribute(c, "icon_uri"),
			}
			items = append(items, item)

		case "hr": // Separator
			item := BookmarkItem{
				Type: "separator",
			}
			items = append(items, item)

		case "h3": // Folder heading
			folderName := getNodeText(c)
			folderItem := BookmarkItem{
				Type:         "folder",
				Name:         folderName,
				ID:           getAttribute(c, "id"),
				AddDate:      getAttribute(c, "add_date"),      // Folders can also have these
				LastModified: getAttribute(c, "last_modified"), // attributes in some exports
				// Icon and IconURI are less common for folders but can be added if needed
			}

			// The <DL> for the folder's children should be the next element sibling
			var folderDLNode *html.Node
			nextNode := c.NextSibling
			for nextNode != nil {
				if nextNode.Type == html.ElementNode {
					if nextNode.Data == "dl" {
						folderDLNode = nextNode
					}
					break // Found next element, stop searching
				}
				nextNode = nextNode.NextSibling
			}

			if folderDLNode != nil {
				children, err := parseDLNode(folderDLNode)
				if err != nil {
					return nil, fmt.Errorf("parsing children for folder '%s': %w", folderName, err)
				}
				folderItem.Children = children // parseDLNode now guarantees non-nil children
				c = folderDLNode               // Advance the loop cursor past the processed <DL>
			} else { // If no <DL> follows <H3>, it's an empty folder
				folderItem.Children = []BookmarkItem{}
			}
			items = append(items, folderItem)
		}
	}
	return items, nil
}

// ParseNetscapeBookmarks reads from the reader and parses the bookmark data.
func ParseNetscapeBookmarks(reader io.Reader) ([]BookmarkItem, error) {
	byteContent, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}
	fileContent := string(byteContent)

	// Preprocess to remove <p> and <DT> tags which can interfere with simple hierarchy parsing
	htmlContent := removeUnwantedElements(fileContent)
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML content: %w", err)
	}

	items := []BookmarkItem{} // Initialize as empty non-nil slice
	var findBodyDL func(*html.Node) *html.Node
	findBodyDL = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.Data == "body" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "dl" {
					return c
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if found := findBodyDL(c); found != nil {
				return found
			}
		}
		return nil
	}
	rootDL := findBodyDL(doc)
	if rootDL != nil {
		parsedItems, parseErr := parseDLNode(rootDL)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse bookmark structure: %w", parseErr)
		}
		items = parsedItems
	} else {
		// No root <DL> found, items will remain the initialized empty slice.
	}
	return items, nil
}
