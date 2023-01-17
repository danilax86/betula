// Package db encapsulates all used queries to the database.
//
// Do not forget to Initialize and Finalize.
//
// All functions in this package might crash vividly.
package db

import (
	"git.sr.ht/~bouncepaw/betula/types"
	"log"
	"time"
)

const schema = `
create table if not exists Posts (
    ID integer primary key autoincrement not null,
    URL text not null,
    Title text not null,
    Description text not null,
    Visibility integer not null,
    CreationTime integer not null                   
);

create table if not exists Categories (
    ID integer primary key autoincrement not null,
    Name text not null
);

create table if not exists CategoriesToPosts (
    CategoryID integer not null,
    PostID integer not null
);

create table if not exists BetulaMeta (
    Key text primary key,
    Value text
);

insert or ignore into BetulaMeta values
	('DB version', 0),
	('Admin username', null),
	('Admin password hash', null);

create table if not exists Sessions (
    Token text primary key,
    CreationTime integer not null
);`

const sqlAddSession = `
insert into Sessions values (?, ?);
`

func AddSession(token string) {
	mustExec(sqlAddSession, token, time.Now().Unix())
}

const sqlHasSession = `
select exists(select 1 from Sessions where Token = ?);
`

func HasSession(token string) (has bool) {
	rows := mustQuery(sqlHasSession, token)
	rows.Next()
	mustScan(rows, &has)
	_ = rows.Close()
	return has
}

const sqlStopSession = `
delete from Sessions where Token = ?;
`

func StopSession(token string) {
	mustExec(sqlStopSession, token)
}

const sqlSetCredentials = `
insert or replace into BetulaMeta values
	('Admin username', ?),
	('Admin password hash', ?);
`

func SetCredentials(name, hash string) {
	mustExec(sqlSetCredentials, name, hash)
}

const sqlGetMetaEntry = `
select Value from BetulaMeta where Key = ? limit 1;
`

func MetaEntry[T any](key string) T {
	return querySingleValue[T](sqlGetMetaEntry, key)
}

const sqlPostsForCategory = `
select
	ID, URL, Title, Description, Visibility, CreationTime
from
	Posts
inner join
	CategoriesToPosts
where
	ID = PostID and CatID = ?;
`

const sqlCatNameByID = `
select Name from Categories where ID = ? limit 1;
`

func AuthorizedPostsForCategoryAndNameByID(authorized bool, id int) (name string, out chan types.Post) {
	rows := mustQuery(sqlPostsForCategory, id)
	out = make(chan types.Post)

	go func() {
		for rows.Next() {
			var post types.Post
			mustScan(rows, &post.ID, &post.URL, &post.Title, &post.Description, &post.Visibility, &post.CreationTime)
			if !authorized && post.Visibility == types.Private {
				continue
			}
			// TODO: Probably can be optimized with a smart query.
			post.Categories = CategoriesForPost(post.ID)
			out <- post
		}
		close(out)
	}()
	return querySingleValue[string](sqlCatNameByID, id), out
}

const sqlCategoriesForPost = `
select
    CatID, Name
from 
    CategoriesToPosts
inner join
    Categories
where
    ID = CatID and PostID = ?;
`

func CategoriesForPost(id int) (cats []types.Category) {
	rows := mustQuery(sqlCategoriesForPost, id)
	for rows.Next() {
		var cat types.Category
		mustScan(rows, &cat.ID, &cat.Name)
		cats = append(cats, cat)
	}
	return cats
}

const sqlGetAllCategories = `
select ID, Name from Categories;
`

func Categories() (cats []types.Category) {
	rows := mustQuery(sqlGetAllCategories)
	for rows.Next() {
		var cat types.Category
		mustScan(rows, &cat.ID, &cat.Name)
		cats = append(cats, cat)
	}
	return cats
}

const sqlGetAllPosts = `
select ID, URL, Title, Description, Visibility, CreationTime from Posts;
`

// YieldAuthorizedPosts returns a channel, from which you can get all posts stored in the database, along with their tags, but only if the viewer is authorized! Otherwise, only public posts will be given.
func YieldAuthorizedPosts(authorized bool) chan types.Post {
	rows := mustQuery(sqlGetAllPosts)
	out := make(chan types.Post)

	go func() {
		for rows.Next() {
			var post types.Post
			mustScan(rows, &post.ID, &post.URL, &post.Title, &post.Description, &post.Visibility, &post.CreationTime)
			if !authorized && post.Visibility == types.Private {
				continue
			}
			// TODO: Probably can be optimized with a smart query.
			post.Categories = CategoriesForPost(post.ID)
			out <- post
		}
		close(out)
	}()
	return out
}

const sqlAddPost = `
insert into Posts (URL, Title, Description, Visibility, CreationTime) VALUES (?, ?, ?, ?, ?);
`

// AddPost adds a new post to the database. Creation time is set by this function, ID is set by the database. The ID is returned.
func AddPost(post types.Post) int64 {
	post.CreationTime = time.Now().Unix()
	res, err := db.Exec(sqlAddPost, post.URL, post.Title, post.Description, post.Visibility, post.CreationTime)
	if err != nil {
		log.Fatalln(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatalln(err)
	}
	return id
}

const sqlEditPost = `
update Posts
set
    URL = ?,
    Title = ?,
    Description = ?,
    Visibility = ?,
    CreationTime = ?
where
    ID = ?;
`

func EditPost(post types.Post) {
	mustExec(sqlEditPost, post.URL, post.Title, post.Description, post.Visibility, post.CreationTime, post.ID)
}

const sqlPostForID = `
select ID, URL, Title, Description, Visibility, CreationTime from Posts where ID = ?;
`

// PostForID returns the post corresponding to the given id, if there is any.
func PostForID(id int) (post types.Post, found bool) {
	rows := mustQuery(sqlPostForID, id)
	rows.Next()
	mustScan(rows, &post.ID, &post.URL, &post.Title, &post.Description, &post.Visibility, &post.CreationTime)
	_ = rows.Close()
	return post, true
}

const sqlURLForID = `
select URL from Posts where ID = ?;
`

// URLForID returns the URL of the post corresponding to the given ID, if there is any post like that.
func URLForID(id int) (url string, found bool) {
	rows := mustQuery(sqlURLForID, id)
	rows.Next()
	mustScan(rows, &url)
	_ = rows.Close()
	return url, true
}

const sqlLinkCount = `select count(ID) from Posts;`
const sqlOldestTime = `select min(CreationTime) from Posts;`
const sqlNewestTime = `select max(CreationTime) from Posts;`

func LinkCount() int        { return querySingleValue[int](sqlLinkCount) }
func OldestTime() time.Time { return time.Unix(querySingleValue[int64](sqlOldestTime), 0) }
func NewestTime() time.Time { return time.Unix(querySingleValue[int64](sqlNewestTime), 0) }
