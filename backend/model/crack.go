package model

type CrackMetadata struct {
	Status    byte
	StartTime int
	EndTime   int

	Progress int
	LastCost int
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
