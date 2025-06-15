# Netscape Bookmarks to JSON Converter

`netscape-bookmarks-json-converter` is a command-line tool written in Go to parse Netscape-format bookmark HTML files and convert them into a structured JSON output.

## Features

- Parses Netscape bookmark HTML files.
- Converts bookmarks and folders into a hierarchical JSON structure.
- Supports reading from a file or standard input.
- Handles nested folders.

## Installation

Ensure you have Go installed (version 1.18 or newer is recommended).

You can install the `netscape-bookmarks-json-converter` using `go install`:
```bash
go install github.com/ocodo/netscape-bookmarks-json-converter@latest
```
This will download the source code, compile it, and place the executable in your `$GOPATH/bin` directory (or `$HOME/go/bin` if `GOPATH` is not set). Make sure this directory is in your system's `PATH`.

Alternatively, you can clone the repository and build it manually:
```bash
git clone https://github.com/ocodo/netscape-bookmarks-json-converter.git
cd netscape-bookmarks-json-converter
go build
```
This will create an executable named `netscape-bookmarks-json-converter` (or `netscape-bookmarks-json-converter.exe` on Windows) in the current directory.

## Usage

The tool can read from a specified file or from standard input.

```
Usage: netscape-bookmarks-json-converter -f <filepath>
   or: cat bookmarks.html | netscape-bookmarks-json-converter

Options:
  -f string
    	Path to the Netscape bookmark file. Reads from stdin if not provided.
```

### Examples

1.  **Convert a bookmark file:**
    ```bash
    netscape-bookmarks-json-converter -f /path/to/your/bookmarks.html > output.json
    ```

2.  **Convert from standard input:**
    ```bash
    cat /path/to/your/bookmarks.html | netscape-bookmarks-json-converter > output.json
    ```

## Output JSON Structure

The output is an array of bookmark items. Each item can be a "bookmark", "folder", or "separator".

```json
[
  {
    "type": "folder",
    "name": "My Folder",
    "id": "folder_id_123",
    "add_date": "1678886400",
    "last_modified": "1678886401",
    "children": [
      {
        "type": "bookmark",
        "name": "Example Site",
        "href": "https://example.com",
        "tags": "tag1,tag2",
        "id": "bookmark_id_456",
        "add_date": "1678886400",
        "last_modified": "1678886401",
        "icon_uri": "https://example.com/icon.png",
        "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
      }
    ]
  },
  {
    "type": "separator"
  },
  {
    "type": "bookmark",
    "name": "Another Site",
    "href": "https://anothersite.org",
    "add_date": "1578886400"
  }
]
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for bugs, feature requests, or improvements.

## License

This project is currently unlicensed. Please consider adding a license file (e.g., MIT, Apache 2.0).
