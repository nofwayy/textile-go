package db

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/textileio/textile-go/repo"
)

type CafeRequestDB struct {
	modelStore
}

func NewCafeRequestStore(db *sql.DB, lock *sync.Mutex) repo.CafeRequestStore {
	return &CafeRequestDB{modelStore{db, lock}}
}

func (c *CafeRequestDB) Add(req *repo.CafeRequest) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stm := `insert into cafe_requests(id, peerId, targetId, cafeId, cafe, type, date) values(?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(stm)
	if err != nil {
		log.Errorf("error in tx prepare: %s", err)
		return err
	}
	defer stmt.Close()

	cafe, err := json.Marshal(req.Cafe)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		req.Id,
		req.PeerId,
		req.TargetId,
		req.Cafe.Peer,
		cafe,
		req.Type,
		req.Date.UnixNano(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (c *CafeRequestDB) List(offset string, limit int) []repo.CafeRequest {
	c.lock.Lock()
	defer c.lock.Unlock()
	var stm string
	if offset != "" {
		stm = "select * from cafe_requests where date>(select date from cafe_requests where id='" + offset + "') order by date asc limit " + strconv.Itoa(limit) + ";"
	} else {
		stm = "select * from cafe_requests order by date asc limit " + strconv.Itoa(limit) + ";"
	}
	return c.handleQuery(stm)
}

func (c *CafeRequestDB) Delete(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_requests where id=?", id)
	return err
}

func (c *CafeRequestDB) DeleteByCafe(cafeId string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.db.Exec("delete from cafe_requests where cafeId=?", cafeId)
	return err
}

func (c *CafeRequestDB) handleQuery(stm string) []repo.CafeRequest {
	var ret []repo.CafeRequest
	rows, err := c.db.Query(stm)
	if err != nil {
		log.Errorf("error in db query: %s", err)
		return nil
	}
	for rows.Next() {
		var id, peerId, targetId, cafeId string
		var typeInt int
		var dateInt int64
		var cafe []byte
		if err := rows.Scan(&id, &peerId, &targetId, &cafeId, &cafe, &typeInt, &dateInt); err != nil {
			log.Errorf("error in db scan: %s", err)
			continue
		}

		var mod repo.Cafe
		if err := json.Unmarshal(cafe, &mod); err != nil {
			log.Errorf("error unmarshaling cafe: %s", err)
			continue
		}

		ret = append(ret, repo.CafeRequest{
			Id:       id,
			PeerId:   peerId,
			TargetId: targetId,
			Cafe:     mod,
			Type:     repo.CafeRequestType(typeInt),
			Date:     time.Unix(0, dateInt),
		})
	}
	return ret
}
