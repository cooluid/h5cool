package service

import (
	"runtime/debug"
	"time"

	"github.com/sencydai/h5cool/log"
	t "github.com/sencydai/h5cool/typedefine"
)

type gameStartHandle func()
type gameCloseHandle func()
type configLoadFinishHandle func(bool)
type systemNewDayHandle func()
type systemTimeChangeHandle func()

type gmHandle func(map[string]string) (int, string)

type actorCreateHandle func(*t.Actor)
type actorDataLoadHandle func(*t.Actor)
type actorBeforeLoginHandle func(*t.Actor, int)
type actorLoginHandle func(*t.Actor, int)
type actorLogoutHandle func(*t.Actor)
type actorNewDayHandle func(*t.Actor)
type actorMinTimerHandle func(*t.Actor, int)
type actorUpgradeHandle func(*t.Actor, int)
type actorTimeChangeHandle func(*t.Actor)

var (
	gameStartHandles        = make([]gameStartHandle, 0)
	gameCloseHandles        = make([]gameCloseHandle, 0)
	configLoadHandles       = make([]configLoadFinishHandle, 0)
	systemNewDayHandles     = make([]systemNewDayHandle, 0)
	systemTimeChangeHandles = make([]systemTimeChangeHandle, 0)

	gmHandles = make(map[string]gmHandle)

	actorCreateHandles      = make([]actorCreateHandle, 0)
	actorDataLoadHandles    = make([]actorDataLoadHandle, 0)
	actorBeforeLoginHandles = make([]actorBeforeLoginHandle, 0)
	actorLoginHandles       = make([]actorLoginHandle, 0)
	actorLogoutHandles      = make([]actorLogoutHandle, 0)
	actorNewDayHandles      = make([]actorNewDayHandle, 0)
	actorMinTimerHandles    = make([]actorMinTimerHandle, 0)
	actorUpgradeHandles     = make([]actorUpgradeHandle, 0)
)

func RegGameStart(handle func()) {
	gameStartHandles = append(gameStartHandles, handle)
}

func RegGameClose(handle func()) {
	gameCloseHandles = append(gameCloseHandles, handle)
}

func RegConfigLoadFinish(handle func(isGameStart bool)) {
	configLoadHandles = append(configLoadHandles, handle)
}

func RegSystemNewDay(handle func()) {
	systemNewDayHandles = append(systemNewDayHandles, handle)
}

func RegSystemTimeChange(handle func()) {
	systemTimeChangeHandles = append(systemTimeChangeHandles, handle)
}

func RegGm(cmd string, handle func(values map[string]string) (int, string)) {
	gmHandles[cmd] = handle
}

func GetGmHandle(cmd string) gmHandle {
	return gmHandles[cmd]
}

func RegActorCreate(handle func(*t.Actor)) {
	actorCreateHandles = append(actorCreateHandles, handle)
}

func RegActorDataLoad(handle func(*t.Actor)) {
	actorDataLoadHandles = append(actorDataLoadHandles, handle)
}

func RegActorBeforeLogin(handle func(actor *t.Actor, offSec int)) {
	actorBeforeLoginHandles = append(actorBeforeLoginHandles, handle)
}

func RegActorLogin(handle func(actor *t.Actor, offSec int)) {
	actorLoginHandles = append(actorLoginHandles, handle)
}

func RegActorLogout(handle func(*t.Actor)) {
	actorLogoutHandles = append(actorLogoutHandles, handle)
}

func RegActorNewDay(handle func(*t.Actor)) {
	actorNewDayHandles = append(actorNewDayHandles, handle)
}

func RegActorMinTimer(handle func(actor *t.Actor, times int)) {
	actorMinTimerHandles = append(actorMinTimerHandles, handle)
}

func RegActorUpgrade(handle func(actor *t.Actor, oldLevel int)) {
	actorUpgradeHandles = append(actorUpgradeHandles, handle)
}

func OnGameStart() {
	for _, handle := range gameStartHandles {
		handle()
	}
}

func OnGameClose() {
	for _, handle := range gameCloseHandles {
		handle()
	}
}

func OnConfigReloadFinish(isGameStart bool) {
	for _, handle := range configLoadHandles {
		handle(isGameStart)
	}
}

func OnSystemNewDay() {
	for _, handle := range systemNewDayHandles {
		handle()
	}
}

func OnSystemTimeChange() {
	for _, handle := range systemTimeChangeHandles {
		handle()
	}
}

func OnActorCreate(actor *t.Actor) {
	for _, handle := range actorCreateHandles {
		handle(actor)
	}
}

func OnActorDataLoad(actor *t.Actor) {
	for _, handle := range actorDataLoadHandles {
		handle(actor)
	}
}

func OnActorBeforeLogin(actor *t.Actor, offSec int) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("actor(%d) before login error: %v,%s", actor.ActorId, err, string(debug.Stack()))
		}
	}()
	for _, handle := range actorBeforeLoginHandles {
		handle(actor, offSec)
	}
}

func OnActorLogin(actor *t.Actor, offSec int) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("actor(%d) login error: %v,%s", actor.ActorId, err, string(debug.Stack()))
		}
	}()
	for _, handle := range actorLoginHandles {
		handle(actor, offSec)
	}
}

func OnActorLogout(actor *t.Actor) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("actor(%d) logout error: %v,%s", actor.ActorId, err, string(debug.Stack()))
		}
	}()
	for _, handle := range actorLogoutHandles {
		handle(actor)
	}
}

func OnActorNewDay(actor *t.Actor) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("actor(%d) newday error: %v,%s", actor.ActorId, err, string(debug.Stack()))
		}
		exData := actor.GetExData()
		exData.NewDay = time.Now().Unix()
	}()
	for _, handle := range actorNewDayHandles {
		handle(actor)
	}
}

func OnActorMinTimer(actor *t.Actor, times int) {
	for _, handle := range actorMinTimerHandles {
		handle(actor, times)
	}
}

func OnActorUpgrade(actor *t.Actor, oldLevel int) {
	for _, handle := range actorUpgradeHandles {
		handle(actor, oldLevel)
	}
}
