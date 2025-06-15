package consts

import (
	"fmt"
	"strings"
	"time"
)

type EVENT string

const (
	APP_NAME = "MaxxGuI"
	DB_NAME  = "maxx_gui.db"

	EVENT_PROGRESS = "EVENT_PROGRESS"
	EVENT_RESULT   = "EVENT_RESULT"
)

var (
	COPYWRIGHT = fmt.Sprintf("Copyright © 2025-%d dusbot.\nAll rights reserved.", time.Now().Year())

	EventEnums = []EVENT{EVENT_PROGRESS, EVENT_RESULT}
)

func (c EVENT) TSName() string {
	return strings.ToUpper(string(c))
}
