package repository

import "github.com/bonyuta0204/personal-agent/go/internal/domain/model"

type MemoryRepository interface {
	SaveMemory(memory *model.Memory) error
	ListMemories() ([]*model.Memory, error)
	FindExistingSHAs(memories []*model.Memory) ([]string, error)
}
