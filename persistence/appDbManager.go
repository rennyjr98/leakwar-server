package persistence

import (
	"fmt"
)

// AppDbManager : DB manager operations
type AppDbManager struct {
	Database int
}

func (manager *AppDbManager) init() {
	fmt.Println("FFF")
}

// Add : Insert entity to database
func (manager *AppDbManager) Add(entity interface{}) {
	fmt.Println("FFF")
}

// Remove : Remove virtually data
func (manager *AppDbManager) Remove(entity interface{}) {
	fmt.Println("FFF")
}

// Delete : Remove physsically data
func (manager *AppDbManager) Delete(entity interface{}) {
	fmt.Println("FFF")
}

// Get : Select one entity from database
func (manager *AppDbManager) Get() interface{} {
	return nil
}
