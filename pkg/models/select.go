package models

func GetUser(username string) (*User, error) { //✅✅✅
	stmt := `SELECT * FROM "main"."users" WHERE "username" = ?`
	row := db.QueryRow(stmt, username)
	user := &User{}
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

func GetAllTags() ([]string, error) {
	stmt := `SELECT tag FROM tags`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	var tag, tags = "", []string{}
	for rows.Next() {
		if err = rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func GetAllPosts() (*[]Post, error) { //✅✅✅
	stmt := `SELECT * FROM "main"."posts" ORDER BY "created_at" DESC`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	post, posts := Post{}, []Post{}
	for rows.Next() {
		err = rows.Scan(&post.ID, &post.UserID, &post.Username, &post.Title, &post.Text, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return &posts, nil
}

func GetPostById(postID int64) (*Post, error) { //✅✅✅
	stmt := `SELECT * FROM "main"."posts" WHERE "id" = ?`
	row := db.QueryRow(stmt, postID)
	var post Post
	err := row.Scan(&post.ID, &post.UserID, &post.Username, &post.Title, &post.Text, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	post.Tags, err = getPostTags(postID)
	if err != nil {
		return nil, err
	}
	post.Comments, err = getPostComments(postID)
	if err != nil {
		return nil, err
	}
	votes, err := getPostVotes(postID)
	post.Votes = *votes
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func GetPostsCreatedByUser(userID int64) (*[]Post, error) {
	stmt := `
	SELECT posts.id, posts.username, posts.title, posts.text, posts.created_at
	FROM posts
	WHERE posts.user_id = ?`
	rows, err := db.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	post, posts := Post{}, []Post{}
	for rows.Next() {
		err = rows.Scan(&post.ID, &post.Username, &post.Title, &post.Text, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		votes, err := getPostVotes(post.ID)
		if err != nil {
			return nil, err
		}
		post.Votes = *votes
		posts = append(posts, post)
	}
	return &posts, nil
}

func GetPostsVotedByUser(userID int64, vote int) (*[]Post, error) { //✅✅✅
	stmt := `
	SELECT posts.id, posts.username, posts.title, posts.text, posts.created_at
	FROM posts
	INNER JOIN post_votes
	ON posts.id=post_votes.post_id
	WHERE post_votes.vote = ? AND post_votes.user_id = ?`

	rows, err := db.Query(stmt, vote, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	post, posts := Post{}, []Post{}

	for rows.Next() {
		err = rows.Scan(&post.ID, &post.Username, &post.Title, &post.Text, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return &posts, nil
}

func GetCommentByID(cmtID int64) (*Comment, error) { //✅✅✅
	stmt := `SELECT * FROM comments WHERE id = ?`
	var cmt Comment
	row := db.QueryRow(stmt, cmtID)
	err := row.Scan(&cmt.ID, &cmt.PostID, &cmt.UserID, &cmt.Username, &cmt.Text, &cmt.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &cmt, nil
}

func GetPostsByTag(tag string) (*[]Post, error) { //???
	tagID, err := getTagID(tag)
	if err != nil {
		return nil, err
	}
	stmt1 := `SELECT "post_id" FROM "main"."posts_and_tags" WHERE "tag_id" = ?`
	stmt2 := `SELECT * FROM "main"."posts" WHERE "id" = ?`
	var postID int64
	post, posts := Post{}, []Post{}
	rows, err := db.Query(stmt1, tagID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err = rows.Scan(&postID); err != nil {
			return nil, err
		}
		row := db.QueryRow(stmt2, postID)
		if err = row.Scan(&post.ID, &post.UserID, &post.Username, &post.Title, &post.Text, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return &posts, nil
}

//////////////////////////////////////////////////
func getTagID(tag string) (int64, error) { //???
	var tagID int64
	stmt := `SELECT "id" FROM "main"."tags" WHERE "tag" = ?`
	row := db.QueryRow(stmt, tag)
	if err := row.Scan(&tagID); err != nil {
		return 0, err
	}
	return tagID, nil
}
func getPostTags(postID int64) ([]string, error) {
	stmt1 := `SELECT "tag_id" FROM "main"."posts_and_tags" WHERE "post_id" = ?`
	stmt2 := `SELECT "tag" FROM "main"."tags" WHERE "id" = ?`
	var tags, tag = []string{}, ""
	var tagID int64
	rows, err := db.Query(stmt1, postID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&tagID)
		if err != nil {
			return nil, err
		}
		row := db.QueryRow(stmt2, tagID)
		err = row.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func getPostVotes(postID int64) (*Vote, error) {
	stmt := `SELECT vote FROM post_votes WHERE post_id = ?`
	rows, err := db.Query(stmt, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	vote, votes := 0, Vote{}
	for rows.Next() {
		err = rows.Scan(&vote)
		if err != nil {
			return nil, err
		}
		if vote == 1 {
			votes.Likes++
		} else {
			votes.Dislikes++
		}
	}
	return &votes, nil
}

func getPostComments(postID int64) ([]Comment, error) {
	stmt := `SELECT * FROM "main"."comments" WHERE "post_id" = ? ORDER BY "created_at" DESC`

	rows, err := db.Query(stmt, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cmt, cmts := Comment{}, []Comment{}
	for rows.Next() {
		err = rows.Scan(&cmt.ID, &cmt.PostID, &cmt.UserID, &cmt.Username, &cmt.Text, &cmt.CreatedAt)
		if err != nil {
			return nil, err
		}
		votes, err := getCommentVotes(cmt.ID)
		if err != nil {
			return nil, err
		}
		cmt.Votes = *votes
		cmts = append(cmts, cmt)
	}
	return cmts, nil
}

func getCommentVotes(cmtID int64) (*Vote, error) {
	stmt := `SELECT vote FROM comment_votes WHERE comment_id = ?`
	rows, err := db.Query(stmt, cmtID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	vote, votes := 0, Vote{}
	for rows.Next() {
		err = rows.Scan(&vote)
		if err != nil {
			return nil, err
		}
		if vote == 1 {
			votes.Likes++
		} else {
			votes.Dislikes++
		}
	}
	return &votes, nil
}

//	CHECK ALL SQL STATEMENTS FOR CORRECTNESS AS YOU CHANGED SOME TABLES

// CHECK DISMATCHES
//func (m *ForumModel) GetVotePosts(login string, vote int) (*[]models.Post, error) {
//func GetVotedPostsByUserID(userID int64, vote int) (*[]Post, error) {
//----------
