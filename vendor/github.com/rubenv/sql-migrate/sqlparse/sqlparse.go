package sqlparse

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"strings"
)

const (
	sqlCmdPrefix        = "-- +migrate "
	optionNoTransaction = "notransaction"
)

type ParsedMigration struct {
	UpStatements   []string
	DownStatements []string

	DisableTransactionUp   bool
	DisableTransactionDown bool
}

// Checks the line to see if the line has a statement-ending semicolon
// or if the line contains a double-dash comment.
func endsWithSemicolon(line string) bool {

	prev := ""
	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		if strings.HasPrefix(word, "--") {
			break
		}
		prev = word
	}

	return strings.HasSuffix(prev, ";")
}

type migrationDirection int

const (
	directionNone migrationDirection = iota
	directionUp
	directionDown
)

type migrateCommand struct {
	Command string
	Options []string
}

func (c *migrateCommand) HasOption(opt string) bool {
	for _, specifiedOption := range c.Options {
		if specifiedOption == opt {
			return true
		}
	}

	return false
}

func parseCommand(line string) (*migrateCommand, error) {
	cmd := &migrateCommand{}

	if !strings.HasPrefix(line, sqlCmdPrefix) {
		return nil, errors.New("ERROR: not a sql-migrate command")
	}

	fields := strings.Fields(line[len(sqlCmdPrefix):])
	if len(fields) == 0 {
		return nil, errors.New(`ERROR: incomplete migration command`)
	}

	cmd.Command = fields[0]

	cmd.Options = fields[1:]

	return cmd, nil
}

// Split the given sql script into individual statements.
//
// The base case is to simply split on semicolons, as these
// naturally terminate a statement.
//
// However, more complex cases like pl/pgsql can have semicolons
// within a statement. For these cases, we provide the explicit annotations
// 'StatementBegin' and 'StatementEnd' to allow the script to
// tell us to ignore semicolons.
func ParseMigration(r io.ReadSeeker) (*ParsedMigration, error) {
	p := &ParsedMigration{}

	_, err := r.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	scanner := bufio.NewScanner(r)

	statementEnded := false
	ignoreSemicolons := false
	currentDirection := directionNone

	for scanner.Scan() {
		line := scanner.Text()

		// handle any migrate-specific commands
		if strings.HasPrefix(line, sqlCmdPrefix) {
			cmd, err := parseCommand(line)
			if err != nil {
				return nil, err
			}

			switch cmd.Command {
			case "Up":
				currentDirection = directionUp
				if cmd.HasOption(optionNoTransaction) {
					p.DisableTransactionUp = true
				}
				break

			case "Down":
				currentDirection = directionDown
				if cmd.HasOption(optionNoTransaction) {
					p.DisableTransactionDown = true
				}
				break

			case "StatementBegin":
				if currentDirection != directionNone {
					ignoreSemicolons = true
				}
				break

			case "StatementEnd":
				if currentDirection != directionNone {
					statementEnded = (ignoreSemicolons == true)
					ignoreSemicolons = false
				}
				break
			}
		}

		if currentDirection == directionNone {
			continue
		}

		if _, err := buf.WriteString(line + "\n"); err != nil {
			return nil, err
		}

		// Wrap up the two supported cases: 1) basic with semicolon; 2) psql statement
		// Lines that end with semicolon that are in a statement block
		// do not conclude statement.
		if (!ignoreSemicolons && endsWithSemicolon(line)) || statementEnded {
			statementEnded = false
			switch currentDirection {
			case directionUp:
				p.UpStatements = append(p.UpStatements, buf.String())

			case directionDown:
				p.DownStatements = append(p.DownStatements, buf.String())

			default:
				panic("impossible state")
			}

			buf.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// diagnose likely migration script errors
	if ignoreSemicolons {
		return nil, errors.New("ERROR: saw '-- +migrate StatementBegin' with no matching '-- +migrate StatementEnd'")
	}

	if currentDirection == directionNone {
		return nil, errors.New(`ERROR: no Up/Down annotations found, so no statements were executed.
			See https://github.com/rubenv/sql-migrate for details.`)
	}

	return p, nil
}
