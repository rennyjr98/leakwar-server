package persistence

import (
	"fmt"
	"leakwarsvr/domain"
)

// AppDbManager : DB manager operations
type AppDbManager struct {
	Database int
}

func (manager *AppDbManager) init() {

}

// Add : Insert entity to database
func (manager *AppDbManager) Add(entity interface) {

}

// Remove : Remove virtually data
func (manager *AppDbManager) Remove(entity interface) {

}

// Delete : Remove physsically data
func (manager *AppDbManager) Delete(entity interface) {

}

// Get : Select one entity from database
func (manager *AppDbManager) Get() interface {

}

// GetPage : Select a page entities from database
func (manager *AppDbManager) GetPage(page domain.PaginatedRequest) domain.PaginatedResult {
	return nil
}