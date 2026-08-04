package load

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// queryRefPattern matches the config.App().QUERY["NAME"] lookups used across the code base.
var queryRefPattern = regexp.MustCompile(`QUERY\["([A-Z0-9_]+)"\]`)

// TestEveryReferencedQueryExists guards against the failure mode that left
// SCHEDULE_LOG_INSERT undefined: a missing map key yields "" and db.Prepare("")
// then fails at runtime, silently, on every execution.
func TestEveryReferencedQueryExists(t *testing.T) {
	root := repoRoot(t)

	file, err := os.Open(filepath.Join(root, "query.sql"))
	if err != nil {
		t.Fatalf("open query.sql: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	queries, err := parseSQLQueries(file, make(map[string]string))
	if err != nil {
		t.Fatalf("parse query.sql: %v", err)
	}

	var missing []string
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		// test files are skipped: this one carries the pattern itself
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, match := range queryRefPattern.FindAllStringSubmatch(string(content), -1) {
			name := match[1]
			if _, ok := queries[name]; !ok {
				rel, _ := filepath.Rel(root, path)
				missing = append(missing, name+" ("+rel+")")
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk repository: %v", err)
	}

	if len(missing) > 0 {
		sort.Strings(missing)
		t.Fatalf("queries referenced in code but not defined in query.sql:\n%s", strings.Join(missing, "\n"))
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found, cannot locate repository root")
		}
		dir = parent
	}
}
