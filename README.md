# Gleaner 📝

Gleaner is a lightweight, terminal-based note-taking application built with Go and the Bubble Tea framework. It provides a minimalistic and efficient way to create, manage, and organize your notes directly from the command line.


## Screenshots of the Terminal based App
### loading screen
<img width="1366" height="768" alt="250711_17h31m51s_screenshot" src="https://github.com/user-attachments/assets/f1715ae4-ee0a-42fd-9d95-9c9dc85ba978" />

### editing with notes
<img width="1366" height="768" alt="250711_17h34m21s_screenshot" src="https://github.com/user-attachments/assets/7e5e1861-16e6-4854-990a-d66bd455fb6d" />

### filtering notes
<img width="1366" height="768" alt="250711_17h35m37s_screenshot" src="https://github.com/user-attachments/assets/f40a8589-fdc4-42fe-9669-55e5d621d047" />




## 🌟 Features

- **Create Notes**: Quickly create new notes with a simple interface
- **Edit Notes**: Modify existing notes with ease
- **Delete Notes**: Remove notes you no longer need
- **List View**: Browse through your notes with a clean, organized list
- **Timestamp Tracking**: Automatically tracks note creation times
- **Keyboard-Driven**: Navigate and manage notes using keyboard shortcuts

## 🛠 Prerequisites

- Go 1.23.5 + 
- Git

## 🚀 Installation

1. Clone the repository:
```bash
git clone https://github.com/debjit-1004/gleaner.git
cd gleaner
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build
```

4. Run Gleaner:
```bash
./gleaner
```

## 📋 Keyboard Shortcuts

- `Ctrl+N`: Create a new note
- `Ctrl+E`: Edit selected note
- `Ctrl+S`: Save note
- `Ctrl+D`: Delete selected note
- `Ctrl+U`: Refresh notes list
- `Tab`: Switch between title and content fields
- `Esc`: Return to list view
- `↑/↓`: Navigate notes
- `Enter`: View note details
- `Ctrl+C`: Quit application

## 📦 Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea): Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss): Terminal styling
- Standard Go libraries for file management and time handling

## 🔧 Note Storage

Notes are stored as Markdown files in `~/.notes` directory. Each note filename includes a timestamp for unique identification and chronological sorting.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

## 🙌 Acknowledgments

- [Charmbracelet](https://github.com/charmbracelet) for amazing Go terminal UI libraries
- Go community for incredible development tools

## 📞 Contact

Debjit - [GitHub Profile](https://github.com/debjit-1004)

Project Link: [https://github.com/debjit-1004/gleaner](https://github.com/debjit-1004/gleaner)
