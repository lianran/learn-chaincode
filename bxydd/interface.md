###chaincode 设计
##### init
	do nothing
##### invoke
###### create
创建用户并获得初始余额
1. 传入参数 (id, balance,timestamp)
2. 操作：
	存入 (id, balance) -> 作为该用户balance记录，将一直存在
	存入 (R+id+timestamp,outid+balance) -> 用于记录用户的交易记录
###### transfer
转账
1. 传入参数 (outid, inid, ammount, timestamp)
2. 存入内容：
	修改 outid和inid balance
	存入 (S+outid+timestamp, inid+ammount)
	存入 (R+inid+timestamp, outid+ammount)
##### query
###### getbalance
查询余额
1. 传入参数 id
2. 返回参数 balance
###### history
查询历史
1. 传入参数 id和time（非必须）
2. 返回参数 list of history

###peer接口使用
#####测试部署
	go build 
	CORE_CHAINCODE_ID_NAME=mycc CORE_PEER_ADDRESS=0.0.0.0:7051 ./chaincode
#####调用
	
	mycc="mycc"
	path="https://github.com/lianran/learn-chaincode/newone"

	deploy:
	peer chaincode deploy -n ${mycc} -c '{"Args": ["init","a"]}'

	invoke：
	创建用户
		peer chaincode invoke -n ${mycc} -c "{\"Args\": [\"create\", \"id1\",\"1000\",  \"1480835530\"]}"  
	转账
		peer chaincode invoke -n ${mycc} -c "{\"Args\": [\"transfer\", \"outid\",\"inid\", \"10\", \"1480835530\"]}" 

	query:
	查询历史
		peer chaincode query -n ${mycc} -c "{\"Args\": [\"history\", \"id\",\"1480835530\"]}"
	查询余额：
		peer chaincode query -n ${mycc} -c "{\"Args\": [\"balance\", \"id\"]}"