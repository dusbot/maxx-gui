package handler

import (
	"context"
	"math"
	"maxxgui/backend/consts"
	"maxxgui/backend/model"
	"maxxgui/backend/query"
	"maxxgui/backend/report"
	"maxxgui/backend/utils"

	"strings"
	"sync"
	"time"

	"github.com/dusbot/maxx/core/types"
	"github.com/dusbot/maxx/libs/slog"
	utils_ "github.com/dusbot/maxx/libs/utils"
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
	if task.MaxRuntime > 0 {
		if task.MaxRuntime < 30 {
			task.MaxRuntime = 30
		}
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(task.MaxRuntime)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	CancelMap.Store(uuid, cancel)
	targets := utils_.ParseNetworkInput(task.Targets)
	innerTask := &types.Task{
		MaxTime:      task.MaxRuntime,
		Interval:     task.Interval,
		Progress:     true,
		Thread:       task.Thread,
		Targets:      targets,
		Users:        strings.Split(task.Usernames, ","),
		Passwords:    strings.Split(task.Passwords, ","),
		ResultChan:   make(chan types.Result, 1<<8),
		ProgressChan: make(chan types.Progress, 1<<8),
	}

	saveTask2DB(c.Query, &task)

	go handleProgress(c.Ctx, uuid, c.Query, innerTask.ProgressChan)

	go handleResult(c.Ctx, uuid, c.Query, innerTask.ResultChan)

	go run.Crack(ctx, innerTask)

	return
}

func (c *CrackHandler) Cancel(id string) (ok bool) {
	var cancel any
	cancel, ok = CancelMap.Load(id)
	if ok {
		defer CancelMap.Delete(id)
		cancel.(context.CancelFunc)()
		slog.Printf(slog.DEBUG, "CrackTask[%s] cancelled", id)
	} else {
		slog.Printf(slog.WARN, "CrackTask[%s] not found", id)
	}
	return
}

func (c *CrackHandler) CancelAll() (ok bool) {
	defer CancelMap.Clear()
	CancelMap.Range(func(key, value any) bool {
		value.(context.CancelFunc)()
		ok = true
		return ok
	})
	if ok {
		slog.Printf(slog.DEBUG, "All CrackTask cancelled")
	} else {
		slog.Printf(slog.WARN, "No CrackTask found")
	}
	return
}

func (c *CrackHandler) GenerateReport(id string, zh bool) (content string) {
	results, err := c.Query.CrackResult.Where(c.Query.CrackResult.TaskID.Eq(id)).Find()
	if err != nil {
		slog.Printf(slog.WARN, "Query database with error[%+v]", err)
		return
	}
	task, err := c.Query.CrackTask.Where(c.Query.CrackTask.ID.Eq(id)).First()
	if err != nil {
		slog.Printf(slog.WARN, "Query CrackTask error[%+v]", err)
		return
	}
	_, content = report.DoGenCrackReport(zh, task, results)
	return
}

func saveTask2DB(q *query.Query, task *model.CrackTask) {
	q.CrackTask.Save(task)
}

func handleProgress(ctx context.Context, id string, q *query.Query, pipe chan types.Progress) {
	start := time.Now()
	var theLastProgress types.Progress
	defer runtime.EventsEmit(ctx, consts.EVENT_PROGRESS, 1)
	for progress := range pipe {
		theLastProgress = progress
		currProgressRate := math.Round(progress.Progress*100) / 100
		if progress.Progress >= 1 {
			currProgressRate = 0.99
		}
		runtime.EventsEmit(ctx, consts.EVENT_PROGRESS, currProgressRate)
	}
	q.CrackTask.
		Where(q.CrackTask.ID.Eq(id)).
		Updates(map[string]interface{}{
			"start_time": start.Unix(),
			"end_time":   time.Now().Unix(),
			"progress":   theLastProgress.Progress,
			"index":      int(theLastProgress.Done),
			"total":      int(theLastProgress.Total),
			"last_cost":  int(time.Since(start).Seconds()),
		})
	slog.Printf(slog.DEBUG, "progress pipe closed")
}

func handleResult(ctx context.Context, id string, q *query.Query, pipe chan types.Result) {
	for result := range pipe {
		slog.Printf(slog.DEBUG, "result: %+v", result)
		result_ := &model.CrackResult{
			TaskID:   id,
			Target:   result.Target,
			Service:  result.Protocol,
			Username: result.User,
			Password: result.Pass,
		}
		go runtime.EventsEmit(ctx, consts.EVENT_RESULT, result_)
		q.CrackResult.Save(result_)
	}
	slog.Printf(slog.DEBUG, "result pipe closed")
}
