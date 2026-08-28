package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"go-boilerplate/internal/config"
)

const migrationsDir = "migrations"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	var databaseURL string
	if args[0] != "create" {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		databaseURL = cfg.DatabaseURL
	}

	var err error
	switch args[0] {
	case "create":
		err = create(args[1:])
	case "up":
		err = run(databaseURL, args[1:], 1)
	case "down":
		err = run(databaseURL, args[1:], -1)
	case "version":
		err = version(databaseURL)
	case "force":
		err = force(databaseURL, args[1:])
	case "fresh":
		err = fresh(databaseURL)
	default:
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func create(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: go run ./cmd/migrate create <name>")
	}

	name := sanitize(args[0])
	if name == "" {
		return errors.New("migration name is empty")
	}

	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return err
	}

	version, err := nextVersion()
	if err != nil {
		return err
	}

	prefix := fmt.Sprintf("%06d_%s", version, name)
	for _, suffix := range []string{"up", "down"} {
		path := filepath.Join(migrationsDir, prefix+"."+suffix+".sql")
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			return err
		}

		if _, err := file.WriteString("-- Write your migration SQL here.\n"); err != nil {
			_ = file.Close()
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}

		fmt.Println("created", path)
	}

	return nil
}

func run(databaseURL string, args []string, direction int) error {
	m, err := newMigrator(databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if direction == 1 && len(args) == 0 {
		err = m.Up()
	} else {
		steps := 1
		if len(args) > 0 {
			steps, err = strconv.Atoi(args[0])
			if err != nil || steps < 1 {
				return errors.New("steps must be a positive number")
			}
		}
		err = m.Steps(direction * steps)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("no migration changes")
		return nil
	}

	return err
}

func version(databaseURL string) error {
	m, err := newMigrator(databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		fmt.Println("version: none")
		return nil
	}
	if err != nil {
		return err
	}

	fmt.Printf("version: %d dirty: %t\n", version, dirty)
	return nil
}

func force(databaseURL string, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: go run ./cmd/migrate force <version>")
	}

	version, err := strconv.Atoi(args[0])
	if err != nil || version < 0 {
		return errors.New("version must be a non-negative number")
	}

	m, err := newMigrator(databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	return m.Force(version)
}

func newMigrator(databaseURL string) (*migrate.Migrate, error) {
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	return migrate.New("file://"+migrationsDir, databaseURL)
}

func nextVersion() (int, error) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return 0, err
	}

	re := regexp.MustCompile(`^(\d+)_.*\.(up|down)\.sql$`)
	versions := make([]int, 0)
	for _, entry := range entries {
		matches := re.FindStringSubmatch(entry.Name())
		if len(matches) == 0 {
			continue
		}

		version, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, err
		}
		versions = append(versions, version)
	}

	if len(versions) == 0 {
		return 1, nil
	}

	sort.Ints(versions)
	return versions[len(versions)-1] + 1, nil
}

func sanitize(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	var builder strings.Builder
	lastUnderscore := false

	for _, r := range name {
		valid := unicode.IsLetter(r) || unicode.IsDigit(r)
		if valid {
			builder.WriteRune(r)
			lastUnderscore = false
			continue
		}

		if !lastUnderscore {
			builder.WriteByte('_')
			lastUnderscore = true
		}
	}

	return strings.Trim(builder.String(), "_")
}
func fresh(databaseURL string) error {
	m, err := newMigrator(databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Drop(); err != nil {
		return err
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func usage() {
	fmt.Println(`Usage:
  go run ./cmd/migrate create <name>
  go run ./cmd/migrate up [steps]
  go run ./cmd/migrate down [steps]
  go run ./cmd/migrate version
  go run ./cmd/migrate force <version>
  go run ./cmd/migrate fresh`)
}
