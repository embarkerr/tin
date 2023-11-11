package b64

import (
	"testing"
)

func TestBase64Encode(t *testing.T) {
	input := "Hello, world!"
	output := Execute([]string{input})
	if output != "SGVsbG8sIHdvcmxkIQ==" {
		t.Fatalf("b64 encode output incorrect\nexpect: SGVsbG8sIHdvcmxkIQ==\nactual: %s\n", output)
	}
}

func TestBase64Decode(t *testing.T) {
	input := "SGVsbG8sIHdvcmxkIQ=="
	output := Execute([]string{"-d", input})
	if output != "Hello, world!" {
		t.Fatalf("b64 decode output incorrect\nexpect: Hello, world!\nactual: %s\n", output)
	}
}

func TestBase64URLEncode(t *testing.T) {
	input := "Hello, world!"
	output := Execute([]string{"-u", input})
	if output != "SGVsbG8sIHdvcmxkIQ" {
		t.Fatalf("b64 URL encode output incorrect\nexpect: SGVsbG8sIHdvcmxkIQ\nactual: %s\n", output)
	}
}

func TestBase64URLDecode(t *testing.T) {
	input := "SGVsbG8sIHdvcmxkIQ"
	output := Execute([]string{"-u", "-d", input})
	if output != "Hello, world!" {
		t.Fatalf("b64 URL decode output incorrect\nexpect: Hello, world!\nactual: %s\n", output)
	}
}
