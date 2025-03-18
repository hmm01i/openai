package commands

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CommandResponse represents a structured response from a command
type CommandResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// Command represents a command or subcommand
type Command struct {
	Execute   func(c ChatClient, args []string) string
	Help      string
	SubCmds   map[string]*Command
	MinAccess AccessLevel // Minimum access level required for this command
}

// AccessLevel represents the feature flag level for commands
type AccessLevel int

const (
	// AccessLegacy represents the old command system
	AccessLegacy AccessLevel = iota
	// AccessBeta represents the new command system with subcommands
	AccessBeta
	// AccessStable represents fully tested and stable commands
	AccessStable
)

// ChatClient interface defines the methods that the chat client must implement
type ChatClient interface {
	ListPersonas() []string
	SavePersona(name, directive string) error
	ShowPersona() string
	LoadPersona(name string) error
	SetDirective(directive string) error
	ClearHistory()
	ListModels() []string
	SetModel(model string)
	SaveConversation(name string) error
	ListConversations() []string
	LoadConversation(name string) error
	GetCurrentPersona() string
	GetHistory() []Message
}

// Message represents a chat message
type Message struct {
	Role    string
	Content string
}

// CommandRegistry manages the available commands and their access levels
type CommandRegistry struct {
	commands    map[string]*Command
	accessLevel AccessLevel
}

// NewCommandRegistry creates a new command registry with the specified access level
func NewCommandRegistry(accessLevel AccessLevel) *CommandRegistry {
	r := &CommandRegistry{
		commands:    make(map[string]*Command),
		accessLevel: accessLevel,
	}
	r.registerCommands()
	return r
}

func formatResponse(success bool, message string, err error) string {
	resp := CommandResponse{
		Success: success,
		Message: message,
	}
	if err != nil {
		resp.Error = err.Error()
	}
	jsonResp, _ := json.Marshal(resp)
	return string(jsonResp)
}

// ExecuteCommand handles command execution with subcommand support
func (r *CommandRegistry) ExecuteCommand(c ChatClient, input string) string {
	if !strings.HasPrefix(input, "/") {
		return ""
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return ""
	}

	cmd, exists := r.commands[parts[0]]
	if !exists {
		return formatResponse(false, "", fmt.Errorf("unknown command: %s (try /help)", parts[0]))
	}

	// Check if the command is available at current access level
	if cmd.MinAccess > r.accessLevel {
		return formatResponse(false, "", fmt.Errorf("command %s is not available in current mode", parts[0]))
	}

	if cmd.SubCmds != nil && r.accessLevel >= AccessBeta {
		return r.executeSubCommand(c, cmd, parts[1:])
	}

	return cmd.Execute(c, parts[1:])
}

func (r *CommandRegistry) executeSubCommand(c ChatClient, cmd *Command, args []string) string {
	if len(args) == 0 {
		return formatResponse(false, "", fmt.Errorf("missing subcommand\n%s", cmd.Help))
	}

	subCmd, exists := cmd.SubCmds[args[0]]
	if !exists {
		return formatResponse(false, "", fmt.Errorf("unknown subcommand: %s\n%s", args[0], cmd.Help))
	}

	if subCmd.MinAccess > r.accessLevel {
		return formatResponse(false, "", fmt.Errorf("subcommand %s is not available in current mode", args[0]))
	}

	return subCmd.Execute(c, args[1:])
}

// GetHelp returns help information for commands
func (r *CommandRegistry) GetHelp(command string) string {
	if command == "" {
		var help []string
		help = append(help, "Available commands:")
		for cmdName, cmd := range r.commands {
			if cmd.MinAccess <= r.accessLevel {
				help = append(help, fmt.Sprintf("%s - %s", cmdName, strings.Split(cmd.Help, "\n")[0]))
			}
		}
		return formatResponse(true, strings.Join(help, "\n"), nil)
	}

	cmd, exists := r.commands[command]
	if !exists {
		return formatResponse(false, "", fmt.Errorf("no help available for: %s", command))
	}

	if cmd.MinAccess > r.accessLevel {
		return formatResponse(false, "", fmt.Errorf("command %s is not available in current mode", command))
	}

	return formatResponse(true, cmd.Help, nil)
}

// addHelpSubCommand adds a help subcommand to a command
func addHelpSubCommand(cmd *Command) {
	if cmd.SubCmds == nil {
		cmd.SubCmds = make(map[string]*Command)
	}
	cmd.SubCmds["help"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			return formatResponse(true, cmd.Help, nil)
		},
		Help:      "Show help for this command",
		MinAccess: cmd.MinAccess,
	}
}

func (r *CommandRegistry) registerCommands() {
	// Legacy commands (always available)
	r.commands["/q"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			return formatResponse(true, "Goodbye!", nil)
		},
		Help:      "Quit the application",
		MinAccess: AccessLegacy,
	}

	r.commands["/help"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			if len(args) > 0 {
				return r.GetHelp(args[0])
			}
			return r.GetHelp("")
		},
		Help:      "Show help information. Use /help <command> for detailed help on a command.",
		MinAccess: AccessLegacy,
	}

	// Beta commands (new command system)
	r.commands["/system"] = &Command{
		Help: `System Commands:
  directive <text> - Set the system directive
  help            - Show this help message`,
		MinAccess: AccessBeta,
		SubCmds:   make(map[string]*Command),
	}
	addHelpSubCommand(r.commands["/system"])

	r.commands["/history"] = &Command{
		Help: `History Commands:
  show   - Show conversation history
  clear  - Clear conversation history
  help   - Show this help message`,
		MinAccess: AccessBeta,
		SubCmds:   make(map[string]*Command),
	}
	addHelpSubCommand(r.commands["/history"])

	r.commands["/persona"] = &Command{
		Help: `Persona Management Commands:
  list         - List all personas (* marks current)
  show         - Show current persona
  save <name>  - Save current system directive as a persona
  load <name>  - Load a persona by name
  help         - Show this help message`,
		MinAccess: AccessBeta,
		SubCmds:   make(map[string]*Command),
	}
	addHelpSubCommand(r.commands["/persona"])

	r.commands["/model"] = &Command{
		Help: `Model Commands:
  list       - List available models
  set <name> - Set the current model
  help       - Show this help message`,
		MinAccess: AccessBeta,
		SubCmds:   make(map[string]*Command),
	}
	addHelpSubCommand(r.commands["/model"])

	r.commands["/conversation"] = &Command{
		Help: `Conversation Commands:
  list       - List saved conversations
  save <name> - Save current conversation
  load <name> - Load a saved conversation
  help       - Show this help message`,
		MinAccess: AccessBeta,
		SubCmds:   make(map[string]*Command),
	}
	addHelpSubCommand(r.commands["/conversation"])

	// Add all the subcommands after help is added
	addPersonaCommands(r.commands["/persona"])
	addSystemCommands(r.commands["/system"])
	addHistoryCommands(r.commands["/history"])
	addModelCommands(r.commands["/model"])
	addConversationCommands(r.commands["/conversation"])
}

func addPersonaCommands(cmd *Command) {
	cmd.SubCmds["list"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			personas := c.ListPersonas()
			currentPersona := c.GetCurrentPersona()
			for i, p := range personas {
				if p == currentPersona {
					personas[i] = p + "*"
				}
			}
			return formatResponse(true, strings.Join(personas, "\n"), nil)
		},
		Help:      "Lists all available personas",
		MinAccess: AccessBeta,
	}
	cmd.SubCmds["show"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			persona := c.ShowPersona()
			return formatResponse(true, persona, nil)
		},
		Help:      "Shows the current persona",
		MinAccess: AccessBeta,
	}
	cmd.SubCmds["save"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			if len(args) < 1 {
				return formatResponse(false, "", fmt.Errorf("usage: /persona save <name>"))
			}
			if err := c.SavePersona(args[0], c.ShowPersona()); err != nil {
				return formatResponse(false, "", fmt.Errorf("failed to save persona: %w", err))
			}
			return formatResponse(true, "Persona saved successfully", nil)
		},
		Help:      "Saves the current system directive as a persona",
		MinAccess: AccessBeta,
	}
	cmd.SubCmds["load"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			if len(args) < 1 {
				return formatResponse(false, "", fmt.Errorf("usage: /persona load <name>"))
			}
			if err := c.LoadPersona(args[0]); err != nil {
				return formatResponse(false, "", fmt.Errorf("failed to load persona: %w", err))
			}
			return formatResponse(true, "Persona loaded successfully", nil)
		},
		Help:      "Loads a persona by name",
		MinAccess: AccessBeta,
	}
}

func addSystemCommands(cmd *Command) {
	cmd.SubCmds["directive"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			if len(args) < 1 {
				return formatResponse(false, "", fmt.Errorf("usage: /system directive <text>"))
			}
			directive := strings.Join(args, " ")
			if err := c.SetDirective(directive); err != nil {
				return formatResponse(false, "", fmt.Errorf("failed to set directive: %w", err))
			}
			return formatResponse(true, "Directive set successfully", nil)
		},
		Help:      "Sets the system directive",
		MinAccess: AccessBeta,
	}
}

func addHistoryCommands(cmd *Command) {
	cmd.SubCmds["show"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			var hist []string
			for _, m := range c.GetHistory() {
				hist = append(hist, fmt.Sprintf("%s: %s", m.Role, m.Content))
			}
			return formatResponse(true, strings.Join(hist, "\n"), nil)
		},
		Help:      "Shows the conversation history",
		MinAccess: AccessBeta,
	}
	cmd.SubCmds["clear"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			c.ClearHistory()
			return formatResponse(true, "History cleared", nil)
		},
		Help:      "Clears the conversation history",
		MinAccess: AccessBeta,
	}
}

func addModelCommands(cmd *Command) {
	cmd.SubCmds["list"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			models := c.ListModels()
			return formatResponse(true, strings.Join(models, "\n"), nil)
		},
		Help:      "Lists available models",
		MinAccess: AccessBeta,
	}
	cmd.SubCmds["set"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			if len(args) < 1 {
				return formatResponse(false, "", fmt.Errorf("usage: /model set <name>"))
			}
			c.SetModel(args[0])
			return formatResponse(true, fmt.Sprintf("Model set to %s", args[0]), nil)
		},
		Help:      "Sets the current model",
		MinAccess: AccessBeta,
	}
}

func addConversationCommands(cmd *Command) {
	cmd.SubCmds["list"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			convos := c.ListConversations()
			return formatResponse(true, strings.Join(convos, "\n"), nil)
		},
		Help:      "Lists saved conversations",
		MinAccess: AccessBeta,
	}
	cmd.SubCmds["save"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			if len(args) < 1 {
				return formatResponse(false, "", fmt.Errorf("usage: /conversation save <name>"))
			}
			if err := c.SaveConversation(args[0]); err != nil {
				return formatResponse(false, "", fmt.Errorf("failed to save conversation: %w", err))
			}
			return formatResponse(true, "Conversation saved successfully", nil)
		},
		Help:      "Saves the current conversation",
		MinAccess: AccessBeta,
	}
	cmd.SubCmds["load"] = &Command{
		Execute: func(c ChatClient, args []string) string {
			if len(args) < 1 {
				return formatResponse(false, "", fmt.Errorf("usage: /conversation load <name>"))
			}
			if err := c.LoadConversation(args[0]); err != nil {
				return formatResponse(false, "", fmt.Errorf("failed to load conversation: %w", err))
			}
			return formatResponse(true, "Conversation loaded successfully", nil)
		},
		Help:      "Loads a saved conversation",
		MinAccess: AccessBeta,
	}
}
