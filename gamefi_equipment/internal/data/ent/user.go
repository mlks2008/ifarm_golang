package ent

import (
	"github.com/go-pg/pg/v10/orm"
	"time"
)

var (
	TableUser = (*UserEnt)(nil)
)

type UserEnt struct {
	tableName   struct{}  `pg:"app.gateway_user"`
	RealUserId  string    `pg:"real_userid" json:"real_userid"`
	TokenUserId string    `pg:"token_userid" json:"token_userid"`
	Web3Address string    `pg:"web3_address,default:''" json:"web3_address"`
	HdAddress   string    `pg:"hd_address,default:''" json:"hd_address"`
	CreateTime  time.Time `pg:"create_time,default:now()"`
	Updatetime  time.Time `pg:"update_time"`
}

func (e *UserEnt) Name() string {
	model, _ := orm.NewModel(e)
	return string(model.(orm.TableModel).Table().SQLName)
}
