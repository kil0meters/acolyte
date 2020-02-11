package forum

import (
	"time"
)

type Comment struct {
	ID        string     `valid:"printableascii,required"`
	Account   string     `valid:"printableascii,required"`
	Body      string     `valid:"type(string),optional"`
	Removed   bool       `valid:"-"`
	CreatedAt time.Time  `valid:"-"`
	Upvotes   int        `valid:"-"`
	Downvotes int        `valid:"-"`
	Children  []*Comment `valid:"-"`
}
