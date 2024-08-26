
package transformer

import (
	"oasis/domain"
)

//  response body
type  struct {
	ID         	   string  `json:"_id,omitempty"` 
	FieldOne       string              `json:"field_one,omitempty"` 
	TimeStamp   domain.TimeStamp    `json:"timestamp"` 
}

// Transform ...
func Transform(t *domain.) * {
	return &{
		ID:        t.ID,
		FieldOne:  t.FieldOne,
		TimeStamp: t.TimeStamp,
	}
}

// TransformList ...
func TransformList(tl []*domain.) []* {
	 := make([]*, 0)
	for _, t := range tl {
		 = append(, Transform(t))
	}
	return 
}
