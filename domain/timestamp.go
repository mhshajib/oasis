package domain

import "time"

// TimeStamp contains createdAt, updatedAt and deletedAt fields
type TimeStamp struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// PopulateCreateTimeStamp populate the timestamp for created_at and updated_at field
func (t *TimeStamp) PopulateCreateTimeStamp() {
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
}

// PopulateUpdateTimeStamp populate the timestamp for updated_at field
func (t *TimeStamp) PopulateUpdateTimeStamp() {
	t.UpdatedAt = time.Now()
}

// PopulateDeleteTimeStamp populate the timestamp for deleted_at field
func (t *TimeStamp) PopulateDeleteTimeStamp() {
	now := time.Now()
	t.DeletedAt = &now
}
