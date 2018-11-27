package data

import "strconv"

type ThreadClassInfo struct {
		Id			int64
		Name 		string
		ThreadCount string
}

func GetThreadClassInfo() (tci []ThreadClassInfo) {
	rows, err := Db.Query(
		"select tc.id, tc.name, t.count " +
			"from thread_class as tc, (select class_id, count(*) as count from threads group by class_id) as t " +
			"where tc.id=t.class_id")
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		conv := ThreadClassInfo{}
		if err = rows.Scan(&conv.Id, &conv.Name, &conv.ThreadCount); err != nil {
			return
		}
		tci = append(tci, conv)
	}
	return
}

func TotalThreadsOfClass(classId int64) (total int64, err error) {
	err = Db.QueryRow("select count(id) from threads where class_id = ?", classId).Scan(&total)
	return
}

func ThreadsByClassAndPage(classId, CurrentPage, PageSize int64) (threads []Thread, err error) {
	rows, err := Db.Query(
		"SELECT id, uuid, topic, body, user_id, created_at FROM threads " +
			"WHERE class_id = ? " +
			"ORDER BY created_at DESC " +
			"LIMIT ?, ?", classId, (CurrentPage-1)*PageSize, PageSize)
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

func AddThreadClass(name string) (err error) {
	statement := "insert into thread_class(name) values(?)"
	stmt, _ := Db.Prepare(statement)
	defer stmt.Close()

	_, err = stmt.Exec(name)
	return
}

func DeleteThreadClassById(id string) (err error) {
	statement := "delete from thread_class where id = ?"
	stmt, _ := Db.Prepare(statement)
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return
}

func ModifyThreadClass(uuid, classId string) (err error) {
	cId, err := strconv.ParseInt(classId, 10, 64)
	if err != nil {
		return
	}

	statement := "update threads set class_id = ? where uuid = ?"
	stmt, _ := Db.Prepare(statement)
	defer stmt.Close()

	_, err = stmt.Exec(cId, uuid)
	return
}

func DeleteThread(uuid string) (err error) {
	statement := "delete from threads where uuid = ?"
	stmt, _ := Db.Prepare(statement)
	defer stmt.Close()

	_, err = stmt.Exec(uuid)
	return
}