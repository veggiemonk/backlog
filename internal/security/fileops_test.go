package security

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestSecureFileOps_ValidatePath(t *testing.T) {
	fs := afero.NewMemMapFs()
	sfo := NewSecureFileOps(fs)

	tests := []struct {
		name      string
		path      string
		wantError bool
		errorType string
	}{
		{
			name:      "valid path",
			path:      "/home/user/task.md",
			wantError: false,
		},
		{
			name:      "empty path",
			path:      "",
			wantError: true,
			errorType: "path cannot be empty",
		},
		{
			name:      "path with null byte",
			path:      "/home/user\x00/task.md",
			wantError: true,
			errorType: "path contains null bytes",
		},
		{
			name:      "path with directory traversal",
			path:      "/home/user/../../../etc/passwd",
			wantError: true,
			errorType: "path contains directory traversal sequences",
		},
		{
			name:      "path too long",
			path:      "/" + strings.Repeat("a", 4100),
			wantError: true,
			errorType: "path exceeds maximum length",
		},
		{
			name:      "suspicious /proc path",
			path:      "/proc/self/mem",
			wantError: true,
			errorType: "path contains suspicious pattern",
		},
		{
			name:      "suspicious /sys path",
			path:      "/sys/kernel/debug",
			wantError: true,
			errorType: "path contains suspicious pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sfo.ValidatePath(tt.path)
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidatePath() expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorType) {
					t.Errorf("ValidatePath() error = %v, want error containing %v", err, tt.errorType)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePath() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSecureFileOps_SecureWrite(t *testing.T) {
	fs := afero.NewMemMapFs()
	sfo := NewSecureFileOps(fs)

	tests := []struct {
		name      string
		path      string
		content   []byte
		perm      os.FileMode
		wantError bool
	}{
		{
			name:      "valid write",
			path:      "/test/file.txt",
			content:   []byte("test content"),
			perm:      0644,
			wantError: false,
		},
		{
			name:      "invalid path",
			path:      "",
			content:   []byte("test"),
			perm:      0644,
			wantError: true,
		},
		{
			name:      "content too large",
			path:      "/test/large.txt",
			content:   make([]byte, 101*1024*1024), // 101MB
			perm:      0644,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sfo.SecureWrite(tt.path, tt.content, tt.perm)
			if tt.wantError {
				if err == nil {
					t.Errorf("SecureWrite() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("SecureWrite() unexpected error: %v", err)
				}
				// Verify file was written
				if exists, _ := afero.Exists(fs, tt.path); !exists {
					t.Errorf("SecureWrite() file was not created")
				}
			}
		})
	}
}

func TestSecureFileOps_SecureRead(t *testing.T) {
	fs := afero.NewMemMapFs()
	sfo := NewSecureFileOps(fs)

	// Create a test file
	testPath := "/test/file.txt"
	testContent := []byte("test content")
	afero.WriteFile(fs, testPath, testContent, 0644)

	// Create a directory (non-regular file)
	dirPath := "/test/dir"
	fs.MkdirAll(dirPath, 0755)

	tests := []struct {
		name      string
		path      string
		wantError bool
		expected  []byte
	}{
		{
			name:      "valid read",
			path:      testPath,
			wantError: false,
			expected:  testContent,
		},
		{
			name:      "invalid path",
			path:      "",
			wantError: true,
		},
		{
			name:      "non-existent file",
			path:      "/test/nonexistent.txt",
			wantError: true,
		},
		{
			name:      "directory (non-regular file)",
			path:      dirPath,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := sfo.SecureRead(tt.path)
			if tt.wantError {
				if err == nil {
					t.Errorf("SecureRead() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("SecureRead() unexpected error: %v", err)
				}
				if string(content) != string(tt.expected) {
					t.Errorf("SecureRead() content = %v, want %v", string(content), string(tt.expected))
				}
			}
		})
	}
}

func TestSecureFileOps_SecureDelete(t *testing.T) {
	fs := afero.NewMemMapFs()
	sfo := NewSecureFileOps(fs)

	// Create test files
	testPath := "/test/file.txt"
	afero.WriteFile(fs, testPath, []byte("test content"), 0644)

	dirPath := "/test/dir"
	fs.MkdirAll(dirPath, 0755)

	tests := []struct {
		name      string
		path      string
		wantError bool
	}{
		{
			name:      "valid delete",
			path:      testPath,
			wantError: false,
		},
		{
			name:      "invalid path",
			path:      "",
			wantError: true,
		},
		{
			name:      "non-existent file",
			path:      "/test/nonexistent.txt",
			wantError: false, // Should not error for non-existent files
		},
		{
			name:      "directory (non-regular file)",
			path:      dirPath,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sfo.SecureDelete(tt.path)
			if tt.wantError {
				if err == nil {
					t.Errorf("SecureDelete() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("SecureDelete() unexpected error: %v", err)
				}
				// For valid deletes, verify file no longer exists
				if tt.path == testPath {
					if exists, _ := afero.Exists(fs, tt.path); exists {
						t.Errorf("SecureDelete() file still exists after deletion")
					}
				}
			}
		})
	}
}

func TestSecureFileOps_SecureMove(t *testing.T) {
	fs := afero.NewMemMapFs()
	sfo := NewSecureFileOps(fs)

	// Create test files
	srcPath := "/test/src.txt"
	dstPath := "/test/dst.txt"
	testContent := []byte("test content")
	afero.WriteFile(fs, srcPath, testContent, 0644)

	tests := []struct {
		name      string
		srcPath   string
		dstPath   string
		wantError bool
	}{
		{
			name:      "valid move",
			srcPath:   srcPath,
			dstPath:   dstPath,
			wantError: false,
		},
		{
			name:      "invalid source path",
			srcPath:   "",
			dstPath:   "/test/valid.txt",
			wantError: true,
		},
		{
			name:      "invalid destination path",
			srcPath:   "/test/another.txt",
			dstPath:   "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For the valid move test, recreate the source file
			if tt.name == "valid move" {
				afero.WriteFile(fs, tt.srcPath, testContent, 0644)
			}

			err := sfo.SecureMove(tt.srcPath, tt.dstPath)
			if tt.wantError {
				if err == nil {
					t.Errorf("SecureMove() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("SecureMove() unexpected error: %v", err)
				}
				// Verify move was successful
				if exists, _ := afero.Exists(fs, tt.srcPath); exists {
					t.Errorf("SecureMove() source file still exists after move")
				}
				if exists, _ := afero.Exists(fs, tt.dstPath); !exists {
					t.Errorf("SecureMove() destination file does not exist after move")
				}
			}
		})
	}
}

func TestSecureFileOps_CreateSecureDirectory(t *testing.T) {
	fs := afero.NewMemMapFs()
	sfo := NewSecureFileOps(fs)

	tests := []struct {
		name      string
		path      string
		perm      os.FileMode
		wantError bool
	}{
		{
			name:      "valid directory creation",
			path:      "/test/secure/dir",
			perm:      0750,
			wantError: false,
		},
		{
			name:      "invalid path",
			path:      "",
			wantError: true,
		},
		{
			name:      "too permissive permissions (should be restricted)",
			path:      "/test/permissive",
			perm:      0777,
			wantError: false, // Should succeed but with restricted permissions
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sfo.CreateSecureDirectory(tt.path, tt.perm)
			if tt.wantError {
				if err == nil {
					t.Errorf("CreateSecureDirectory() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("CreateSecureDirectory() unexpected error: %v", err)
				}
				// Verify directory was created
				if exists, _ := afero.DirExists(fs, tt.path); !exists {
					t.Errorf("CreateSecureDirectory() directory was not created")
				}
			}
		})
	}
}

func TestSecureFileOps_CheckFileIntegrity(t *testing.T) {
	fs := afero.NewMemMapFs()
	sfo := NewSecureFileOps(fs)

	// Create test files
	testPath := "/test/file.txt"
	afero.WriteFile(fs, testPath, []byte("test content"), 0644)

	dirPath := "/test/dir"
	fs.MkdirAll(dirPath, 0755)

	tests := []struct {
		name      string
		path      string
		wantError bool
	}{
		{
			name:      "valid file",
			path:      testPath,
			wantError: false,
		},
		{
			name:      "invalid path",
			path:      "",
			wantError: true,
		},
		{
			name:      "non-existent file",
			path:      "/test/nonexistent.txt",
			wantError: true,
		},
		{
			name:      "directory (non-regular file)",
			path:      dirPath,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sfo.CheckFileIntegrity(tt.path)
			if tt.wantError {
				if err == nil {
					t.Errorf("CheckFileIntegrity() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("CheckFileIntegrity() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSecurityError_Error(t *testing.T) {
	err := SecurityError{
		Operation: "test",
		Path:      "/test/path",
		Reason:    "test reason",
	}

	expected := "security error in test operation on '/test/path': test reason"
	if err.Error() != expected {
		t.Errorf("SecurityError.Error() = %v, want %v", err.Error(), expected)
	}
}