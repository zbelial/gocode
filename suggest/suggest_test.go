package suggest_test

import (
	"bytes"
	"go/importer"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mdempsky/gocode/suggest"
)

func TestRegress(t *testing.T) {
	testDirs, err := filepath.Glob("testdata/test.*")
	if err != nil {
		t.Fatal(err)
	}

	for _, testDir := range testDirs {
		testDir := testDir // capture
		name := strings.TrimPrefix(testDir, "testdata/")
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testRegress(t, testDir)
		})
	}
}

func testRegress(t *testing.T, testDir string) {
	testDir, err := filepath.Abs(testDir)
	if err != nil {
		t.Errorf("Abs failed: %v", err)
		return
	}

	filename := filepath.Join(testDir, "test.go.in")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("ReadFile failed: %v", err)
		return
	}

	cursor := bytes.IndexByte(data, '@')
	if cursor < 0 {
		t.Errorf("Missing @")
		return
	}
	data = append(data[:cursor], data[cursor+1:]...)

	cfg := suggest.Config{
		Importer: importer.Default(),
	}
	if testing.Verbose() {
		cfg.Logf = t.Logf
	}
	candidates, prefixLen := cfg.Suggest(filename, data, cursor)

	var out bytes.Buffer
	suggest.NiceFormat(&out, candidates, prefixLen)

	want, _ := ioutil.ReadFile(filepath.Join(testDir, "out.expected"))
	if got := out.Bytes(); !bytes.Equal(got, want) {
		t.Errorf("%s:\nGot:\n%s\nWant:\n%s\n", testDir, got, want)
		return
	}
}
