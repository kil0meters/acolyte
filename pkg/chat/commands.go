package chat

import (
	"fmt"
	"strings"
	"time"

	"github.com/kil0meters/acolyte/pkg/authorization"
)

// CommandFunction a datatype for commands
type CommandFunction func(*Client, []string) string

// Command a struct containing data for a command
type Command struct {
	Name               string                        `json:"name"`
	Description        string                        `json:"description"`
	RequiredPermission authorization.PermissionLevel `json:"-"`
	Function           CommandFunction               `json:"-"`
}

// Commands array storing commands
var Commands map[*Command]int = make(map[*Command]int)

// InitializeCommands initializes commands
func InitializeCommands() {
	Commands = map[*Command]int{
		&Command{
			Name:               "/ping",
			Description:        "tests ping",
			RequiredPermission: authorization.Standard,
			Function: func(*Client, []string) string {
				return "pong"
			},
		}: 0,
		&Command{
			Name:               "/addcommand",
			Description:        "[name] [output...]",
			RequiredPermission: authorization.Moderator,
			Function:           AddCommand,
		}: 0,
		&Command{
			Name:               "/removecommand",
			Description:        "[name]",
			RequiredPermission: authorization.Moderator,
			Function:           RemoveCommand,
		}: 0,
		&Command{
			Name:               "/ban",
			Description:        "[user] [duration] [reason...]",
			RequiredPermission: authorization.Moderator,
			Function:           BanAccount,
		}: 0,
		&Command{
			Name:               "/mod",
			Description:        "[user]",
			RequiredPermission: authorization.Admin,
			Function:           BanAccount,
		}: 0,
		&Command{
			Name:               "/unmod",
			Description:        "[user]",
			RequiredPermission: authorization.Admin,
			Function:           BanAccount,
		}: 0,
	}
}

// ParseCommand parses a command
func ParseCommand(client *Client, command string) string {
	tokens := strings.Split(command, " ")

	for command := range Commands {
		if tokens[0] == command.Name && client.Account.Permissions.AtLeast(command.RequiredPermission) {
			return command.Function(client, tokens)
		}
	}

	return "Unknown command"
}

// AddCommand adds a command
func AddCommand(client *Client, tokens []string) string {
	if len(tokens) < 2 {
		return fmt.Sprintf("error: Expected at least 2 arguments, got %d instead", len(tokens))
	}

	name := "/" + tokens[1]
	output := strings.Join(tokens[2:], " ")

	Commands[&Command{
		Name:               name,
		Description:        "",
		RequiredPermission: authorization.Standard,
		Function: func(*Client, []string) string {
			return output
		},
	}] = 0

	for client := range client.Pool.Clients {
		client.UpdateCommands()
	}

	return fmt.Sprintf("Created command \"%s\"", name)
}

// RemoveCommand removes a command
func RemoveCommand(client *Client, tokens []string) string {
	if len(tokens) != 2 {
		return fmt.Sprintf("error: Expected 2 arguments, got %d instead", len(tokens))
	}

	for command := range Commands {
		if command.Name == "/"+tokens[1] {
			delete(Commands, command)

			for client := range client.Pool.Clients {
				client.UpdateCommands()
			}

			return fmt.Sprintf("Removed command \"%s\"", tokens[1])
		}
	}

	return fmt.Sprintf("error: Command \"%s\" not found", tokens[1])
}

// BanAccount bans an account
func BanAccount(client *Client, tokens []string) string {
	if len(tokens) < 4 {
		return fmt.Sprintf("error: Expected at least 4 arguments, got %d instead", len(tokens))
	}

	account := authorization.AccountFromUsername(tokens[1])
	if account == nil {
		return "error: A user with that username does not exist"
	}
	duration, err := time.ParseDuration(tokens[2])
	if err != nil {
		return "error: Invalid time"
	}
	reason := strings.Join(tokens[3:], " ")

	client.Pool.Broadcast <- Message{
		Type: 1,
		Data: MessageData{
			Username: "Ban",
			Text:     fmt.Sprintf("User \"%s\" has been banned for %s. \"%s\"", account.Username, duration, reason),
		},
	}

	return "User successfully banned"
}
