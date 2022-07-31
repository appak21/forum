package models

import (
	"database/sql"
)

func CreateUser(user User) error { //✅✅✅
	stmt := `INSERT INTO "main"."users"(
		"username",
		"email",
		"password",
		"created_at"
	) VALUES (?, ?, ?, ?)`

	_, err := db.Exec(stmt, user.Username, user.Email, user.Password, user.CreatedAt)
	return err
}

func CreatePost(post *Post) error { //✅✅✅
	stmt := `INSERT INTO "main"."posts"
	("user_id", "username", "title", "text", "created_at")
	VALUES (?, ?, ?, ?, ?);`
	res, err := db.Exec(stmt, post.UserID, post.Username, post.Title, post.Text, post.CreatedAt)
	if err != nil {
		return err
	}
	post.ID, err = res.LastInsertId()
	return err
}

func CreateTags(postID int64, tags []string) error { //✅✅✅
	stmt1 := `INSERT INTO "main"."tags" (Tag) VALUES (?)`
	stmt2 := `INSERT INTO "main"."posts_and_tags"
	("post_id", "tag_id")
	VALUES (?, ?);`

	for _, tag := range tags {
		res, err := db.Exec(stmt1, tag)
		if err != nil {
			return err
		}
		tagID, err := res.LastInsertId()
		if err != nil {
			return err
		}
		_, err = db.Exec(stmt2, postID, tagID)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateComment(cmt Comment) error { //✅✅✅
	stmt := `INSERT INTO "main"."comments"
	("user_id", "post_id", "username", "text", "created_at")
	VALUES (?, ?, ?, ?, ?);`
	_, err := db.Exec(stmt, cmt.UserID, cmt.PostID, cmt.Username, cmt.Text, cmt.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func VotePost(userID, postID int64, vote int) error { //✅✅✅
	stmtSelect := `SELECT id, vote FROM post_votes WHERE user_id = ? AND post_id = ?;`
	stmtExec := `INSERT INTO "main"."post_votes" (
		"user_id",
		"post_id",
		"vote")
		VALUES (?, ?, ?)`
	stmtDelete := `DELETE FROM "main"."post_votes" WHERE "id" = ?`
	var id, like int64
	row := db.QueryRow(stmtSelect, userID, postID)
	err := row.Scan(&id, &like)
	if err == sql.ErrNoRows {
		if _, err = db.Exec(stmtExec, userID, postID, vote); err != nil {
			return err
		}
		return nil
	}
	if _, err = db.Exec(stmtDelete, id); err != nil {
		return err
	}
	if int64(vote) != like {
		if _, err = db.Exec(stmtExec, userID, postID, vote); err != nil {
			return err
		}
	}
	return nil
}

func VoteComment(userID, cmtID int64, vote int) error { //✅✅✅
	stmtSelect := `SELECT id, vote FROM comment_votes WHERE user_id = ? AND comment_id = ?;`
	stmtExec := `INSERT INTO "main"."comment_votes" (
		"user_id",
		"comment_id",
		"vote")
		VALUES (?, ?, ?)`
	stmtDelete := `DELETE FROM "main"."comment_votes" WHERE "id" = ?`

	var id, like int64
	row := db.QueryRow(stmtSelect, userID, cmtID)
	err := row.Scan(&id, &like)
	if err != nil {
		return err
	}
	if _, err = db.Exec(stmtExec, userID, cmtID, vote); err != nil {
		return err
	}

	if _, err = db.Exec(stmtDelete, id); err != nil {
		return err
	}

	if int64(vote) != like {
		if _, err = db.Exec(stmtExec, userID, cmtID, vote); err != nil {
			return err
		}
	}
	return nil
}
