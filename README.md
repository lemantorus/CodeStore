# CodeStore 🚀

A powerful code storage and management utility built with Go! 🎉

## Overview

CodeStore is a versatile tool designed to help you manage and store your code snippets efficiently. Whether you're a developer looking to organize your code or need a quick way to dump and retrieve code, CodeStore has got you covered!

## Features

- 📦 **Efficient Code Storage**: Store and manage your code snippets with ease
- ⚡ **Fast Retrieval**: Quickly access your stored code when you need it
- 🔍 **Search Functionality**: Find exactly what you're looking for
- 📁 **Organized Structure**: Clean and intuitive organization of your code
- 🛠 **Go-Powered**: Built with the speed and reliability of Go

## Screenshots

Check out CodeStore in action! 📸

### Search Function
![Search Function](images/search_func.png)

### Working Results
![Working Results](images/Working_result.png)

## Installation

1. Make sure you have Go installed on your system
2. Clone this repository:
   ```bash
   git clone https://github.com/lemantorus/CodeStore.git
   ```
3. Navigate to the project directory:
   ```bash
   cd CodeStore
   ```
4. Run the application:
   ```bash
   go run main.go
   ```

## Usage

### Direct Execution
Simply run the application directly:
```bash
go run main.go
```

### Install Globally
To make CodeStore available from anywhere on your system:

1. Build the binary:
   ```bash
   go build -o codestore main.go
   ```

2. Move the binary to a directory in your PATH:
   ```bash
   # Option 1: Move to /usr/local/bin (requires sudo)
   sudo mv codestore /usr/local/bin/

   # Option 2: Move to $HOME/bin (doesn't require sudo, but ensure $HOME/bin is in your PATH)
   mv codestore $HOME/bin/
   ```

3. Now you can run CodeStore from anywhere:
   ```bash
   codestore
   ```

### Build and Install with Go
Alternatively, you can install it directly using Go:
```bash
go install .
```
This will build and install the binary to your `$GOBIN` directory (or `$GOPATH/bin` if `$GOBIN` is not set), which should be in your PATH.

## How It Works

CodeStore is built as an interactive terminal application using the Bubble Tea framework, which provides a clean and intuitive user interface for navigating and managing your code files.

### Core Functionality

1. **Directory Navigation**: The application starts in your current working directory and allows you to navigate through your file system using arrow keys or vim-style navigation (j/k for up/down).

2. **Code Collection**: When you press 'R', CodeStore performs the following operations:
   - Creates a timestamped dump file (e.g., `dump_projectname_20260102_153045.txt`)
   - Walks through all files in the current directory and subdirectories
   - Filters out unwanted files and directories based on predefined configurations
   - Combines all relevant code files into a single text file with clear separators

3. **Search & Filter**: Press `Ctrl+W` to enter search mode, allowing you to filter files and directories by name in real-time.

4. **File System Filtering**: The application intelligently skips certain directories and file types to avoid collecting unnecessary data:
   - Blacklisted directories: `node_modules`, `.git`, `venv`, `.venv`, `target`, `dist`, `build`, `vendor`
   - Ignored file extensions: `.exe`, `.dll`, `.so`, `.png`, `.jpg`, `.jpeg`, `.gif`, `.pdf`, `.zip`, `.pyc`, `.ico`, `.ttf`, `.woff`, `.woff2`

### Internal Architecture

The application is structured around a Bubble Tea model with the following key components:

- **Model**: Manages the application state including current path, file entries, cursor position, and search functionality
- **View**: Renders the terminal UI with color-coded directories and files, search indicators, and status messages
- **Update**: Handles user input events (navigation, search, file operations) and updates the application state accordingly

### UI Components and Navigation

CodeStore features an intuitive terminal-based user interface with the following elements:

**Header Section**:
- Displays the "CODE COLLECTOR" title with a pink background
- Shows the current directory path being browsed

**Navigation Controls**:
- **Arrow Keys** / **J/K**: Move cursor up and down through the file list
- **Enter**: Open selected directory or navigate into it
- **R**: Collect all code in the current directory into a dump file
- **Ctrl+W**: Enter search/filter mode
- **ESC**: Clear search filter or exit search mode
- **Ctrl+C** / **Q**: Quit the application

**Visual Elements**:
- **Directories**: Displayed in bold purple (`#7D56F4`)
- **Files**: Displayed in white (`#FAFAFA`)
- **Selected Item**: Highlighted with white text on purple background
- **Search Indicator**: Shows current search term with a yellow "🔍" icon
- **Status Messages**: Success messages appear in green when code collection is complete

**Scrolling Behavior**: The interface automatically manages scrolling when there are more items than can fit in the visible area (18 lines), keeping the cursor in view as you navigate.

### Code Collection Process

When collecting code, the application:

1. Creates a new dump file with a timestamp in the current directory
2. Adds metadata headers including project name, path, and collection time
3. Recursively walks through all subdirectories using `filepath.WalkDir`
4. Applies filtering rules to exclude unwanted files and directories
5. Reads content from each valid file and appends it to the dump file with clear path separators
6. Automatically opens the containing folder after collection is complete

#### Filtering Mechanism

The filtering system operates at multiple levels:

- **Directory Blacklisting**: Directories like `node_modules`, `.git`, and `vendor` are completely skipped during the file walk to avoid collecting unnecessary files
- **File Extension Filtering**: Binary files and non-code files (images, executables, archives) are excluded based on their extensions
- **Hidden File Filtering**: Files starting with a dot (`.`) are excluded to avoid collecting configuration files that aren't part of the main codebase
- **Self-Exclusion**: The application automatically excludes any previously generated dump files to prevent recursive inclusion

#### Configuration Details

CodeStore uses two main configuration maps to control what gets collected:

**Blacklisted Directories (`blacklist`)**:
```
"node_modules": true,  // Node.js dependency directory
".git": true,          // Git version control directory
"venv": true,          // Python virtual environment
".venv": true,         // Alternative Python virtual environment
"target": true,        // Rust/Cargo build directory
"dist": true,          // Distribution build directory
"build": true,         // Build output directory
"vendor": true         // Dependency vendor directory
```

**Ignored File Extensions (`ignoreExt`)**:
```
".exe": true,    // Executable files
".dll": true,    // Dynamic link libraries
".so": true,     // Shared object files (Linux)
".png": true,    // Image files
".jpg": true,    // Image files
".jpeg": true,   // Image files
".gif": true,    // Image files
".pdf": true,    // PDF documents
".zip": true,    // Archive files
".pyc": true,    // Python compiled files
".ico": true,    // Icon files
".ttf": true,    // Font files
".woff": true,   // Web font files
".woff2": true   // Web font files
```

The collected code is formatted with clear section markers (`####`) that include the relative file path, making it easy to identify where each code segment originated.

## Tech Stack

- **Language**: Go (Golang)
- **Architecture**: Terminal user interface (TUI) with Bubble Tea framework
- **UI Styling**: Lipgloss for beautiful terminal formatting
- **Dependencies**: Managed via go.mod

## Technical Architecture

CodeStore follows a Model-View-Update (MVU) architecture pattern, which is central to the Bubble Tea framework:

**Model (`model` struct)**:
- Stores application state (current path, file entries, cursor position, search query)
- Manages the file system navigation state
- Handles filtering and search functionality

**View (`View()` method)**:
- Renders the terminal interface with proper styling
- Displays directory and file listings with color coding
- Shows search status and user prompts
- Handles layout and formatting using Lipgloss

**Update (`Update()` method)**:
- Processes user input events (keyboard commands)
- Updates application state based on user actions
- Implements the core business logic for navigation and file operations

**Key Technical Features**:
- Cross-platform file system operations using Go's `filepath` package
- Recursive directory traversal with intelligent filtering
- Platform-specific folder opening using `os/exec` (Explorer on Windows, `open` on macOS, `xdg-open` on Linux)
- Real-time search and filtering with immediate UI updates
- Scrollable interface with automatic offset management for large directories

## Contributing

We welcome contributions! 🤝

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with ❤️ using Go
- Inspired by the need for better code organization
- Thanks to all contributors who make this project better!

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

⭐ **Star this repo if you find it useful!** ⭐
