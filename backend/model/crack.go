package model

type CrackMetadata struct {
	Status    byte
	StartTime int64
	EndTime   int64

	Index, Total int
	Progress     float64
	LastCost     int
}

type CrackTask struct {
	CrackMetadata
	ID                           string `gorm:"type:varchar(32);index;not null"`
	Targets                      string `gorm:"not null"`
	Usernames                    string `gorm:"null"`
	Passwords                    string `gorm:"null"`
	Proxies                      string `gorm:"null"`
	Thread, Interval, MaxRuntime int    `gorm:"null"`
}

type CrackResult struct {
	ID       string `gorm:"index;not null"`
	Target   string `gorm:"index;not null"`
	Service  string `gorm:"index;null"`
	Username string `gorm:"null"`
	Password string `gorm:"null"`
}

type CrackProvider struct {
	CrackTask   CrackTask
	CrackResult CrackResult
}

func (p *CrackProvider) Register(task CrackTask, result CrackResult) {
	p.CrackTask = task
	p.CrackResult = result
}

var DefaultCrackProvider = &CrackProvider{}
