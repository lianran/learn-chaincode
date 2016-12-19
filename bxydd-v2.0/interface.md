###chaincode 设计
##### init
	do nothing
##### invoke
###### create
创建用户并获得初始余额    
1. 传入参数 (id, balance,timestamp)  
2. 操作：  
	存入 (id, balance) -> 作为该用户balance记录，将一直存在
	存入 (id+"numoftx", 1) -> 用于存储与用户有关的交易数量的条数
	m:存入 (id+1,R+outid+balance+timestamp) -> 用于记录用户的交易记录
	t:只查创建用户，而不进行初始转账
###### transfer
转账  
1. 传入参数 (outid, inid, amount, timestamp)  
2. 存入内容:  
	修改 outid和inid balance 和txnum
	存入 (outid+numoftx, S+inid+amount+timestamp)  
	存入 (inid+numoftx, R+outid+amount+timestamp)  
##### query
###### getbalance
查询余额
1. 传入参数 id
2. 返回参数 balance
###### gettxnum
查询相关交易数量
1. 传入参数 id
2. 返回参数 num
###### history
查询历史
1. 传入参数 id和num  
2. 返回参数   
	格式为 R/S+sp+outid/inid + sp + amount + sp + timestamp    
	这里sp暂时设置为"\n" 
3. todo：检查num的合理性   

###peer接口使用
#####测试部署
	go build   
	CORE_CHAINCODE_ID_NAME=mycc CORE_PEER_ADDRESS=0.0.0.0:7051 ./chaincode  
#####使用方法
	mycc="mycc"  
	path="https://github.com/lianran/learn-chaincode/newone"  
	unixtime=$(date +%s)  

	deploy:  
	peer chaincode deploy -n ${mycc} -c '{"Args": ["init","a"]}'  

	invoke： 
	创建用户  
		peer chaincode invoke -n ${mycc} -c "{\"Args\": [\"create\", \"id1\",\"1000\",  \"${unixtime}\"]}"    
	转账
		peer chaincode invoke -n ${mycc} -c "{\"Args\": [\"transfer\", \"outid\",\"inid\", \"10\", \"${unixtime}\"]}"   

	query:
	查询历史
		peer chaincode query -n ${mycc} -c "{\"Args\": [\"history\", \"id\",\"1\"]}"
	查询余额：
		peer chaincode query -n ${mycc} -c "{\"Args\": [\"getbalance\", \"id\"]}"
	查询交易数量：
		peer chaincode query -n ${mycc} -c "{\"Args\": [\"gettxnum\", \"id\"]}"
#####连续测试命令233
	see the test.sh
