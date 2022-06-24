// Package replacer implements tag's search and replace functionality.
package replacer

import (
	"bytes"
	"fmt"
	"os"
)

// Replace finds all occurrences of 'search' in the file 'path' and replaces
// them with 'replace'. If the search term is not in the file, it
// will return an error rather than do nothing.
func Replace(path, search, replace string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if !bytes.Contains(contents, []byte(search)) {
		return fmt.Errorf("Could not find %s in %s", search, path)
	}

	newContent := bytes.ReplaceAll(contents, []byte(search), []byte(replace))

	if err := os.WriteFile(path, newContent, 0755); err != nil {
		return err
	}
	return nil
}
