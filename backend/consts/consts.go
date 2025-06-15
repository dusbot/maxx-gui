package consts

import (
	"fmt"
	"time"
)

const (
	APP_NAME = "MaxxGuI"
	DB_NAME  = "maxx_gui.db"

	EVENT_PROGRESS = "EVENT_PROGRESS"
	EVENT_RESULT   = "EVENT_RESULT"
)

var COPYWRIGHT = fmt.Sprintf("Copyright © 2025-%d dusbot.\nAll rights reserved.", time.Now().Year())
