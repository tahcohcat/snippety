package git

import (
	"testing"
)

func TestExtractTicketPrefix(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		expected   string
	}{
		{
			name:       "Standard JIRA-like ticket at start",
			branchName: "FEAT-1234-add-new-feature",
			expected:   "FEAT-1234: ",
		},
		{
			name:       "Bug fix ticket",
			branchName: "BUG-5678-fix-login-issue",
			expected:   "BUG-5678: ",
		},
		{
			name:       "Backlog item",
			branchName: "BP-9999-user-story",
			expected:   "BP-9999: ",
		},
		{
			name:       "DevOps ticket with slash prefix",
			branchName: "feature/DEVOPS-123",
			expected:   "DEVOPS-123: ",
		},
		{
			name:       "Chore with ticket",
			branchName: "chore/PROJ-456-cleanup",
			expected:   "PROJ-456: ",
		},
		{
			name:       "Hotfix branch",
			branchName: "hotfix/URGENT-789-critical-fix",
			expected:   "URGENT-789: ",
		},
		{
			name:       "Main branch - no prefix",
			branchName: "main",
			expected:   "",
		},
		{
			name:       "Master branch - no prefix",
			branchName: "master",
			expected:   "",
		},
		{
			name:       "Develop branch - no prefix",
			branchName: "develop",
			expected:   "",
		},
		{
			name:       "Feature branch without ticket",
			branchName: "feature/add-authentication",
			expected:   "",
		},
		{
			name:       "Random branch name",
			branchName: "my-random-branch",
			expected:   "",
		},
		{
			name:       "Branch with numbers but no ticket format",
			branchName: "feature-123-test",
			expected:   "",
		},
		{
			name:       "Empty branch name",
			branchName: "",
			expected:   "",
		},
		{
			name:       "Single letter ticket",
			branchName: "A-1-test",
			expected:   "A-1: ",
		},
		{
			name:       "Long ticket prefix",
			branchName: "VERYLONGPREFIX-12345-description",
			expected:   "VERYLONGPREFIX-12345: ",
		},
		{
			name:       "Nested path with ticket",
			branchName: "feature/epic/STORY-999-implementation",
			expected:   "STORY-999: ",
		},
		{
			name:       "Ticket with lowercase letters",
			branchName: "feat-123-test",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTicketPrefix(tt.branchName)
			if result != tt.expected {
				t.Errorf("extractTicketPrefix(%q) = %q, want %q", tt.branchName, result, tt.expected)
			}
		})
	}
}

func TestAnalyzeAndGenerateMessage(t *testing.T) {
	tests := []struct {
		name     string
		diff     string
		expected string
	}{
		{
			name: "New file added",
			diff: `diff --git a/new-file.go b/new-file.go
new file mode 100644
index 0000000..1234567
--- /dev/null
+++ b/new-file.go
@@ -0,0 +1,5 @@
+package main
+
+func main() {
+    fmt.Println("Hello World")
+}`,
			expected: "Add new-file.go",
		},
		{
			name: "File deleted",
			diff: `diff --git a/old-file.go b/old-file.go
deleted file mode 100644
index 1234567..0000000
--- a/old-file.go
+++ /dev/null
@@ -1,5 +0,0 @@
-package main
-
-func main() {
-    fmt.Println("Hello World")
-}`,
			expected: "Remove old-file.go",
		},
		{
			name: "File modified with more additions",
			diff: `diff --git a/main.go b/main.go
index 1234567..7890abc 100644
--- a/main.go
+++ b/main.go
@@ -1,3 +1,8 @@
 package main
 
+import "fmt"
+
+func newFunction() {
+    fmt.Println("New function")
+}
+
 func main() {`,
			expected: "Enhance main.go",
		},
		{
			name: "File modified with more deletions",
			diff: `diff --git a/main.go b/main.go
index 1234567..7890abc 100644
--- a/main.go
+++ b/main.go
@@ -1,10 +1,3 @@
 package main
 
-import "fmt"
-
-func oldFunction() {
-    fmt.Println("Old function")
-}
-
 func main() {`,
			expected: "Refactor main.go",
		},
		{
			name: "File modified with balanced changes",
			diff: `diff --git a/main.go b/main.go
index 1234567..7890abc 100644
--- a/main.go
+++ b/main.go
@@ -1,5 +1,5 @@
 package main
 
-import "fmt"
+import "log"
 
 func main() {`,
			expected: "Update main.go",
		},
		{
			name: "Multiple files added",
			diff: `diff --git a/file1.go b/file1.go
new file mode 100644
index 0000000..1234567
--- /dev/null
+++ b/file1.go
@@ -0,0 +1,3 @@
+package main
+
+// File 1
diff --git a/file2.go b/file2.go
new file mode 100644
index 0000000..7890abc
--- /dev/null
+++ b/file2.go
@@ -0,0 +1,3 @@
+package main
+
+// File 2`,
			expected: "Add 2 new files",
		},
		{
			name: "Multiple files modified",
			diff: `diff --git a/file1.go b/file1.go
index 1234567..7890abc 100644
--- a/file1.go
+++ b/file1.go
@@ -1,3 +1,4 @@
 package main
 
 // Modified file 1
+// New line
diff --git a/file2.go b/file2.go
index 2345678..8901bcd 100644
--- a/file2.go
+++ b/file2.go
@@ -1,3 +1,4 @@
 package main
 
 // Modified file 2
+// New line`,
			expected: "Update 2 files",
		},
		{
			name:     "Empty diff",
			diff:     "",
			expected: "Update project files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzeAndGenerateMessage(tt.diff)
			if result != tt.expected {
				t.Errorf("analyzeAndGenerateMessage() = %q, want %q", result, tt.expected)
			}
		})
	}
}
