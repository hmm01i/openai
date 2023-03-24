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
	switch cmd {
	case "/listper":
		p := c.listPersonas()
		fmt.Println(strings.Join(p, "\n"))
	case "/saveper":
		c.savePersona(args[0], c.systemDirective)
	case "/showper":
		c.showPersona()
	case "/loadper":
		if len(args) < 1 {
			println(`ERR: No persona provided.
		Usage: /loadper <persona>`)
			break
		}
		c.loadPersona(args[0])
	case "/setdir":
		c.setDirective(strings.Trim(strings.Join(args, " "), `"`))
	case "/history", "/hist":
		for _, m := range c.history {
			fmt.Printf("%s: %s\n", m.Role, m.Content)
		}
	case "/clearhist":
		c.clearHistory()
	case "/listmod":
		c.listModels()
	case "/q", "/quit":
		os.Exit(0)
	}
	return true
}
