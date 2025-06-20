package storage

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bonyuta0204/personal-agent/go/internal/infrastructure/util"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
)

// logDuration logs the time taken by a function with the given name
func logDuration(start time.Time, name string) {
	duration := time.Since(start)
	log.Printf("%s took %s", name, duration)
}

// GitHubStorage implements the storage.Storage interface for GitHub
// This is an infrastructure layer component
type GitHubStorage struct {
	client     *github.Client
	repoOwner  string
	repoName   string
	tmpDirPath string // Path to the local repository clone
}

// NewGitHubStorage creates a new GitHub storage instance
func NewGitHubStorage(repo string) (*GitHubStorage, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	repoParts := strings.Split(repo, "/")
	if len(repoParts) != 2 {
		return nil, fmt.Errorf("invalid repo format, expected 'owner/repo'")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &GitHubStorage{
		client:    github.NewClient(tc),
		repoOwner: repoParts[0],
		repoName:  repoParts[1],
	}, nil
}

// SaveDocument implements the Storage interface
func (s *GitHubStorage) SaveDocument(document *model.Document) error {
	if document == nil {
		return fmt.Errorf("document cannot be nil")
	}

	ctx := context.Background()

	// Prepare file content
	content := []byte(document.Content)
	message := fmt.Sprintf("Add/update document: %s", document.Path)

	// Check if file exists to determine if this is an update
	_, _, _, err := s.client.Repositories.GetContents(ctx, s.repoOwner, s.repoName, document.Path, nil)
	if err != nil && !strings.Contains(err.Error(), "404 Not Found") {
		return fmt.Errorf("error checking if file exists: %w", err)
	}

	// Create or update file
	_, _, err = s.client.Repositories.CreateFile(ctx, s.repoOwner, s.repoName, document.Path, &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: content,
		SHA:     nil, // Will be populated by GitHub for updates
	})

	if err != nil {
		return fmt.Errorf("error creating/updating file: %w", err)
	}

	return nil
}

// SaveMemory implements the Storage interface
func (s *GitHubStorage) SaveMemory(memory *model.Memory) error {
	if memory == nil {
		return fmt.Errorf("memory cannot be nil")
	}

	// For memories, we'll store them in a .memories directory
	path := filepath.Join(".memories", memory.Path)
	if !strings.HasSuffix(path, ".md") {
		path += ".md"
	}

	ctx := context.Background()

	// Prepare file content
	content := []byte(memory.Content)
	message := fmt.Sprintf("Add/update memory: %s", path)

	// Check if file exists to determine if this is an update
	_, _, _, err := s.client.Repositories.GetContents(ctx, s.repoOwner, s.repoName, path, nil)
	if err != nil && !strings.Contains(err.Error(), "404 Not Found") {
		return fmt.Errorf("error checking if memory file exists: %w", err)
	}

	// Create or update file
	_, _, err = s.client.Repositories.CreateFile(ctx, s.repoOwner, s.repoName, path, &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: content,
		SHA:     nil, // Will be populated by GitHub for updates
	})

	if err != nil {
		return fmt.Errorf("error creating/updating memory file: %w", err)
	}

	return nil
}

// fetchFileContent is a helper method that handles fetching file content either from local file system or GitHub API
func (s *GitHubStorage) fetchFileContent(path string) (content string, modTime time.Time, err error) {

	if s.tmpDirPath == "" {
		if err := s.downloadRepository(); err != nil {
			return "", time.Time{}, fmt.Errorf("error downloading repository: %w", err)
		}
	}
	// If we have a local clone, read from the file system
	fullPath := filepath.Join(s.tmpDirPath, path)
	contentBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error reading file from local clone: %w", err)
	}

	// Check if the file is binary
	if isBinary(contentBytes) {
		return "", time.Time{}, fmt.Errorf("skipping binary file: %s", path)
	}

	// Get file info for modification time
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error getting file info: %w", err)
	}

	return string(contentBytes), fileInfo.ModTime(), nil
}

// isBinary checks if a byte slice contains binary data by looking for null bytes
// and checking for common binary file signatures.
func isBinary(data []byte) bool {
	// Empty files are not binary
	if len(data) == 0 {
		return false
	}

	// Check if it's valid UTF-8 text
	if !utf8.Valid(data) {
		return true
	}

	// Only consider files with .pptx, .xlsx, .docx, .pdf, .jpg, .png, etc. as binary
	// For markdown files with Japanese text, we want to process them as text

	// For image files and other known binary formats, check magic numbers
	if len(data) > 8 {
		// PNG signature
		if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}) {
			return true
		}
		// JPEG signatures
		if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
			return true
		}
		// PDF signature
		if bytes.HasPrefix(data, []byte{0x25, 0x50, 0x44, 0x46}) { // %PDF
			return true
		}
		// ZIP signatures (used by Office files like .docx, .xlsx, .pptx)
		if bytes.HasPrefix(data, []byte{0x50, 0x4B, 0x03, 0x04}) {
			return true
		}
	}

	// For text files with UTF-8 characters (like Japanese text),
	// we'll be more lenient and not consider them binary
	return false
}

// FetchDocument implements the Storage interface
func (s *GitHubStorage) FetchDocument(storeId model.StoreId, path string) (*model.Document, error) {
	content, modTime, err := s.fetchFileContent(path)
	if err != nil {
		return nil, err
	}

	sha := util.CalculateSHA256(content)

	return &model.Document{
		Path:       path,
		StoreId:    storeId,
		Content:    content,
		SHA:        sha,
		ModifiedAt: modTime,
	}, nil
}

// FetchMemory implements the Storage interface
func (s *GitHubStorage) FetchMemory(path string) (*model.Memory, error) {
	// For memories, we'll look in the .memories directory
	if !strings.HasPrefix(path, ".memories/") {
		path = filepath.Join(".memories", path)
	}
	if !strings.HasSuffix(path, ".md") {
		path += ".md"
	}

	// Use fetchFileContent directly to avoid unnecessary document creation
	content, modTime, err := s.fetchFileContent(path)
	if err != nil {
		return nil, err
	}

	sha := util.CalculateSHA256(content)

	// Strip the .memories/ prefix and .md suffix for the memory path
	memoryPath := strings.TrimSuffix(strings.TrimPrefix(path, ".memories/"), ".md")

	return &model.Memory{
		Path:      memoryPath,
		Content:   content,
		SHA:       sha,
		CreatedAt: modTime,
		UpdatedAt: modTime,
	}, nil
}

// GetDocumentEntriesFromFS recursively gets all file paths from the local file system
func (s *GitHubStorage) GetDocumentEntriesFromFS(dir string) ([]model.DocumentEntry, error) {

	var documentEntries []model.DocumentEntry

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			subEntries, err := s.GetDocumentEntriesFromFS(fullPath)
			if err != nil {
				return nil, err
			}
			documentEntries = append(documentEntries, subEntries...)
		} else if entry.Type().IsRegular() {
			// Convert to relative path from repo root
			relPath, err := filepath.Rel(s.tmpDirPath, fullPath)
			if err != nil {
				return nil, fmt.Errorf("error getting relative path: %w", err)
			}
			fileInfo, err := os.Stat(fullPath)
			if err != nil {
				return nil, fmt.Errorf("error getting file info: %w", err)
			}
			documentEntries = append(documentEntries, model.DocumentEntry{
				Path:       relPath,
				ModifiedAt: fileInfo.ModTime(),
			})
		}
	}

	return documentEntries, nil
}

// GetDocumentEntries implements the Storage interface
func (s *GitHubStorage) GetDocumentEntries() ([]model.DocumentEntry, error) {
	if s.tmpDirPath == "" {
		log.Printf("No local clone found, downloading repository...")
		if err := s.downloadRepository(); err != nil {
			return nil, fmt.Errorf("error downloading repository: %w", err)
		}
	} else {
		log.Printf("Using existing local clone at %s", s.tmpDirPath)
	}

	// If we have a local clone, read from the file system
	log.Printf("Getting document entries from local filesystem...")
	paths, err := s.GetDocumentEntriesFromFS(s.tmpDirPath)
	if err != nil {
		return nil, fmt.Errorf("error getting paths from local clone: %w", err)
	}
	log.Printf("Found %d document entries", len(paths))
	return paths, nil
}

// downloadRepository downloads the repository tarball and extracts it to a temporary directory
func (s *GitHubStorage) downloadRepository() error {
	start := time.Now()
	defer logDuration(start, "downloadRepository")

	ctx := context.Background()
	log.Printf("Starting repository download for %s/%s", s.repoOwner, s.repoName)

	// Create a temporary directory to store the downloaded tarball
	tmpDir, err := os.MkdirTemp("", "github-repo-*")
	if err != nil {
		return fmt.Errorf("error creating temp directory: %w", err)
	}

	// Get the tarball URL for the repository
	url, _, err := s.client.Repositories.GetArchiveLink(ctx, s.repoOwner, s.repoName, github.Tarball, &github.RepositoryContentGetOptions{}, 1)
	if err != nil {
		os.RemoveAll(tmpDir) // Clean up temp dir on error
		return fmt.Errorf("error getting archive link: %w", err)
	}

	// Download the tarball
	resp, err := s.client.Client().Get(url.String())
	if err != nil {
		os.RemoveAll(tmpDir) // Clean up temp dir on error
		return fmt.Errorf("error downloading repository: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		os.RemoveAll(tmpDir) // Clean up temp dir on error
		return fmt.Errorf("error downloading repository: %s", resp.Status)
	}

	// Save the tarball to a temporary file
	tarballPath := filepath.Join(tmpDir, "repo.tar.gz")
	file, err := os.Create(tarballPath)
	if err != nil {
		os.RemoveAll(tmpDir) // Clean up temp dir on error
		return fmt.Errorf("error creating tarball file: %w", err)
	}
	defer file.Close()

	_, err = file.ReadFrom(resp.Body)
	if err != nil {
		os.RemoveAll(tmpDir) // Clean up temp dir on error
		return fmt.Errorf("error saving tarball: %w", err)
	}

	// Create a directory to extract the tarball
	extractDir := filepath.Join(tmpDir, "extracted")
	err = os.Mkdir(extractDir, 0755)
	if err != nil {
		os.RemoveAll(tmpDir) // Clean up temp dir on error
		return fmt.Errorf("error creating extraction directory: %w", err)
	}

	// Extract the tarball
	cmd := exec.Command("tar", "-xzf", tarballPath, "-C", extractDir, "--strip-components=1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(tmpDir) // Clean up temp dir on error
		return fmt.Errorf("error extracting tarball: %v\nOutput: %s", err, string(output))
	}

	// Clean up the tarball
	err = os.Remove(tarballPath)
	if err != nil {
		// Non-fatal error, just log it
		fmt.Printf("warning: failed to remove tarball: %v\n", err)
	}

	s.tmpDirPath = filepath.Join(extractDir)
	return nil
}

// GetMemoryEntries implements the Storage interface for memories
func (s *GitHubStorage) GetMemoryEntries() ([]model.MemoryEntry, error) {
	if s.tmpDirPath == "" {
		log.Printf("No local clone found, downloading repository...")
		if err := s.downloadRepository(); err != nil {
			return nil, fmt.Errorf("error downloading repository: %w", err)
		}
	} else {
		log.Printf("Using existing local clone at %s", s.tmpDirPath)
	}

	// Look for memories in the .memories directory
	memoriesDir := filepath.Join(s.tmpDirPath, ".memories")
	
	// Check if .memories directory exists
	if _, err := os.Stat(memoriesDir); os.IsNotExist(err) {
		log.Printf(".memories directory does not exist, returning empty list")
		return []model.MemoryEntry{}, nil
	}

	// Get all markdown files from the .memories directory
	var memoryEntries []model.MemoryEntry
	
	err := filepath.Walk(memoriesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories and non-markdown files
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		
		// Get relative path from .memories directory
		relPath, err := filepath.Rel(memoriesDir, path)
		if err != nil {
			return fmt.Errorf("error getting relative path: %w", err)
		}
		
		// Remove .md extension for the memory path
		memoryPath := strings.TrimSuffix(relPath, ".md")
		
		memoryEntries = append(memoryEntries, model.MemoryEntry{
			Path:       memoryPath,
			ModifiedAt: info.ModTime(),
		})
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("error walking memories directory: %w", err)
	}
	
	log.Printf("Found %d memory entries", len(memoryEntries))
	return memoryEntries, nil
}
