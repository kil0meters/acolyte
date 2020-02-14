package forum

import (
	"github.com/kil0meters/acolyte/pkg/authorization"
	"github.com/kil0meters/acolyte/pkg/database"
	"log"
	"time"
)

type Comment struct {
	ID              string                 `db:"comment_id" valid:"printableascii,required"`
	Account         *authorization.Account `db:"-"          valid:"-"`
	AccountID       string                 `db:"account_id" valid:"-"`
	Username        string                 `db:"username"   valid:"-"`
	ParentID        string                 `db:"parent_id"  valid:"-"`
	PostID          string                 `db:"post_id"    valid:"_"`
	Body            string                 `db:"body"       valid:"type(string),optional"`
	CreatedAt       time.Time              `db:"created_at" valid:"-"`
	Removed         bool                   `db:"removed"    valid:"-"`
	Upvotes         int                    `db:"upvotes"    valid:"-"`
	Downvotes       int                    `db:"downvotes"  valid:"-"`
	HasMoreChildren bool                   `db:"-"          valid:"-"`
	Replies         []*Comment             `db:"-"          valid:"-"`
}

func CreateComment(account *authorization.Account, parentID string, postID string, body string) (string, error) {
	commentID := authorization.GenerateID("c", 6)

	_, err := database.DB.Exec("INSERT INTO comments (comment_id, parent_id, post_id, account_id, username, body) VALUES ($1, $2, $3, $4, $5, $6)", commentID, parentID, postID, account.ID, account.Username, body)
	if err != nil {
		return "", err
	}

	return commentID, nil
}

func GetComment(commentID string, depth int) *Comment {
	if depth == 0 {
		return nil
	}

	comment := new(Comment)

	err := database.DB.QueryRowx("SELECT * FROM comments WHERE comment_id = $1", commentID).StructScan(comment)
	if err != nil {
		log.Println(err)
		return nil
	}

	comment.Replies = GetCommentChildren(commentID, depth)

	return comment
}

func GetCommentChildren(commentID string, depth int) []*Comment {
	log.Println("Getting children of", commentID, "at depth", depth)

	if depth <= 0 {
		return nil
	}

	rows, err := database.DB.Queryx("SELECT * FROM comments WHERE parent_id = $1", commentID)
	if err != nil {
		log.Println(err)
		return nil
	}

	comments := make([]*Comment, 0)

	for rows.Next() {
		comment := new(Comment)
		comment.HasMoreChildren = false

		err = rows.StructScan(comment)
		if err != nil {
			log.Println(err)
			return nil
		}

		if depth != 1 {
			comment.Replies = GetCommentChildren(comment.ID, depth-1)
		} else {
			comment.HasMoreChildren = true
		}

		comments = append(comments, comment)
	}

	return comments
}
