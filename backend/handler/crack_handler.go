package handler

import (
	"context"
	"maxxgui/backend/consts"
	"maxxgui/backend/model"
	"maxxgui/backend/query"
	"maxxgui/backend/utils"
	"sync"

	"github.com/dusbot/maxx/core/types"
	"github.com/dusbot/maxx/libs/slog"
	"github.com/dusbot/maxx/run"
	"github.com/google/wire"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var CrackHandlerSet = wire.NewSet(wire.Struct(new(CrackHandler), "*"))

// id:cancelFunc
var CancelMap sync.Map

type CrackHandler struct {
	Ctx   context.Context
	Query *query.Query
}

func (c *CrackHandler) Scan(task model.CrackTask) (ok bool) {
	if task.Targets == "" {
		return
	}
	ok = true
	uuid, err := utils.GenerateTimestampUUID(8)
	if err != nil {
		slog.Printf(slog.WARN, "Failed to GenerateTimestampUUID with err[%+v]", err)
		return
	}
	task.ID = uuid
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if task.MaxRuntime != 0 {
		ctx, cancel = context.WithCancel(context.Background())
		CancelMap.Store(uuid, cancel)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	innerTask := &types.Task{
		MaxTime:      task.MaxRuntime,
		Interval:     task.Interval,
		Progress:     true,
		Thread:       task.Thread,
		Targets:      []string{task.Targets},
		Users:        []string{task.Usernames},
		Passwords:    []string{task.Passwords},
		ResultChan:   make(chan types.Result),
		ProgressChan: make(chan int),
	}

	innerTask.ProgressChan = make(chan int, 1<<8)
	innerTask.ResultChan = make(chan types.Result, 1<<8)

	go handleProgress(c.Ctx, uuid, c.Query, innerTask.ProgressChan)

	go handleResult(c.Ctx, uuid, c.Query, innerTask.ResultChan)

	go run.Crack(ctx, innerTask)

	return
}

func handleProgress(ctx context.Context, id string, q *query.Query, pipe chan int) {
	defer runtime.EventsEmit(ctx, consts.EVENT_PROGRESS, 1)
	for progress := range pipe {
		runtime.EventsEmit(ctx, consts.EVENT_PROGRESS, progress)
	}
	slog.Printf(slog.WARN, "progress pipe closed")
}

func handleResult(ctx context.Context, id string, q *query.Query, pipe chan types.Result) {
	for result := range pipe {
		result_ := &model.CrackResult{
			ID:       id,
			Target:   result.Target,
			Service:  result.Protocol,
			Username: result.User,
			Password: result.Pass,
		}
		runtime.EventsEmit(ctx, consts.EVENT_RESULT, result_)
		q.Transaction(func(tx *query.Query) error {
			tx.CrackResult.Save()
			return nil
		})
	}
}
