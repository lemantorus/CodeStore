package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Конфигурация фильтров ---
var blacklist = map[string]bool{
	"node_modules": true, ".git": true, "venv": true, ".venv": true,
	"target": true, "dist": true, "build": true, "vendor": true,
}

var ignoreExt = map[string]bool{
	".exe": true, ".dll": true, ".so": true, ".png": true, ".jpg": true,
	".jpeg": true, ".gif": true, ".pdf": true, ".zip": true, ".pyc": true,
	".ico": true, ".ttf": true, ".woff": true, ".woff2": true,
}

// --- Улучшенные Стили ---
var (
	// Папки в обычном списке
	dirStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)
	// Файлы в обычном списке
	fileStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA"))
	// Стиль выбранной строки (Белый текст на Фиолетовом фоне)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#7D56F4")).Bold(true)
	// Подсказки клавиш (Светло-серый)
	tipStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Italic(true)
	
	headerStyle   = lipgloss.NewStyle().Background(lipgloss.Color("#FF5F87")).Foreground(lipgloss.Color("#FFFFFF")).Padding(0, 1).Bold(true)
	searchStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
)

type entry struct {
	name  string
	isDir bool
}

type model struct {
	path            string
	entries         []entry
	filteredEntries []entry
	cursor          int
	scrollOffset    int
	searching       bool
	searchQuery     string
	terminalMsg     string
}

func initialModel() model {
	path, _ := os.Getwd()
	m := model{path: path}
	m.updateEntries()
	return m
}

func (m *model) updateEntries() {
	files, _ := os.ReadDir(m.path)
	m.entries = []entry{{name: "..", isDir: true}}
	for _, f := range files {
		m.entries = append(m.entries, entry{name: f.Name(), isDir: f.IsDir()})
	}
	m.applyFilter()
}

func (m *model) applyFilter() {
	if m.searchQuery == "" {
		m.filteredEntries = m.entries
	} else {
		m.filteredEntries = []entry{}
		for _, e := range m.entries {
			if strings.Contains(strings.ToLower(e.name), strings.ToLower(m.searchQuery)) || e.name == ".." {
				m.filteredEntries = append(m.filteredEntries, e)
			}
		}
	}
	m.cursor = 0
	m.scrollOffset = 0
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searching {
			switch msg.String() {
			case "esc", "enter":
				m.searching = false
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.applyFilter()
				}
			case "ctrl+c":
				return m, tea.Quit
			default:
				if len(msg.String()) == 1 {
					m.searchQuery += msg.String()
					m.applyFilter()
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+w":
			m.searching = true
		case "esc":
			m.searchQuery = ""
			m.applyFilter()
			m.terminalMsg = ""
		case "up", "k":
			if m.cursor > 0 { m.cursor-- }
		case "down", "j":
			if m.cursor < len(m.filteredEntries)-1 { m.cursor++ }
		case "enter":
			if len(m.filteredEntries) == 0 { return m, nil }
			selected := m.filteredEntries[m.cursor]
			if selected.isDir {
				m.path = filepath.Join(m.path, selected.name)
				m.path = filepath.Clean(m.path)
				m.searchQuery = ""
				m.updateEntries()
			}
		case "r", "R":
			filename := m.collectCode()
			openFolder(m.path)
			m.terminalMsg = successStyle.Render("✅ Created: " + filename)
		}
	}

	height := 18
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	} else if m.cursor >= m.scrollOffset+height {
		m.scrollOffset = m.cursor - height + 1
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder
	s.WriteString(headerStyle.Render(" CODE COLLECTOR ") + " " + m.path + "\n")
	
	if m.searching {
		s.WriteString(searchStyle.Render("🔍 Search: ") + m.searchQuery + "█\n")
	} else if m.searchQuery != "" {
		s.WriteString(searchStyle.Render("🔍 Filter: ") + m.searchQuery + tipStyle.Render(" [ESC to clear]") + "\n")
	} else {
		s.WriteString(tipStyle.Render("CTRL+W: Search | R: Collect & Open | Enter: Open Dir") + "\n")
	}
	s.WriteString("\n")

	viewHeight := 18
	end := m.scrollOffset + viewHeight
	if end > len(m.filteredEntries) { end = len(m.filteredEntries) }

	visible := m.filteredEntries[m.scrollOffset:end]
	for i, entry := range visible {
		absIndex := i + m.scrollOffset
		pointer := "  "
		if m.cursor == absIndex { pointer = "> " }

		var line string
		if m.cursor == absIndex {
			// При выборе - используем единый стиль выделения для папок и файлов
			name := entry.name
			if entry.isDir { name += "/" }
			line = selectedStyle.Render(pointer + name)
		} else {
			// Обычный вид
			if entry.isDir {
				line = pointer + dirStyle.Render(entry.name+"/")
			} else {
				line = pointer + fileStyle.Render(entry.name)
			}
		}
		s.WriteString(line + "\n")
	}

	if m.terminalMsg != "" {
		s.WriteString("\n" + m.terminalMsg + "\n")
	}

	return s.String()
}

func (m model) collectCode() string {
	timestamp := time.Now().Format("20060102_150405")
	folderName := filepath.Base(m.path)
	fileName := fmt.Sprintf("dump_%s_%s.txt", folderName, timestamp)
	dumpPath := filepath.Join(m.path, fileName)

	f, err := os.Create(dumpPath)
	if err != nil { return "error" }
	defer f.Close()

	fmt.Fprintf(f, "================================================\n")
	fmt.Fprintf(f, "📁 PROJECT: %s\n", folderName)
	fmt.Fprintf(f, "📍 PATH: %s\n", m.path)
	fmt.Fprintf(f, "🕒 COLLECTED: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f, "================================================\n\n")

	filepath.WalkDir(m.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			if d != nil && d.IsDir() && blacklist[d.Name()] { return filepath.SkipDir }
			return nil
		}
		if strings.HasPrefix(d.Name(), "dump_") && strings.HasSuffix(d.Name(), ".txt") { return nil }
		ext := filepath.Ext(d.Name())
		if ignoreExt[ext] || strings.HasPrefix(d.Name(), ".") { return nil }

		content, _ := os.ReadFile(path)
		rel, _ := filepath.Rel(m.path, path)
		fmt.Fprintf(f, "####\n%s\n%s\n#####\n\n", rel, string(content))
		return nil
	})

	return fileName
}

func openFolder(path string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows": cmd = exec.Command("explorer", path)
	case "darwin": cmd = exec.Command("open", path)
	default: cmd = exec.Command("xdg-open", path)
	}
	cmd.Run()
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
