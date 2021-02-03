package runner

import (
	"testing"
)

func TestIsWatchedFile(t *testing.T) {
	tests := []struct {
		file     string
		expected bool
	}{
		{"test.go", true},
		{"test.tpl", true},
		{"test.tmpl", true},
		{"test.html", true},
		{"test.css", false},
		{"test-executable", false},
		{"./tmp/test.go", false},
	}

	for _, test := range tests {
		actual := isWatchedFile(test.file)

		if actual != test.expected {
			t.Errorf("Expected %v for %v, got %v", test.expected, test.file, actual)
		}
	}
}

func TestShouldRebuild(t *testing.T) {
	tests := []struct {
		eventName string
		expected  bool
	}{
		{`"test.go": MODIFIED`, true},
		{`"test.tpl": MODIFIED`, false},
		{`"test.tmpl": DELETED`, false},
		{`"unknown.extension": DELETED`, true},
		{`"no_extension": ADDED`, true},
		{`"./a/path/test.go": MODIFIED`, true},
	}

	for _, test := range tests {
		actual := shouldRebuild(test.eventName)

		if actual != test.expected {
			t.Errorf("Expected %v, got %v (event was '%s')", test.expected, actual, test.eventName)
		}
	}
}

func TestIsIgnoredFolder(t *testing.T) {
	settings["ignored"] = "tmp, assets, app/controllers"
	tests := []struct {
		dir      string
		expected bool
	}{
		{"assets", true},
		{"assets/node_modules", true},
		{"tmp", true},
		{"tmp/pid", true},
		{"app", false},
		{"app/controllers", true},
		{"app/controllers/user", true},
		{"app/views", false},
	}

	for _, test := range tests {
		actual := isIgnoredFolder(test.dir)
		if actual != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}
