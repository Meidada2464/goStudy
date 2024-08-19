### 白山mallard监控go语言版本使用规范库  
白山go-bspub库统一使用方法，详见url：http://jr.baishancloud.com:8090/display/FLOW/golang

```
使用方法：
    //获取mallard客户端实例  
        mClient, err := bsmallard.MustNewClient(&app.ZConfig.Mallard) 
    //进行信息推送
        // 构造推送数据
    	var mallardSlice []*bsmallard.MetricRaw
     	fieldsMap := make(map[string]interface{})
        fieldsMap["test_field"] = 1
        tagsMap := make(map[string]string)
        tagsMap["test_tag"] = "1"
        mallardSlice = append(mallardSlice, &bsmallard.MetricRaw{
            Name:     "test_name", //需找监控平台申请表名，申请地址https://ares.bs58i.baishancloud.com/
            Time:     time.Now().Unix(),
            Value:    1,
            Fields:   fieldsMap,
            Tags:     tagsMap,
            Endpoint: "test_endpoint",
        })
        // 推送上报数据
        if err := mClient.Report(mallardSlice); err != nil {
            fmt.Println(err)
        }
    
```
#### 相关上报指导文档
1，本机接口上报，包括报警跟踪配置文档：  
    http://jr.baishancloud.com:8090/pages/viewpage.action?pageId=163557264  

2，远程上报，针对无本机agent服务，需要进行远程上报：  
    http://jr.baishancloud.com:8090/pages/viewpage.action?pageId=29295737  

3，本机写文件上报，不推荐，当前sdk也还不支持  

