package models

// GPage : Struct for JSON page of items from google photos
type GPage struct {
	MediaItems   []*GPhoto `json:"mediaItems"`
	HasPageToken bool
	PageToken    string `json:"nextPageToken"`
	UserID       string
}
