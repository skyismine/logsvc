package storage

import (
	"CommonUtil/src/GYGUtils"
	"encoding/json"
)

type StorageNanomsg struct {
	pubc GYGUtils.GSocket
}

func NewStorageStorageNanomsg(addr string) *StorageNanomsg {
	nanomsg := new(StorageNanomsg)
	nanomsg.pubc = GYGUtils.PubNode(addr, nil)
	return nanomsg
}

func (store *StorageNanomsg) Save(msg *Logmsg) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return GYGUtils.GSocketSend(store.pubc, data, 0)
}

func (store *StorageNanomsg) Close() {
	_ = store.pubc.Close()
}
