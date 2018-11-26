package data

import (
	"time"
	)

type Thread struct {
	Id        int
	Uuid      string
	Topic     string
	Body      string
	UserId    int
	ClassName string
	CreatedAt time.Time
}

type Post struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int
	ThreadId  int
	CreatedAt time.Time
}

// format the CreatedAt date to display nicely on the screen
func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt.Format("2006-01-02 15:04:05")
}

func (post *Post) CreatedAtDate() string {
	return post.CreatedAt.Format("2006-01-02 15:04:05")
}

// get the number of posts in a thread
func (thread *Thread) NumReplies() (count int) {
	rows, err := Db.Query("SELECT count(*) FROM posts where thread_id = ?", thread.Id)
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return
		}
	}
	return
}

// get posts to a thread
func (thread *Thread) Posts() (posts []Post, err error) {
	rows, err := Db.Query("SELECT id, uuid, body, user_id, thread_id, created_at FROM posts where thread_id = ?", thread.Id)
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt); err != nil {
			return
		}
		posts = append(posts, post)
	}
	return
}

// Create a new thread
func (user *User) CreateThread(topic, body string, classId int64) (conv Thread, err error) {
	statement := "insert into threads (uuid, topic, body, user_id, class_id, created_at) values (?, ?, ?, ?, ?, ?)"
	stmt, _ := Db.Prepare(statement)
	defer stmt.Close()

	uuid := createUUID()
	// use QueryRow to return a row and scan the returned id into the Session struct
	_, err = stmt.Exec(uuid, topic, body, user.Id, classId, time.Now())
	if err != nil {
		return
	}
	conv, err = ThreadByUUID(uuid)
	return
}

// Create a new post to a thread
func (user *User) CreatePost(conv Thread, body string) (post Post, err error) {
	statement := "insert into posts (uuid, body, user_id, thread_id, created_at) values (?, ?, ?, ?, ?)"
	stmt, _ := Db.Prepare(statement)
	defer stmt.Close()

	uuid := createUUID()
	// use QueryRow to return a row and scan the returned id into the Session struct
	_, err = stmt.Exec(uuid, body, user.Id, conv.Id, time.Now())
	if err != nil {
		return
	}
	post, err = PostByUUID(uuid)
	return
}

// Get all threads in the database and returns it
func Threads() (threads []Thread, err error) {
	rows, err := Db.Query(
		"SELECT t.id, t.uuid, t.topic, t.body, t.user_id, c.name, t.created_at " +
			"FROM threads t, thread_class c " +
			"WHERE t.class_id = c.id " +
			"ORDER BY created_at DESC")
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		conv := Thread{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.Body, &conv.UserId, &conv.CreatedAt); err != nil {
			return
		}
		threads = append(threads, conv)
	}
	return
}

// Get a thread by the UUID
func ThreadByUUID(uuid string) (conv Thread, err error) {
	conv = Thread{}
	err = Db.QueryRow(
		"SELECT t.id, t.uuid, t.topic, t.body, t.user_id, c.name, t.created_at " +
			"FROM threads t, thread_class c " +
			"WHERE t.uuid = ? and t.class_id = c.id", uuid).
		Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.Body, &conv.UserId, &conv.ClassName, &conv.CreatedAt)
	return
}

func ThreadsByPage(CurrentPage, PageSize int64) (threads []Thread, err error) {
	rows, err := Db.Query(
		"SELECT id, uuid, topic, body, user_id, created_at FROM threads ORDER BY created_at DESC " +
			   "LIMIT ?, ?", (CurrentPage-1)*PageSize, PageSize)
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		conv := Thread{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.Body, &conv.UserId, &conv.CreatedAt); err != nil {
			return
		}
		threads = append(threads, conv)
	}
	return
}

func TotalThreads() (total int64, err error) {
	err = Db.QueryRow("select count(id) from threads").Scan(&total)
	return
}

// Get a post by the UUID
func PostByUUID(uuid string) (conv Post, err error) {
	conv = Post{}
	err = Db.QueryRow("select id, uuid, body, user_id, thread_id, created_at from posts where uuid = ?", uuid).
		Scan(&conv.Id, &conv.Uuid, &conv.Body, &conv.UserId, &conv.ThreadId, &conv.CreatedAt)
	return
}

// Get the user who started this thread
func (thread *Thread) User() (user User) {
	user = User{}
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = ?", thread.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

// Get the user who wrote the post
func (post *Post) User() (user User) {
	user = User{}
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = ?", post.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

func (thread *Thread) CompatibleNonCkeidtor() (b bool) {
	if thread.Id >= 1 && thread.Id <= 8 {
		return true
	}

	return false
}

func (post *Post) CompatibleNonCkeidtor() (b bool) {
	if post.ThreadId >= 1 && post.ThreadId <= 8 {
		return true
	}

	return false
}

type ThreadArchive struct {
	Month	string
	Count 	int64
}

func ThreadArchived() (t []ThreadArchive) {
	rows, err := Db.Query(
		"select DATE_FORMAT(created_at,'%Y年%m月') months, count(id) count from threads group by months")
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		conv := ThreadArchive{}
		if err = rows.Scan(&conv.Month, &conv.Count); err != nil {
			return
		}
		t = append(t, conv)
	}
	return
}
