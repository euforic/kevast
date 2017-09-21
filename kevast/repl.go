package kevast

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	errMissingKey   = "Expected: %s <key>"
	errMissingVal   = "Expected: %s <key> <val>"
	errNotFound     = "No key found for: %s"
	errReadingStdin = "Error reading from Stdin: %s"
	errInvalidCmd   = "Unknown command \"%s\""
)

// Session is the base struct for the current
// REPL session
type Session struct {
	db *Kevast
}

// NewSession initilizes a new Session with default values
func NewSession() Session {
	return Session{
		db: &Kevast{
			idx: 0,
			stores: []store{
				store{},
			},
		},
	}
}

// Run starts the REPL session and listens on Stdin for new
// commands
func (s Session) Run() error {
	s.printCursor()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		in := scanner.Text()

		if in == "" {
			s.printCursor()
			continue
		}

		p := make([]string, 3)
		copy(p, strings.Fields(in))

		err := s.eval(p[0], p[1], p[2])
		if err != nil {
			printErr(err)
		}

		s.printCursor()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf(errReadingStdin, err)
	}

	return nil
}

// eval parses the command and determines which function
// to call
func (s Session) eval(cmd string, key string, val string) error {
	if cmd == "write" && val == "" {
		return fmt.Errorf(errMissingVal, cmd)
	}

	if (cmd == "read" || cmd == "delete") && key == "" {
		return fmt.Errorf(errMissingVal, cmd)
	}

	switch cmd {
	case "write":
		return s.db.Write(key, val)
	case "read":
		v, err := s.db.Read(key)
		if err != nil {
			return err
		}
		fmt.Print(v)
	case "delete":
		return s.db.Del(key)
	case "start":
		return s.db.Start()
	case "abort":
		return s.db.Abort()
	case "commit":
		return s.db.Commit()
	case "quit":
		return s.quit()
	case "help":
		return s.help()
	case "dump":
		s.dump()
		return nil
	default:
		return fmt.Errorf(errInvalidCmd, cmd)
	}

	return nil
}

// dump is a debug tool that dumps current
// state of the kevast instance
func (s Session) dump() {
	fmt.Print("idx:", s.db.idx)
	pad := "\n"
	for _, stores := range s.db.stores {
		pad += "-"
		fmt.Print(pad)
		for key, val := range stores {
			fmt.Print("|", key, ":", val, "|")
		}
	}
}

// printCursor is a simple helper function to print out REPL the cursor
// with multiple `>` to show the current context (aka transaction depth)
func (s Session) printCursor() {
	fmt.Print("\n" + strings.Repeat(">", len(s.db.stores)) + " ")
}

// quit will exit the current REPL session
func (Session) quit() error {
	os.Exit(0)
	return nil
}

// help prints out the available REPL commands usages
func (Session) help() error {
	fmt.Print(`Available commands are:
		
	READ   <key>            Gets the value for the given key
	WRITE  <key> <val>      Sets the value for the given key
	DELETE <key>            Deletes the value for a given key
	START                   Starts a transaction
	ABORT                   Rolls back and exits the current transaction
	COMMIT                  Commits the transactions chages
	QUIT                    Exits the REPL
	HELP                    Returns available commands and usage
	`)
	return nil
}

// printErr lazy helper function for printing to Stderr
func printErr(a ...interface{}) {
	fmt.Fprint(os.Stderr, a...)
}
