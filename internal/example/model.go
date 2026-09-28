package example

import "time"

// ExampleItem reprezentē datu bāzes modeli vai entītiju.
type ExampleItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName norāda tabulas nosaukumu datubāzē.
func (ExampleItem) TableName() string {
	return "example_items"
}

// CreateExampleDTO ir datu pārsūtīšanas objekts (DTO) jauna ieraksta izveidei no HTTP pieprasījuma.
type CreateExampleDTO struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}
