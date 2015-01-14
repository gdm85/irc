package irc

import "strings"

// An identity represents the prefix of a message, generally the user who sent it
type Identity struct {
	// This is what the Identity was parsed from
	Raw string

	// The nick will either contain the nick of who sent the message or a blank string
	Nick string

	// The nick will either contain the user who sent the message or a blank string
	User string

	// The nick will either contain the host of who sent the message or a blank string
	Host string
}

type Event struct {
	// This is where the Command was parsed from.
	Raw string

	// The Identity is also the prefix of the message.
	Identity *Identity

	// The prefix is essentially a copy of the Raw identity.
	Prefix string

	// Command is which command is being called.
	Command string

	// Arguments are all the arguments for the command.
	Args []string
}

// https://github.com/kylelemons/blightbot/blob/master/bot/parser.go#L34
func ParseIdentity(line string) *Identity {
	// Start by creating an Identity with nothing but the host
	id := &Identity{
		Raw:  line,
		Host: line,
	}

	uh := strings.SplitN(id.Host, "@", 2)
	if len(uh) != 2 {
		return id
	}
	id.User, id.Host = uh[0], uh[1]

	nu := strings.SplitN(id.User, "!", 2)
	if len(nu) != 2 {
		return id
	}
	id.Nick, id.User = nu[0], nu[1]

	return id
}

func (i *Identity) Copy() *Identity {
	newIdent := &Identity{}

	*newIdent = *i

	return newIdent
}

// https://github.com/kylelemons/blightbot/blob/master/bot/parser.go#L55
func ParseEvent(line string) *Event {
	// Trim the line and make sure we have data
	line = strings.TrimSpace(line)
	if len(line) <= 0 {
		return nil
	}

	c := &Event{
		Raw: line,
	}

	if line[0] == ':' {
		split := strings.SplitN(line, " ", 2)
		if len(split) <= 1 {
			return nil
		}
		c.Prefix = string(split[0][1:])
		line = split[1]
	}

	// While the current arg doesn't start with ":"
	for len(line) > 0 && line[0] != ':' {
		split := strings.SplitN(line, " ", 2)
		if len(split) < 2 {
			c.Args = append(c.Args, split...)
			break
		}

		c.Args = append(c.Args, split[0])

		line = split[1]
	}

	if len(line) > 0 && line[0] == ':' {
		c.Args = append(c.Args, line[1:])
	}

	// The first arg will be the command
	if len(c.Args) < 1 {
		return nil
	}

	c.Command = c.Args[0]
	c.Args = c.Args[1:]

	c.Identity = ParseIdentity(c.Prefix)

	return c
}

// This returns the last argument in the Event or an empty string
// if there are no args
func (e *Event) Trailing() string {
	if len(e.Args) < 1 {
		return ""
	}

	return e.Args[len(e.Args)-1]
}

// This is mostly for PRIVMSG events (and similar derived events)
// It will check if the event came from a channel or a person.
func (e *Event) FromChannel() bool {
	if len(e.Args) < 1 || len(e.Args[0]) < 1 {
		return false
	}

	switch e.Args[0][0] {
	case '#', '&':
		return true
	default:
		return false
	}
}

func (e *Event) Copy() *Event {
	// Create a new event
	newEvent := &Event{}

	// Copy stuff from the old event
	*newEvent = *e

	// Copy the Identity
	newEvent.Identity = e.Identity.Copy()

	// Copy the Args slice
	newEvent.Args = append(make([]string, 0, len(e.Args)), e.Args...)

	return newEvent
}

func (e *Event) Message() string {
	if len(e.Args) == 0 {
		return ""
	}
	return e.Args[len(e.Args)-1]
}
