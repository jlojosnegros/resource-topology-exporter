package notification

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fsnotify/fsnotify"
)

func TestMakeFilter(t *testing.T) {
	type testCase struct {
		name     string
		filters  []FilterEvent
		event    fsnotify.Event
		expected bool
	}

	testCases := []testCase{
		{
			name:     "nothing, nothing",
			filters:  []FilterEvent{FilterNothing, FilterNothing},
			expected: false,
		},
		{
			name:     "everything, nothing",
			filters:  []FilterEvent{FilterEverything, FilterNothing},
			expected: true,
		},
		{
			name:     "nothing, everything",
			filters:  []FilterEvent{FilterNothing, FilterEverything},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := MakeFilter(tc.filters...)
			got := f(tc.event)
			if tc.expected != got {
				t.Fatalf("%s event=%v expected=%t got=%t", tc.name, tc.event, tc.expected, got)
			}
		})
	}
}

func TestEnsureNotifyFilePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testnotif")
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "nested", "notify")

	err = EnsureNotifyFilePath(testPath, false)
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	// try again on the same path: should fail now.
	err = EnsureNotifyFilePath(testPath, false)
	if err == nil {
		t.Fatalf("unexpected SUCCESS: %v", err)
	}
}

func TestEnsureNotifyFilePathAllowExistingWithSafeContent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testnotif")
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "nested", "notify")

	err = EnsureNotifyFilePath(testPath, false)
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	err = EnsureNotifyFilePath(testPath, true)
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}
}

// TODO: test symlink
func TestEnsureNotifyFilePathAllowExistingWithUNSAFEContent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testnotif")
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "nested", "notify")

	err = EnsureNotifyFilePath(testPath, false)
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	err = os.WriteFile(testPath, []byte("the file is no longer empty"), 0644)
	if err != nil {
		t.Fatalf("error setting content: %v", err)
	}

	err = EnsureNotifyFilePath(testPath, true)
	if err == nil {
		t.Fatalf("unexpected SUCCESS: %v", err)
	}
}
