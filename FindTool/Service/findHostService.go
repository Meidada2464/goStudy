package Service

import (
	"fmt"
	"goStudy/FindTool/InitDir"
	"goStudy/FindTool/Model"
)

func FindHost() {
	db, err := InitDir.InitDataBase()
	if err != nil {
		return
	}
	//fmt.Println("db:", db)
	var st Model.Strategy
	db.Model(&Model.Strategy{}).Where("id = ?", 2276458).Find(&st)
	fmt.Println("st:", st)
	//var strategy Model.Strategy
	//db.Table("strategy").Where("id = ?", 1965899).Find(&strategy)
	//println("findStrategy", strategy)
}
