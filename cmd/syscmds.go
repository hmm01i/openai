package main

import (
	"fmt"
	"os"
	"strings"
)

func (c *chatClient) command(input string) (string, bool) {
	if !strings.HasPrefix(input, "/") {
		return "", false
	}
	inputSlice := strings.Split(input, " ")
	var cmd string
	var args []string
	cmd = inputSlice[0]
	if len(inputSlice) > 1 {
		args = inputSlice[1:]
	}

	if _, ok := cmds[cmd]; ok {
		return cmds[cmd](c, args), true
	}
	switch cmd {

	case "/listcmd":
		commands := []string{}
		return func(c *chatClient, args []string) string {
			for o, _ := range cmds {
				commands = append(commands, o)
			}
			return strings.Join(commands, "\n")
		}(c, args), true
	}
	return "Unrecognized command (try /listcmd)", true
}

var cmds = map[string]func(c *chatClient, args []string) string{
	//Personas
	"/listper": func(c *chatClient, args []string) string {
		personas := c.listPersonas()
		for i, p := range personas {
			if p == c.persona {
				personas[i] = p + "*"
			}
		}
		return strings.Join(personas, "\n")
	},
	"/saveper": func(c *chatClient, args []string) string { c.savePersona(args[0], c.systemDirective); return "ok" },
	"/showper": func(c *chatClient, args []string) string { persona := c.showPersona(); return persona },
	"/loadper": func(c *chatClient, args []string) string {
		if len(args) < 1 {
			return `ERR: No persona provided. Usage: /loadper <persona>`
		}
		if err := c.loadPersona(args[0]); err != nil {
			return "failed to load persona"
		}
		return "ok"
	},
	// system directive
	"/setdir": func(c *chatClient, args []string) string {
		if err := c.setDirective(strings.Trim(strings.Join(args, " "), `"`)); err != nil {
			return "failed to set directive"
		}
		return "ok"
	},
	// conversation history
	"/hist": func(c *chatClient, args []string) string {
		var hist []string
		for _, m := range c.history {
			hist = append(hist, fmt.Sprintf("%s: %s\n", m.Role, m.Content))
		}
		return strings.Join(hist, "\n")
	},
	"/clearhist": func(c *chatClient, args []string) string { c.clearHistory(); return "ok" },
	// models
	"/listmod": func(c *chatClient, args []string) string { mod := c.listModels(); return strings.Join(mod, "\n") },
	"/setmod":  func(c *chatClient, args []string) string { c.model = args[0]; return "model set" },
	// conversations
	"/saveconv": func(c *chatClient, args []string) string { c.saveConversation(args[0]); return "ok" },
	"/listconv": func(c *chatClient, args []string) string {
		convos := c.listConversations()
		return strings.Join(convos, "\n")
	},
	"/loadconv": func(c *chatClient, args []string) string { c.loadConversation(args[0]); return "ok" },

	// quit
	"/q": func(c *chatClient, args []string) string { os.Exit(0); return "ok" },
}
