package main

import (
	"fmt"
	"os"
	"strings"
)

func (c *client) command(input string) bool {
	if !strings.HasPrefix(input, "/") {
		return false
	}
	inputSlice := strings.Split(input, " ")
	var cmd string
	var args []string
	cmd = inputSlice[0]
	if len(inputSlice) > 1 {
		args = inputSlice[1:]
	}

	if _, ok := cmds[cmd]; ok {
		cmds[cmd](c, args)
		return true
	}
	switch cmd {

	case "/listcmd":
		func(c *client, args []string) {
			for o, _ := range cmds {
				fmt.Println(o)
			}
		}(c, args)
	default:

		fmt.Println("Unrecognized command")
	}
	return true
}

var cmds = map[string]func(c *client, args []string){
	"/listper": func(c *client, args []string) {
		personas := c.listPersonas()
		for _, p := range personas {
			if p == c.persona {
				fmt.Println(p + "*")
			} else {
				fmt.Println(p)
			}
		}
	},
	"/saveper": func(c *client, args []string) { c.savePersona(args[0], c.systemDirective) },
	"/showper": func(c *client, args []string) { c.showPersona() },
	"/loadper": func(c *client, args []string) {
		if len(args) < 1 {
			fmt.Println(`ERR: No persona provided. Usage: /loadper <persona>`)
			return
		}
		c.loadPersona(args[0])
	},
	"/setdir": func(c *client, args []string) {
		c.setDirective(strings.Trim(strings.Join(args, " "), `"`))
	},
	"/hist": func(c *client, args []string) {
		for _, m := range c.history {
			fmt.Printf("%s: %s\n", m.Role, m.Content)
		}
	},
	"/clearhist": func(c *client, args []string) { c.clearHistory() },
	"/listmod":   func(c *client, args []string) { c.listModels() },
	"/q":         func(c *client, args []string) { os.Exit(0) },
	"/saveconv":  func(c *client, args []string) { c.saveConversation(args[0]) },
	"/listconv":  func(c *client, args []string) { c.listConversations() },
	"/loadconv":  func(c *client, args []string) { c.loadConversation(args[0]) },
}
