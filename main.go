package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Note struct represents individual notes with their metadata
type note struct {
	title     string  // Title of the note
	path      string  // File path of the note
	createdAt int64   // Timestamp of note creation
}

// Implement list.Item interface methods for seamless list integration
func (n note) Title() string       { return n.title }
func (n note) Description() string { return time.Unix(n.createdAt, 0).Format("2006-01-02 15:04:05") }
func (n note) FilterValue() string { return n.title }

// Model defines the entire application state
type model struct {
	list          list.Model      // Notes list view
	textInput     textinput.Model // Input for note titles
	textarea      textarea.Model  // Content editing area
	notes         []note          // Slice of all notes
	mode          string          // Current application mode (list/new/edit)
	selectedNote  *note           // Currently selected note
	width, height int             // Window dimensions
	titleEntered  bool            // Tracks title input state
}

// Define application-wide styling for consistent UI
var (
	// Directory to store notes
	notesDir = filepath.Join(os.Getenv("HOME"), ".notes")

	// Document container style
	docStyle = lipgloss.NewStyle().Padding(1, 2)

	// Split view style with rounded borders
	splitStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 1)

	// Help text styling
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Width(80).
			MarginTop(1).
			MarginBottom(0).
			PaddingLeft(2).
			PaddingRight(2)

	// Title styling
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	// Content styling
	contentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			MarginTop(1)
)

// Help text provides quick reference for user interactions
const helpText = `Navigation: ↑/↓:Navigate | enter:View | esc:Back | ctrl+n:New | ctrl+s:Save | ctrl+e:Edit | ctrl+d:Delete | ctrl+u:Refresh | ctrl+c:Quit`

// initialModel sets up the initial application state
func initialModel() model {
	// Create text input for note titles
	ti := textinput.New()
	ti.Placeholder = "Note title (Press Tab to enter content)"
	ti.CharLimit = 50
	ti.Focus()

	// Create text area for note content
	ta := textarea.New()
	ta.Placeholder = "Enter note content (Ctrl+S to save)..."
	ta.ShowLineNumbers = false
	ta.Prompt = "┃ "

	// Configure list with a custom delegate
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true  // Show creation timestamps

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Notes"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	return model{
		list:      l,
		textInput: ti,
		textarea:  ta,
		mode:      "list",
	}
}

// Init prepares initial commands when the application starts
func (m model) Init() tea.Cmd {
	return tea.Batch(
		loadNotes,  // Load existing notes
		textarea.Blink,  // Enable text area cursor blinking
	)
}

// Update handles all application state changes and user interactions
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Adjust UI components based on window size
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width/2-4, msg.Height-10)
		m.textarea.SetWidth(msg.Width/2 - 4)
		m.textarea.SetHeight(msg.Height - 12)

	case tea.KeyMsg:
		switch {
		// Quit application
		case msg.Type == tea.KeyCtrlC:
			return m, tea.Quit

		// Refresh notes list
		case msg.String() == "ctrl+u":
			return m, loadNotes

		// Switch from title input to content input for both new and edit modes
		case msg.Type == tea.KeyTab && (m.mode == "new" || m.mode == "edit") && m.textInput.Focused():
			m.titleEntered = true
			m.textInput.Blur()
			m.textarea.Focus()
			return m, nil

		// Enter new note mode
		case msg.Type == tea.KeyCtrlN:
			m.mode = "new"
			m.textInput.Reset()
			m.textarea.Reset()
			m.titleEntered = false
			m.textInput.Focus()
			m.selectedNote = nil

		// Save note (new or edited)
		case msg.Type == tea.KeyCtrlS && (m.mode == "new" || m.mode == "edit"):
			if m.textInput.Value() != "" {
				cmd = saveNote(m.textInput.Value(), m.textarea.Value(), m.selectedNote)
				m.mode = "list"
				m.textInput.Reset()
				m.textarea.Reset()
				m.titleEntered = false
				m.selectedNote = nil
				return m, tea.Batch(cmd, loadNotes)
			}

		// Delete selected note
		case msg.Type == tea.KeyCtrlD && m.selectedNote != nil:
			return m, tea.Batch(deleteNote(m.selectedNote.path), loadNotes)

		// Edit selected note
		case msg.Type == tea.KeyCtrlE && m.selectedNote != nil:
			m.mode = "edit"
			m.textInput.SetValue(m.selectedNote.title)
			content, _ := os.ReadFile(m.selectedNote.path)
			m.textarea.SetValue(string(content))
			m.textInput.Focus()
			m.titleEntered = true

		// Enhanced list navigation
		case (msg.Type == tea.KeyUp || msg.Type == tea.KeyDown) && m.mode == "list":
			m.list, cmd = m.list.Update(msg)
			
			// Update selected note content immediately
			if selected := m.list.SelectedItem(); selected != nil {
				currentNote := selected.(note)
				m.selectedNote = &currentNote
				
				content, err := os.ReadFile(currentNote.path)
				if err == nil {
					m.textarea.SetValue(string(content))
				}
			}
			
			return m, cmd

		// View note details
		case msg.Type == tea.KeyEnter && m.mode == "list":
			if selected := m.list.SelectedItem(); selected != nil {
				note := selected.(note)
				m.selectedNote = &note
				content, _ := os.ReadFile(note.path)
				m.textarea.SetValue(string(content))
			}

		// Return to list mode
		case msg.Type == tea.KeyEsc:
			m.mode = "list"
			m.textInput.Reset()
			m.textarea.Reset()
			m.titleEntered = false
			m.textInput.Blur()
			m.textarea.Blur()
			m.selectedNote = nil
		}

	// Handle notes loading
	case []note:
		// Sort notes by creation time (newest first)
		sort.Slice(msg, func(i, j int) bool {
			return msg[i].createdAt > msg[j].createdAt
		})
		m.notes = msg
		m.list.SetItems(itemsFromNotes(msg))

		// Select first note if available
		if len(msg) > 0 {
			if m.selectedNote == nil {
				m.list.Select(0)
				m.selectedNote = &msg[0]
				content, _ := os.ReadFile(msg[0].path)
				m.textarea.SetValue(string(content))
			} else {
				// Try to maintain previous note selection
				found := false
				for i, n := range msg {
					if n.path == m.selectedNote.path {
						m.list.Select(i)
						m.selectedNote = &n
						content, _ := os.ReadFile(n.path)
						m.textarea.SetValue(string(content))
						found = true
						break
					}
				}
				if !found {
					m.list.Select(0)
					m.selectedNote = &msg[0]
					content, _ := os.ReadFile(msg[0].path)
					m.textarea.SetValue(string(content))
				}
			}
		}
	}

	// Update input components based on current mode
	if m.mode == "new" || m.mode == "edit" {
		if m.textInput.Focused() {
			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			m.textarea, cmd = m.textarea.Update(msg)
			cmds = append(cmds, cmd)
		}
	} else {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the entire application UI
func (m model) View() string {
	// Create list view
	listView := splitStyle.
		Width(m.width/2 - 36).
		Height(m.height - 6).
		Render(m.list.View())

	// Create content view
	var contentView string
	if m.mode == "new" || m.mode == "edit" {
		contentView = splitStyle.Width(m.width/2 +30).Render(
			lipgloss.JoinVertical(lipgloss.Top,
				titleStyle.Render(m.textInput.View()),
				contentStyle.Render(m.textarea.View()),
			),
		)
	} else {
		contentView = splitStyle.
			Width(m.width/2 +30).
			Height(m.height - 6).
			Render(contentStyle.Render(m.textarea.View()))
	}

	// Render help text
	helpView := helpStyle.Render(helpText)
	
	// Combine all views
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, listView, contentView)
	return docStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top, mainView, helpView),
	)
}

// Main application entry point
func main() {
	// Ensure notes directory exists
	if _, err := os.Stat(notesDir); os.IsNotExist(err) {
		os.Mkdir(notesDir, 0755)
	}

	// Start the Bubble Tea program
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

// Convert notes to list items for display
func itemsFromNotes(notes []note) []list.Item {
	items := make([]list.Item, len(notes))
	for i, n := range notes {
		items[i] = n
	}
	return items
}

// Load notes from the notes directory
func loadNotes() tea.Msg {
	files, _ := os.ReadDir(notesDir)
	var notes []note

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".md" {
			nameParts := strings.SplitN(f.Name(), "-", 2)
			if len(nameParts) < 2 {
				continue
			}

			timestamp, err := strconv.ParseInt(nameParts[0], 10, 64)
			if err != nil {
				continue
			}

			cleanName := strings.TrimSuffix(nameParts[1], ".md")
			cleanName = strings.ReplaceAll(cleanName, "-", " ")
			notes = append(notes, note{
				title:     cleanName,
				path:      filepath.Join(notesDir, f.Name()),
				createdAt: timestamp,
			})
		}
	}
	return notes
}

// Save a note, preserving original timestamp for existing notes
func saveNote(title, content string, existingNote *note) tea.Cmd {
	return func() tea.Msg {
		sanitized := sanitizeFileName(title)
		var path string

		if existingNote != nil {
			// Preserve the original creation timestamp
			filenameParts := strings.SplitN(filepath.Base(existingNote.path), "-", 2)
			originalTimestamp := filenameParts[0]
			
			path = filepath.Join(notesDir, fmt.Sprintf("%s-%s.md", originalTimestamp, sanitized))
			os.Remove(existingNote.path)
		} else {
			path = filepath.Join(notesDir, fmt.Sprintf("%d-%s.md", time.Now().Unix(), sanitized))
		}

		// Directly save the full content
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			fmt.Printf("Error saving note: %v", err)
		}
		return loadNotes()
	}
}

// Delete a note from the filesystem
func deleteNote(path string) tea.Cmd {
	return func() tea.Msg {
		os.Remove(path)
		return loadNotes()
	}
}

// Sanitize filename to remove invalid characters
func sanitizeFileName(input string) string {
	name := strings.TrimSuffix(input, ".md")
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' || r == '_' {
			return r
		}
		return '-'
	}, name)
}