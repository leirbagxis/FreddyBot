package parser

import (
	"os"
	"testing"
)

func TestLoadMessagesResilience(t *testing.T) {
	// 1. Test with existing file in root
	// Assuming we are running from the root of the project
	// pkg/parser is 2 levels deep from root
	// So it should find ../../config/messages.yml
	
	// We need to make sure we don't trigger the real loadOnce
	// But since it's a test, we can just call loadMessages directly if it wasn't private
	// Oh, loadMessages is private.
	
	// Let's create a dummy messages.yml in the current test dir
	os.MkdirAll("config", 0755)
	dummyContent := `- name: test
  text: "Hello World"
`
	os.WriteFile("config/messages.yml", []byte(dummyContent), 0644)
	defer os.RemoveAll("config")

	// Now we can call loadMessages via GetMessage
	// But loadOnce might have already been triggered if other tests ran
	// In this session, it's the first time.
	
	text, _ := GetMessage("test", nil)
	if text != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", text)
	}
}
