package html

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testEmailData struct {
	Subject     string
	MainMessage string
	CtaLink     string
	CtaText     string
}

func TestTemplateHtml(t *testing.T) {
	file, err := os.ReadFile("./test-files/test-template.html")
	if err != nil {
		t.Fatal(err)
	}
	reader := strings.NewReader(string(file))

	data := testEmailData{
		Subject:     "TEST",
		MainMessage: "This is a longer test string",
		CtaLink:     "https://metasoundtools.com",
		CtaText:     "Please click me",
	}

	buf, err := TemplateEmail(reader, data)
	if err != nil {
		t.Fatal(err)
	}
	goldenFile, err := os.ReadFile("./test-files/example-template.html")
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Printf("Formatted Template: %v", buf.String())

	if buf.String() != string(goldenFile) {
		diff := cmp.Diff(string(goldenFile), buf.String())
		t.Fatalf("Files do not match (-want, +got) \n%s", diff)
	}
}
