package actormgr

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"

	"github.com/sencydai/h5cool/data"
	"github.com/sencydai/h5cool/dispatch"
	"github.com/sencydai/h5cool/engine"
	g "github.com/sencydai/h5cool/gconfig"
	"github.com/sencydai/h5cool/log"
	"github.com/sencydai/h5cool/proto/pack"
	proto "github.com/sencydai/h5cool/proto/protocol"
	"github.com/sencydai/h5cool/service"
	"github.com/sencydai/h5cool/timer"
	t "github.com/sencydai/h5cool/typedefine"

	"github.com/json-iterator/go"
)

var (
	maxActorId int64
	actorNames map[string]int64

	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

const (
	loginKey = "QazXswEdc&141009522@"
)

func OnLoadMaxActorId() {
	var err error
	maxActorId, err = engine.GetMaxActorId()
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
	}

	if maxActorId == 0 {
		maxActorId = g.ServerIdML
	}
}

func newActorId() int64 {
	maxActorId++
	return maxActorId
}

func OnLoadAllActorNames() {
	var err error
	actorNames, err = engine.GetAllActorNames()
	if err != nil {
		panic(err)
	}

	for name := range actorNames {
		g.UseRandomName(name)
	}
}

func IsActorNameExist(name string) bool {
	_, ok := actorNames[name]
	return ok
}

func GetActorId(name string) (int64, bool) {
	actorId, ok := actorNames[name]
	return actorId, ok
}

func AppendActorName(name string, actorId int64) {
	actorNames[name] = actorId
}

func RemoveActorName(name string) {
	delete(actorNames, name)
}

func init() {
	dispatch.RegAccountMsgHandle(proto.AccountCLogin, onAccountLogin)

	service.RegGameStart(onGameStart)
	service.RegGm("stat", onGetActorCount)
}

type AccountLoginInfo struct {
	Code    int
	Id      int
	Level   byte
	Timeout int64
	Session []byte
}

//账号登录
func onAccountLogin(account *t.Account, reader *bytes.Reader) {
	retCode := -1
	var accountId int
	defer func() {
		if retCode != 0 {
			account.Reply(pack.EncodeWriter(pack.AllocPack(proto.AccountSLogin, retCode)))
			account.Reply(nil)
			return
		}
		account.Reply(pack.EncodeWriter(pack.AllocPack(proto.AccountSLogin, retCode, accountId)))
	}()

	if account.AccountId != 0 {
		return
	}
	var text string
	pack.Read(reader, &text)
	log.Info(text)

	loginInfo := &AccountLoginInfo{}
	err := json.UnmarshalFromString(text, loginInfo)
	if err != nil {
		log.Error(err)
		return
	}
	if loginInfo.Code != 0 {
		log.Error("login account,code != 0")
		return
	}
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d%d%d%d%s",
		loginInfo.Code, loginInfo.Id, loginInfo.Level, loginInfo.Timeout, loginKey)))
	session := hash.Sum(nil)
	if bytes.Compare(session, loginInfo.Session) != 0 {
		log.Error("login account,no match session")
		return
	}
	if time.Now().Unix() > loginInfo.Timeout {
		log.Error("login account,time out")
		return
	}

	if data.GetAccount(loginInfo.Id) != nil {
		log.Errorf("login account, %d is login", loginInfo.Id)
		return
	}
	account.AccountId = loginInfo.Id
	account.GmLevel = loginInfo.Level
	data.AppendAccount(account)
	retCode = 0
	accountId = account.AccountId
}

func onActorLogin(account *t.Account, reader *bytes.Reader) {
}

func actorLogout(actor *t.Actor) {
	timer.StopActorTimers(actor)
	service.OnActorLogout(actor)
	data.RemoveOnlineActor(actor)
}

func OnAccountLogout(account *t.Account) {
	if account.Actor != nil {
		actorLogout(account.Actor)
		account.Actor = nil
	}
	if account.AccountId != 0 {
		data.RemoveAccount(account.AccountId)
	}
}

func onGameStart() {
	timer.Loop(nil, "flushactors", time.Minute*5, time.Minute*5, -1, func() {
		go func() {
			engine.FlushActorBuffers()
		}()
	})
}

func OnGameClose() {
	tick := time.Now()
	data.LoopActors(func(actor *t.Actor) bool {
		service.OnActorLogout(actor)
		actor.LogoutTime = time.Now()
		engine.UpdateActor(actor)
		return true
	})
	engine.FlushActorBuffers()
	log.Infof("save actors data: %v", time.Since(tick))
}

func onGetActorCount(map[string]string) (int, string) {
	return 0, fmt.Sprintf(
		"max: %d, real: %d, account: %d, online: %d, offline: %d, engineBuff: %d",
		g.GetMaxCount(),
		g.GetRealCount(),
		data.GetAccountCount(),
		data.GetOnlineCount(),
		data.GetCacheCount(),
		engine.GetCacheCount(),
	)
}
