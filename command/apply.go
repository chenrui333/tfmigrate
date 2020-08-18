package command

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/minamijoyo/tfmigrate/config"
)

// ApplyCommand is a command which computes a new state and pushes it to remote state.
type ApplyCommand struct {
	Meta
	path string
}

// Run runs the procedure of this command.
func (c *ApplyCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("apply", flag.ContinueOnError)
	if err := cmdFlags.Parse(args); err != nil {
		c.UI.Error(fmt.Sprintf("failed to parse arguments: %s", err))
		return 1
	}

	if len(cmdFlags.Args()) != 1 {
		c.UI.Error(fmt.Sprintf("The command expects 1 argument, but got %d", len(cmdFlags.Args())))
		c.UI.Error(c.Help())
		return 1
	}

	c.Option = newOption()
	// The option may contains sensitive values such as environment variables.
	// So logging the option set log level to DEBUG instead of INFO.
	log.Printf("[DEBUG] [command] option: %#v\n", c.Option)

	c.path = cmdFlags.Arg(0)
	log.Printf("[INFO] [command] read migration file: %s\n", c.path)
	source, err := ioutil.ReadFile(c.path)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	log.Printf("[DEBUG] [command] parse migration file: %#v\n", string(source))
	config, err := config.ParseMigrationFile(c.path, source)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	log.Printf("[INFO] [command] new migrator: %#v\n", config)
	m, err := config.NewMigrator(c.Option)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	err = m.Apply(context.Background())
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

// Help returns long-form help text.
func (c *ApplyCommand) Help() string {
	helpText := `
Usage: tfmigrate apply <PATH>

Apply computes a new state and pushes it to remote state.
It will fail if terraform plan detects any diffs with the new state.

Arguments
  PATH               A path of migration file
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns one-line help text.
func (c *ApplyCommand) Synopsis() string {
	return "Computes a new state and pushes it to remote state"
}