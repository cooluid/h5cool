package typedefine

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type AccountActor struct {
	ActorId   float64
	ActorName string
	Camp      int
	Sex       int
	Level     int
}

type Account struct {
	AccountId int
	Actor     *Actor
	GmLevel   byte

	conn   *websocket.Conn
	closed bool

	lock sync.RWMutex

	datas [][]byte
}

func NewAccount(conn *websocket.Conn) *Account {
	account := &Account{conn: conn, datas: make([][]byte, 0)}
	go func() {
		write := account.conn.WriteMessage
		bm := websocket.BinaryMessage
		loopTime := time.Millisecond * 10
		timeout := time.Millisecond
		for {
			select {
			case <-time.After(loopTime):
				account.lock.Lock()
				if account.closed {
					account.lock.Unlock()
					return
				}

				if len(account.datas) == 0 {
					account.lock.Unlock()
					break
				}

				tick := time.Now()
				var index int
				for _, data := range account.datas {
					if len(data) == 0 {
						account.closed = true
						account.conn.Close()
						account.lock.Unlock()
						return
					}
					if write(bm, data) != nil {
						break
					}
					index++
					if time.Since(tick) > timeout {
						break
					}
				}
				account.datas = account.datas[index:]

				account.lock.Unlock()
			}
		}
	}()

	return account
}

func (account *Account) Close() {
	account.lock.Lock()
	defer account.lock.Unlock()

	if account.closed {
		return
	}
	account.closed = true
	account.conn.Close()
}

func (account *Account) IsClose() bool {
	account.lock.RLock()
	defer account.lock.RUnlock()

	return account.closed
}

func (account *Account) Reply(data []byte) {
	if account.IsClose() {
		return
	}

	account.datas = append(account.datas, data)
}
