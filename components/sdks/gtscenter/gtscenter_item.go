package gtscenter

import (
	"components/sdks/global"
	"context"
	"fmt"
)

func (c *GtsCenter) AddItem(ctx context.Context, realUserId string, itemId string, itemSeq []string, itemNum []int64, recordId int64) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if itemId == global.Item_Id_Hero {
		err = c.platformDustCli.Hero_Add(ctx, realUserId, itemSeq, itemNum, recordId)
	} else {
		err = fmt.Errorf("item_id:%v is not exist", itemId)
	}

	return
}

func (c *GtsCenter) FreezeItem(ctx context.Context, realUserId string, itemId string, userItemSeq []string, recordId int64) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if itemId == global.Item_Id_Hero {
		err = c.platformDustCli.Hero_Freeze(ctx, realUserId, userItemSeq, recordId)
	} else {
		err = fmt.Errorf("%v is not exist", itemId)
	}

	return
}

func (c *GtsCenter) SubItem(ctx context.Context, realUserId string, itemId string, itemSeq []string, itemNum []int64, recordId int64) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if itemId == global.Item_Id_Hero {
		err = c.platformDustCli.Hero_Sub(ctx, realUserId, recordId)
	} else {
		err = fmt.Errorf("item_id:%v is not exist", itemId)
	}

	return
}

func (c *GtsCenter) ReturnItem(ctx context.Context, realUserId string, itemId string, recordId int64) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if itemId == global.Item_Id_Hero {
		err = c.platformDustCli.Hero_Return(ctx, realUserId, recordId)
	} else {
		err = fmt.Errorf("item_id:%v is not exist", itemId)
	}

	return
}
