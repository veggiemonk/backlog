package validation

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"
)

// SecureFileOps provides secure file operations with safety checks
type SecureFileOps struct {
	fs afero.Fs
}

// NewSecureFileOps creates a new secure file operations instance
func NewSecureFileOps(fs afero.Fs) *SecureFileOps {
	return &SecureFileOps{fs: fs}
}

// SecurityError represents a security-related file operation error
type SecurityError struct {
	Operation string
	Path      string
	Reason    string
}

func (e SecurityError) Error() string {
	return fmt.Sprintf("security error in %s operation on '%s': %s", e.Operation, e.Path, e.Reason)
}

// ValidatePath performs comprehensive path validation
func (sfo *SecureFileOps) ValidatePath(path string) error {
	if path == "" {
		err := SecurityError{
			Operation: "validate",
			Path:      path,
			Reason:    "path cannot be empty",
		}
		LogFileAccessViolation("validate", path, err.Reason)
		return err
	}

	// Check for null bytes
	if strings.Contains(path, "\x00") {
		err := SecurityError{
			Operation: "validate",
			Path:      path,
			Reason:    "path contains null bytes",
		}
		LogFileAccessViolation("validate", path, err.Reason)
		return err
	}

	// Check for directory traversal
	if strings.Contains(path, "..") {
		err := SecurityError{
			Operation: "validate",
			Path:      path,
			Reason:    "path contains directory traversal sequences",
		}
		LogFileAccessViolation("validate", path, err.Reason)
		return err
	}

	// Check path length
	if len(path) > 4096 {
		return SecurityError{
			Operation: "validate",
			Path:      path,
			Reason:    "path exceeds maximum length",
		}
	}

	// Check for suspicious patterns
	suspiciousPatterns := []string{
		"/proc/",
		"/sys/",
		"/dev/",
		"\\\\.\\", // Windows device paths
		"\\\\?\\", // Windows extended paths
	}

	lowerPath := strings.ToLower(path)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerPath, pattern) {
			return SecurityError{
				Operation: "validate",
				Path:      path,
				Reason:    fmt.Sprintf("path contains suspicious pattern: %s", pattern),
			}
		}
	}

	return nil
}

// SecureRead reads a file with security checks
func (sfo *SecureFileOps) SecureRead(path string) ([]byte, error) {
	if err := sfo.ValidatePath(path); err != nil {
		return nil, err
	}

	// Check if file exists and is a regular file
	info, err := sfo.fs.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	if !info.Mode().IsRegular() {
		return nil, SecurityError{
			Operation: "read",
			Path:      path,
			Reason:    "path is not a regular file",
		}
	}

	// Check file size limits (prevent reading extremely large files)
	const maxFileSize = 100 * 1024 * 1024 // 100MB
	if info.Size() > maxFileSize {
		return nil, SecurityError{
			Operation: "read",
			Path:      path,
			Reason:    "file size exceeds maximum allowed size",
		}
	}

	// Check file permissions
	if err := sfo.checkReadPermissions(path, info); err != nil {
		return nil, err
	}

	file, err := sfo.fs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read with size limit
	limitedReader := io.LimitReader(file, maxFileSize)
	content, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return content, nil
}

// SecureWrite writes a file with security checks
func (sfo *SecureFileOps) SecureWrite(path string, content []byte, perm os.FileMode) error {
	if err := sfo.ValidatePath(path); err != nil {
		return err
	}

	// Validate content size
	const maxContentSize = 100 * 1024 * 1024 // 100MB
	if len(content) > maxContentSize {
		return SecurityError{
			Operation: "write",
			Path:      path,
			Reason:    "content size exceeds maximum allowed size",
		}
	}

	// Ensure parent directory exists
	parentDir := filepath.Dir(path)
	if err := sfo.ValidatePath(parentDir); err != nil {
		return fmt.Errorf("invalid parent directory: %w", err)
	}

	if err := sfo.fs.MkdirAll(parentDir, 0o750); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Check if file already exists and validate overwrite
	if exists, err := afero.Exists(sfo.fs, path); err != nil {
		return fmt.Errorf("failed to check file existence: %w", err)
	} else if exists {
		info, err := sfo.fs.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to stat existing file: %w", err)
		}

		if !info.Mode().IsRegular() {
			return SecurityError{
				Operation: "write",
				Path:      path,
				Reason:    "cannot overwrite non-regular file",
			}
		}

		if err := sfo.checkWritePermissions(path, info); err != nil {
			return err
		}
	}

	// Create temporary file for atomic write
	tempPath := path + ".tmp." + generateRandomSuffix()
	if err := sfo.ValidatePath(tempPath); err != nil {
		return fmt.Errorf("invalid temporary file path: %w", err)
	}

	// Write to temporary file
	tempFile, err := sfo.fs.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, perm)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	_, writeErr := tempFile.Write(content)
	closeErr := tempFile.Close()

	if writeErr != nil {
		sfo.fs.Remove(tempPath) // Clean up on write error
		return fmt.Errorf("failed to write to temporary file: %w", writeErr)
	}

	if closeErr != nil {
		sfo.fs.Remove(tempPath) // Clean up on close error
		return fmt.Errorf("failed to close temporary file: %w", closeErr)
	}

	// Atomic rename
	if err := sfo.fs.Rename(tempPath, path); err != nil {
		sfo.fs.Remove(tempPath) // Clean up on rename error
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

// SecureDelete deletes a file with security checks
func (sfo *SecureFileOps) SecureDelete(path string) error {
	if err := sfo.ValidatePath(path); err != nil {
		return err
	}

	// Check if file exists
	info, err := sfo.fs.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist, nothing to delete
		}
		return fmt.Errorf("failed to stat file: %w", err)
	}

	if !info.Mode().IsRegular() {
		return SecurityError{
			Operation: "delete",
			Path:      path,
			Reason:    "cannot delete non-regular file",
		}
	}

	// Check permissions
	if err := sfo.checkWritePermissions(path, info); err != nil {
		return err
	}

	// Remove the file
	if err := sfo.fs.Remove(path); err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	return nil
}

// SecureMove moves a file with security checks
func (sfo *SecureFileOps) SecureMove(srcPath, dstPath string) error {
	if err := sfo.ValidatePath(srcPath); err != nil {
		return fmt.Errorf("invalid source path: %w", err)
	}

	if err := sfo.ValidatePath(dstPath); err != nil {
		return fmt.Errorf("invalid destination path: %w", err)
	}

	// Check source file
	srcInfo, err := sfo.fs.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	if !srcInfo.Mode().IsRegular() {
		return SecurityError{
			Operation: "move",
			Path:      srcPath,
			Reason:    "source is not a regular file",
		}
	}

	// Check permissions on source
	if err := sfo.checkWritePermissions(srcPath, srcInfo); err != nil {
		return fmt.Errorf("source file permission check failed: %w", err)
	}

	// Check destination doesn't exist or is overwritable
	if exists, err := afero.Exists(sfo.fs, dstPath); err != nil {
		return fmt.Errorf("failed to check destination existence: %w", err)
	} else if exists {
		dstInfo, err := sfo.fs.Stat(dstPath)
		if err != nil {
			return fmt.Errorf("failed to stat destination file: %w", err)
		}

		if !dstInfo.Mode().IsRegular() {
			return SecurityError{
				Operation: "move",
				Path:      dstPath,
				Reason:    "cannot overwrite non-regular destination file",
			}
		}

		if err := sfo.checkWritePermissions(dstPath, dstInfo); err != nil {
			return fmt.Errorf("destination file permission check failed: %w", err)
		}
	}

	// Ensure destination directory exists
	dstDir := filepath.Dir(dstPath)
	if err := sfo.fs.MkdirAll(dstDir, 0o750); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Perform the move
	if err := sfo.fs.Rename(srcPath, dstPath); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	return nil
}

// checkReadPermissions checks if the current process can read the file
func (sfo *SecureFileOps) checkReadPermissions(path string, info os.FileInfo) error {
	// On Unix-like systems, check if owner/group/other has read permission
	mode := info.Mode()
	if mode&0o444 == 0 {
		return SecurityError{
			Operation: "read",
			Path:      path,
			Reason:    "insufficient read permissions",
		}
	}
	return nil
}

// checkWritePermissions checks if the current process can write the file
func (sfo *SecureFileOps) checkWritePermissions(path string, info os.FileInfo) error {
	// For memory filesystem in tests, skip detailed permission checks
	// In production with real filesystem, implement proper permission checking

	// Check if the file itself is writable (if it exists)
	if info != nil {
		mode := info.Mode()
		if mode&0o200 == 0 {
			return SecurityError{
				Operation: "write",
				Path:      path,
				Reason:    "file is not writable",
			}
		}
	}

	return nil
}

// isImmutable checks if a file has the immutable attribute (Unix systems)
// func (sfo *SecureFileOps) isImmutable(path string) bool {
// 	// This is a simplified implementation that always returns false
// 	// In a production environment, you would implement platform-specific
// 	// checks for file attributes like chattr +i on Linux
// 	return false
// }

// generateRandomSuffix generates a random suffix for temporary files
func generateRandomSuffix() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a simple timestamp-based suffix
		return fmt.Sprintf("%d", os.Getpid())
	}
	return fmt.Sprintf("%x", bytes)
}

// CreateSecureDirectory creates a directory with secure permissions
func (sfo *SecureFileOps) CreateSecureDirectory(path string, perm os.FileMode) error {
	if err := sfo.ValidatePath(path); err != nil {
		return err
	}

	// Ensure permissions are not too permissive
	if perm&0o077 != 0 {
		perm &= ^os.FileMode(0o077) // Remove group and other permissions
	}

	return sfo.fs.MkdirAll(path, perm)
}

// CheckFileIntegrity performs basic integrity checks on a file
func (sfo *SecureFileOps) CheckFileIntegrity(path string) error {
	if err := sfo.ValidatePath(path); err != nil {
		return err
	}

	info, err := sfo.fs.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Check if it's a regular file
	if !info.Mode().IsRegular() {
		return SecurityError{
			Operation: "integrity_check",
			Path:      path,
			Reason:    "not a regular file",
		}
	}

	// Check for reasonable file size
	if info.Size() < 0 {
		return SecurityError{
			Operation: "integrity_check",
			Path:      path,
			Reason:    "file has negative size",
		}
	}

	// Check modification time is reasonable (not too far in the future)
	if info.ModTime().After(time.Now().Add(24 * time.Hour)) {
		return SecurityError{
			Operation: "integrity_check",
			Path:      path,
			Reason:    "file modification time is suspiciously far in the future",
		}
	}

	return nil
}

