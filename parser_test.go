package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseNetscapeBookmarks(t *testing.T) {
	tests := []struct {
		name        string
		inputHTML   string
		expected    []BookmarkItem
		expectError bool
	}{
		{
			name: "Simple bookmark",
			inputHTML: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
				<HTML>
				<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
				<TITLE>Bookmarks</TITLE>
				<H1>Bookmarks</H1>
				<DL><p>
					<DT><A HREF="https://example.com" ADD_DATE="1678886400" LAST_MODIFIED="1678886401" TAGS="tag1,tag2" ICON_URI="https://example.com/icon.png" ICON="data:image/png;base64,iVBORw0KGgo=" ID="test_id_1">Example</A>
				</DL><p>
				</HTML>`,
			expected: []BookmarkItem{
				{Type: "bookmark", Name: "Example", Href: "https://example.com", AddDate: "1678886400", LastModified: "1678886401", Tags: "tag1,tag2", IconURI: "https://example.com/icon.png", Icon: "data:image/png;base64,iVBORw0KGgo=", ID: "test_id_1"},
			},
			expectError: false,
		},
		{
			name: "Simple folder",
			inputHTML: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
				<HTML><TITLE>Bookmarks</TITLE><H1>Bookmarks</H1>
				<DL><p>
					<DT><H3 ADD_DATE="1678886400" LAST_MODIFIED="1678886401" ID="folder_id_1">My Folder</H3>
					<DL><p>
						<DT><A HREF="https://child.com">Child Link</A>
					</DL><p>
				</DL><p>
				</HTML>`,
			expected: []BookmarkItem{
				{
					Type:         "folder",
					Name:         "My Folder",
					AddDate:      "1678886400",
					LastModified: "1678886401",
					ID:           "folder_id_1",
					Children: []BookmarkItem{
						{Type: "bookmark", Name: "Child Link", Href: "https://child.com"},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Nested folders",
			inputHTML: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
				<HTML><TITLE>Bookmarks</TITLE><H1>Bookmarks</H1>
				<DL><p>
					<DT><H3>Parent</H3>
					<DL><p>
						<DT><H3>Child</H3>
						<DL><p>
							<DT><A HREF="https://grandchild.com">Grandchild Link</A>
						</DL><p>
					</DL><p>
				</DL><p>
				</HTML>`,
			expected: []BookmarkItem{
				{
					Type: "folder", Name: "Parent",
					Children: []BookmarkItem{
						{
							Type: "folder", Name: "Child",
							Children: []BookmarkItem{
								{Type: "bookmark", Name: "Grandchild Link", Href: "https://grandchild.com"},
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Separator",
			inputHTML: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
				<HTML><TITLE>Bookmarks</TITLE><H1>Bookmarks</H1>
				<DL><p>
					<DT><A HREF="https://site1.com">Site 1</A>
					<HR>
					<DT><A HREF="https://site2.com">Site 2</A>
				</DL><p>
				</HTML>`,
			expected: []BookmarkItem{
				{Type: "bookmark", Name: "Site 1", Href: "https://site1.com"},
				{Type: "separator"},
				{Type: "bookmark", Name: "Site 2", Href: "https://site2.com"},
			},
			expectError: false,
		},
		{
			name: "Empty DL",
			inputHTML: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
				<HTML><TITLE>Bookmarks</TITLE><H1>Bookmarks</H1>
				<DL><p>
				</DL><p>
				</HTML>`,
			expected:    []BookmarkItem{},
			expectError: false,
		},
		{
			name: "No root DL in body (should return empty)",
			inputHTML: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
				<HTML><TITLE>Bookmarks</TITLE><H1>Bookmarks</H1>
				<BODY>
					<P>No bookmarks here.</P>
				</BODY>
				</HTML>`,
			expected:    []BookmarkItem{}, // Current behavior is to return empty if no root DL
			expectError: false,
		},
		{
			name: "Mixed content with P and DT tags",
			inputHTML: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
				<HTML><TITLE>Bookmarks</TITLE><H1>Bookmarks</H1>
				<DL><P>
					<DT><H3>Folder 1</H3>
					<DL><P>
						<DT><A HREF="http://link1.com">Link 1</A></P>
						<P><DT><A HREF="http://link2.com">Link 2</A></P>
					</DL>
					<DT><HR>
					<DT><A HREF="http://link3.com">Link 3</A>
				</DL></P>
				</HTML>`,
			expected: []BookmarkItem{
				{
					Type: "folder", Name: "Folder 1",
					Children: []BookmarkItem{
						{Type: "bookmark", Name: "Link 1", Href: "http://link1.com"},
						{Type: "bookmark", Name: "Link 2", Href: "http://link2.com"},
					},
				},
				{Type: "separator"},
				{Type: "bookmark", Name: "Link 3", Href: "http://link3.com"},
			},
		},
		{
			name: "Bookmark with only Href and Name",
			inputHTML: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
				<HTML><TITLE>Bookmarks</TITLE><H1>Bookmarks</H1>
				<DL><p>
					<DT><A HREF="https://minimal.com">Minimal</A>
				</DL><p>
				</HTML>`,
			expected: []BookmarkItem{
				{Type: "bookmark", Name: "Minimal", Href: "https://minimal.com"},
			},
			expectError: false,
		},
		{
			name: "Folder with attributes",
			inputHTML: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
				<HTML><TITLE>Bookmarks</TITLE><H1>Bookmarks</H1>
				<DL><p>
					<DT><H3 ADD_DATE="123" LAST_MODIFIED="456" ID="f1">Folder With Attrs</H3>
					<DL><p></DL><p>
				</DL><p>
				</HTML>`,
			expected: []BookmarkItem{
				{Type: "folder", Name: "Folder With Attrs", AddDate: "123", LastModified: "456", ID: "f1", Children: []BookmarkItem{}},
			},
			expectError: false,
		},
		// Add more test cases for edge cases, malformed HTML (if specific errors are expected), etc.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.inputHTML)
			got, err := ParseNetscapeBookmarks(reader)

			if (err != nil) != tt.expectError {
				t.Errorf("ParseNetscapeBookmarks() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				// For easier debugging, print them out if they don't match
				t.Errorf("ParseNetscapeBookmarks() got = %v, want %v", got, tt.expected)
				// For very complex structures, you might want to marshal to JSON and compare strings
				// gotJSON, _ := json.MarshalIndent(got, "", "  ")
				// expectedJSON, _ := json.MarshalIndent(tt.expected, "", "  ")
				// t.Errorf("ParseNetscapeBookmarks() gotJSON = %s, wantJSON %s", string(gotJSON), string(expectedJSON))
			}
		})
	}
}
