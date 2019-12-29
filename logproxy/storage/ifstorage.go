package storage

type IFStorage interface {
	Save(msg *Logmsg) error
	Close()
}

type Logmsg struct {
	App 	string	`json:"app"`
	Level 	string 	`json:"level"`
	Tag 	string 	`json:"tag"`
	Msg 	string 	`json:"msg"`
	Ctime 	string 	`json:"ctime"`
	Stime   string	`json:"stime"`
}