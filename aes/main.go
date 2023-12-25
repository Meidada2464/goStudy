package main

import (
	"fmt"
	"goStudy/aes/models"
	"goStudy/aes/utils"
)

func main() {

	data, err := utils.AESDecryptFromFile(models.Cake, "/Users/meifengfeng/workSpace/study/goStudy/goStudy/aes/file/.version.hot.yaml")
	if err != nil {
		fmt.Println("err -> ", err)
		return
	}

	fmt.Println(string(data))

}
