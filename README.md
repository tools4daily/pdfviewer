# PDF Reader

A modern PDF reader built with Fyne + go-fitz, featuring multi-tab support.

## Features

### Core Features
- ✅ Open and read PDF files
- ✅ **Multi-tab support** - Open multiple PDF files simultaneously, switch between tabs ⭐ **v1.1 New**
- ✅ Page navigation (Previous/Next/First/Last/Jump to page)
- ✅ Zoom functionality (50%-300%: Zoom in/Zoom out/Reset)
- ✅ Keyboard shortcuts support
- ✅ Cross-platform (Windows/Linux/macOS)

### Multi-tab Features (v1.1) ⭐
- ✅ **Open multiple PDFs simultaneously** - Each file in independent tab
- ✅ **Smart tab management** - Empty tabs automatically reused, avoid redundant tabs ⭐ **v1.1.1 Optimized**
- ✅ **Tab switching** - Click tabs or use shortcuts to switch
- ✅ **Independent state management** - Each tab has independent page number and zoom level
- ✅ **Tab management** - Create and close tabs
- ✅ **Filename display** - Tab title shows filename

### Interface Enhancements (v1.0)
- ✅ **Standard menu bar**
  - File menu: Open, New Tab, Close Tab, Exit
  - View menu: Navigation and zoom options
  - Help menu: Keyboard shortcuts, About information
- ✅ **Icon-based toolbar** - Beautiful icon buttons instead of text
- ✅ **Mouse wheel page flipping** - Scroll up/down to flip pages
- ✅ **Double-click to open file** - Double-click blank area to open file selection dialog
- ✅ **Enhanced status bar** - Displays filename, page number, zoom ratio, file size
- ✅ **Beautiful program icon** - Modern minimalist PDF document icon

## Windows Build

### Prerequisites

1. Install Go 1.19 or higher
2. Install GCC compiler (TDM-GCC or MinGW-w64 recommended)
   - Download TDM-GCC: https://jmeubank.github.io/tdm-gcc/
   - Or download MinGW-w64: https://www.mingw-w64.org/

### Build Steps

```bash
# 1. Navigate to project root directory
cd /path/to/pdfviewer
# 2. Download dependencies
go mod download

# 3. Build (execute in project root directory)
go build -o pdfviewer.exe ./pdfviewer

# Or using PowerShell
go build -o bin\pdfviewer.exe .\pdfviewer
```

### Build Optimized Version (Reduce Size)

```bash
go build -ldflags="-s -w" -o pdfviewer.exe ./pdfviewer
```

## Usage

### Launch Application

```bash
# Method 1: Direct run (will show file selection dialog or double-click blank area)
pdfviewer.exe

# Method 2: Specify PDF file path
pdfviewer.exe document.pdf
```

### Interface Operations

#### Multi-tab (v1.1 New) ⭐
- **Open multiple files** - Menu → File → Open, each open creates new tab
- **Create new tab** - Menu → File → New Tab, creates empty tab
- **Switch tabs** - Click tab title to switch
- **Close tab** - Menu → File → Close Tab, closes current tab
- **Independent operations** - Each tab has independent page navigation and zoom

#### Menu Bar
- **File Menu**
  - Open... - Select PDF file (creates new tab)
  - New Tab - Create empty tab
  - Save As... - Save current document copy to new location
  - Close Tab - Close current tab
  - Exit - Exit program

- **View Menu**
  - First/Previous/Next/Last Page - Page navigation (operates on current tab)
  - Zoom In/Zoom Out/Actual Size - Zoom control (operates on current tab)

- **Help Menu**
  - Shortcuts - View all keyboard shortcuts
  - About - View version and license information

#### Toolbar (Icon-based)
- 📂 Open button - Select PDF file to open (creates new tab)
- 💾 Save As button - Save current document copy to new location
- ❌ Close Tab button - Close current tab (v1.2.2 New)
- ⏮️ First Page - Jump to first page (current tab)
- ◀️ Previous - Go to previous page (current tab)
- ▶️ Next - Go to next page (current tab)
- ⏭️ Last Page - Jump to last page (current tab)
- 🔍➖ Zoom Out - Decrease zoom ratio (current tab)
- 100% - Reset to actual size (current tab)
- 🔍➕ Zoom In - Increase zoom ratio (current tab)
- Page number input - Enter page number and press Enter to jump (current tab)

#### Mouse Operations
- **Wheel page flipping** - Scroll up for previous page, scroll down for next page (current tab)
- **Double-click blank area** - When no document is open, double-click blank area to open file selection dialog

#### Status Bar
Displays detailed document information for currently active tab:
```
document.pdf  |  Page 5 / 120  |  Zoom: 125%  |  Size: 3.2 MB
```

### Keyboard Shortcuts

- `Left Arrow` / `PageUp`: Previous page
- `Right Arrow` / `PageDown` / `Space`: Next page
- `Home`: Jump to first page
- `End`: Jump to last page
- `Ctrl+W`: Close current tab (v1.2.2 New)

### Interface Operations

- **Open button**: Select PDF file to open
- **Navigation buttons**: First page, Previous, Next, Last page
- **Page number input**: Enter page number and press Enter to jump
- **Zoom buttons**: `-` Zoom out, `100%` Reset, `+` Zoom in

## Project Structure

```
/pdfviewer/
├── main.go              # Program entry point
├── ui.go                # User interface (menus, toolbar, event handlers)
├── controller.go        # Controller logic (page management, zoom control)
├── pdf_engine.go        # PDF engine (go-fitz wrapper)
├── theme.go             # Theme configuration
├── icon.go              # Program icon generation
├── i18n.go              # Internationalization support
├── README.md            # Chinese documentation
├── README_EN.md         # English documentation (this file)
```

## Technical Highlights

1. **KISS Principle** - Simple and maintainable code, approximately 1000 lines
2. **Modular Design** - Separation of UI/Controller/Engine
3. **Custom Widgets** - Implements mouse wheel page flipping and double-click events
4. **Code-generated Icons** - No external image resources needed
5. **Cross-platform Support** - Build once for Windows/Linux/macOS
6. **Multi-language Support** - English and Chinese UI (v1.3 New)

## Dependencies

- **Fyne**: GUI framework (v2.4+)
- **go-fitz**: Go wrapper for MuPDF, used for PDF rendering (AGPL license)

## Common Issues

### Q: Build error "gcc: command not found"
A: Need to install GCC compiler and ensure it's added to system PATH.

### Q: Linker error during build
A: Ensure you're using 64-bit version of GCC and that the version is compatible with Go version.

### Q: Failed to open PDF file
A: Ensure PDF file is not corrupted and is not an encrypted PDF.

### Q: Slow page rendering
A: You can reduce DPI setting (modify baseDPI value in controller.go).

## License

- go-fitz: AGPL-3.0
- Fyne: BSD-3-Clause

## Contact

For issues, please refer to official documentation for Fyne and go-fitz:
- Fyne: https://fyne.io/
- go-fitz: https://github.com/gen2brain/go-fitz
