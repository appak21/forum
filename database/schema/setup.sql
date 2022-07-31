-- USERS TABLE --
CREATE TABLE IF NOT EXISTS "users" (
  "id" INTEGER NOT NULL UNIQUE, 
  "username" TEXT NOT NULL UNIQUE, 
  "email" TEXT NOT NULL UNIQUE, 
  "password" TEXT NOT NULL, 
  "created_at" TEXT NOT NULL, 
  PRIMARY KEY("id" AUTOINCREMENT)
);

-- POSTS TABLE --
CREATE TABLE IF NOT EXISTS "posts" (
  "id" INTEGER NOT NULL UNIQUE, 
  "user_id" INTEGER NOT NULL, 
  "username" TEXT NOT NULL, 
  "title" TEXT NOT NULL, 
  "text" TEXT NOT NULL, 
  "created_at" TEXT NOT NULL, 
  PRIMARY KEY("id" AUTOINCREMENT), 
  FOREIGN KEY("user_id") REFERENCES "users"("id")
);

-- POST VOTES TABLE --
CREATE TABLE IF NOT EXISTS "post_votes" (
  "id" INTEGER NOT NULL, 
  "user_id" INTEGER NOT NULL, 
  "post_id" INTEGER NOT NULL, 
  "vote" INTEGER NOT NULL, 
  FOREIGN KEY("post_id") REFERENCES "posts"("id"), 
  FOREIGN KEY("user_id") REFERENCES "users"("id"), 
  PRIMARY KEY("id" AUTOINCREMENT)
);

-- COMMENTS TABLE --
CREATE TABLE IF NOT EXISTS "comments" (
  "id" INTEGER NOT NULL UNIQUE, 
  "post_id" INTEGER NOT NULL, 
  "user_id" INTEGER NOT NULL, 
  "username" TEXT NOT NULL, 
  "text" TEXT NOT NULL, 
  "created_at" TEXT NOT NULL, 
  PRIMARY KEY("id" AUTOINCREMENT), 
  FOREIGN KEY("user_id") REFERENCES "users"("id"), 
  FOREIGN KEY("post_id") REFERENCES "posts"("id")
);

-- COMMENT VOTES TABLE --
CREATE TABLE IF NOT EXISTS "comment_votes" (
  "id" INTEGER NOT NULL, 
  "user_id" INTEGER NOT NULL, 
  "comment_id" INTEGER NOT NULL, 
  "vote" INTEGER NOT NULL, 
  FOREIGN KEY("user_id") REFERENCES "users"("id"), 
  FOREIGN KEY("comment_id") REFERENCES "comments"("id"), 
  PRIMARY KEY("id" AUTOINCREMENT)
);

-- TAGS TABLE --
CREATE TABLE IF NOT EXISTS "tags" (
  "id" INTEGER NOT NULL UNIQUE, 
  "tag" TEXT NOT NULL UNIQUE, 
  PRIMARY KEY("id" AUTOINCREMENT)
);

-- POSTS AND TAGS TABLE --
CREATE TABLE IF NOT EXISTS "posts_and_tags" (
  "id" INTEGER NOT NULL UNIQUE, 
  "post_id" INTEGER NOT NULL, 
  "tag_id" INTEGER NOT NULL, 
  PRIMARY KEY("id" AUTOINCREMENT), 
  FOREIGN KEY("post_id") REFERENCES "posts"("id"), 
  FOREIGN KEY("tag_id") REFERENCES "tags"("id")
);
