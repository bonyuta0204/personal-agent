package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	port "github.com/bonyuta0204/personal-agent/go/internal/domain/port/storage"
)

// GitHubStorage implements the storage.Storage interface for GitHub
// This is an infrastructure layer component
type GitHubStorage struct {
	client    *github.Client
	repoOwner string
	repoName  string
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

	doc := &model.Document{
		Path:    path,
		Content: memory.Content,
	}

	return s.SaveDocument(doc)
}

// FetchDocument implements the Storage interface
func (s *GitHubStorage) FetchDocument(path string) (*model.Document, error) {
	ctx := context.Background()

	fileContent, _, _, err := s.client.Repositories.GetContents(ctx, s.repoOwner, s.repoName, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching file: %w", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, fmt.Errorf("error getting file content: %w", err)
	}

	// Get commit info to get timestamps
	commits, _, err := s.client.Repositories.ListCommits(ctx, s.repoOwner, s.repoName, &github.CommitsListOptions{
		Path: path,
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error getting commit info: %w", err)
	}

	var createdAt, updatedAt time.Time
	if len(commits) > 0 {
		updatedAt = commits[0].Commit.Committer.GetDate().Time
		createdAt = updatedAt // Default to same as updatedAt

		// Try to find the first commit to get creation time
		firstCommit, _, err := s.client.Repositories.GetCommit(ctx, s.repoOwner, s.repoName, *fileContent.SHA, nil)
		if err == nil && firstCommit != nil && firstCommit.Commit != nil && firstCommit.Commit.Committer != nil {
			createdAt = firstCommit.Commit.Committer.GetDate().Time
		}
	}

	return &model.Document{
		Path:      path,
		Content:   content,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
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

	doc, err := s.FetchDocument(path)
	if err != nil {
		return nil, err
	}

	// Strip the .memories/ prefix and .md suffix for the memory path
	memoryPath := strings.TrimSuffix(strings.TrimPrefix(doc.Path, ".memories/"), ".md")

	return &model.Memory{
		Path:      memoryPath,
		Content:   doc.Content,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}, nil
}

// GetAllPaths implements the Storage interface
func (s *GitHubStorage) GetAllPaths() ([]string, error) {
	ctx := context.Background()
	var allPaths []string

	// Get all files in the repository
	_, dirContents, _, err := s.client.Repositories.GetContents(ctx, s.repoOwner, s.repoName, ".", nil)
	if err != nil {
		return nil, fmt.Errorf("error getting repository contents: %w", err)
	}

	// Recursively get all file paths
	var traverseDir func(path string) error
	traverseDir = func(path string) error {
		_, dirContents, _, err := s.client.Repositories.GetContents(ctx, s.repoOwner, s.repoName, path, nil)
		if err != nil {
			return fmt.Errorf("error getting contents of %s: %w", path, err)
		}

		for _, content := range dirContents {
			if content.GetType() == "dir" {
				if err := traverseDir(content.GetPath()); err != nil {
					return err
				}
			} else if content.GetType() == "file" {
				allPaths = append(allPaths, content.GetPath())
			}
		}
		return nil
	}

	// Start traversal from root
	for _, content := range dirContents {
		if content.GetType() == "dir" {
			if err := traverseDir(content.GetPath()); err != nil {
				return nil, err
			}
		} else if content.GetType() == "file" {
			allPaths = append(allPaths, content.GetPath())
		}
	}

	return allPaths, nil
}

// GitHubStorageFactory implements the StorageFactory interface for GitHub
type GitHubStorageFactory struct{}

// NewGitHubStorageFactory creates a new GitHub storage factory
func NewGitHubStorageFactory() *GitHubStorageFactory {
	return &GitHubStorageFactory{}
}

// CreateStorage creates a new GitHub storage instance
func (f *GitHubStorageFactory) CreateStorage(store model.DocumentStore) (port.Storage, error) {
	if store.Type() != model.StoreTypeGitHub {
		return nil, fmt.Errorf("unsupported store type: %s", store.Type())
	}

	// Type assert to GitHubStore to access GitHub-specific fields
	githubStore, ok := store.(*model.GitHubStore)
	if !ok {
		return nil, fmt.Errorf("invalid store type for GitHub")
	}

	return NewGitHubStorage(githubStore.Repo())
}
