package control

import (
	"dmGameServer/model"
	pb "dmGameServer/pb"
)

// 根据号码查询主播  ir是用于获取的时候提示错误
func GetAnchorById(AccountId string, ir ...interface{}) (av *pb.AnchorDBInfo, err error) {
	if len(ir) > 0 {
		av, err = model.ModelGetAnchorById(AccountId, "1")
	} else {
		av, err = model.ModelGetAnchorById(AccountId)
	}
	return

}
