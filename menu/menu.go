package menu

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// Main struct to handle options for Command, Description, and the
// function that should be called
type CommandOption struct {
	Command, Description string
	Function             func(args ...string) error
}

// Menu options -- right now only sets prompt
type MenuOptions struct {
	Prompt     string
	MenuLength int
}

// Menu struct encapsulates Commands and Options
type Menu struct {
	Commands []CommandOption
	Options  MenuOptions
}

// Setup the options for the menu.
//
// An empty string for prompt and a length of 0 will use the
// default "> " prompt and 100 character wide menu
func NewMenuOptions(prompt string, length int) MenuOptions {
	if prompt == "" {
		prompt = "> "
	}

	if length == 0 {
		length = 100
	}

	return MenuOptions{prompt, length}
}

// Trim whitespace, newlines, and create command+arguments slice
func cleanCommand(cmd string) ([]string, error) {
	cmd_args := strings.Split(strings.Trim(cmd, " \n"), " ")
	return cmd_args, nil
}

// Creates a new menu with options
func NewMenu(cmds []CommandOption, options MenuOptions) *Menu {
	return &Menu{cmds, options}
}

func (m *Menu) prompt() {
	fmt.Print(m.Options.Prompt)
}

// Write menu from CommandOptions with tabwriter
func (m *Menu) menu() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 1, ' ', 0)
	layoutMenu(w, m.Commands, m.Options.MenuLength)
}

// Wrapper for providing Stdin to the main menu loop
func (m *Menu) Start() {
	m.start(os.Stdin)
}

// Main loop
func (m *Menu) start(reader io.Reader) {
	m.menu()
MainLoop:
	for {
		input := bufio.NewReader(reader)
		// Prompt for input
		m.prompt()

		inputString, err := input.ReadString('\n')
		if err != nil {
			// If we didn't receive anything from ReadString
			// we shouldn't continue because we're not blocking
			// anymore but we also don't have any data
			break MainLoop
		}

		// convert CRLF to LF
		inputString = strings.Replace(inputString, "\r\n", "", -1)

		cmd, _ := cleanCommand(inputString)
		if len(cmd) < 1 {
			break MainLoop
		}
		// Route the first index of the cmd slice to the appropriate case
	Route:
		switch cmd[0] {
		case "exit", "quit":
			fmt.Println("Exiting...")
			break MainLoop

		case "menu":
			m.menu()
			break

		default:
			// Loop through commands and find the right one
			// Probably a more efficient way to do this, but unless we have
			// tons of commands, it probably doesn't matter
			for i := range m.Commands {
				if m.Commands[i].Command == cmd[0] {
					err := m.Commands[i].Function(cmd[1:]...)
					if err != nil {
						panic(err)
					}

					break Route
				}
			}
			// Shouldn't get here if we found a command
			fmt.Println("Unknown command")
		}
	}
}
