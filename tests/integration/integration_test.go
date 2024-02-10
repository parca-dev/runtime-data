
package integration

import (
	"flag"
	"os"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

// sanitizeIdentifier sanitizes the identifier to be used as a filename.
func sanitizeIdentifier(identifier string) string {
	return strings.ReplaceAll(identifier, ".", "_")
}
