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
1. 传入参数 (outid, inid, amount, timestamp)
2. 存入内容：
	修改 outid和inid balance
	存入 (S+outid+timestamp, inid+amount)
	存入 (R+inid+timestamp, outid+amount)
##### query
###### getbalance
查询余额
1. 传入参数 id
2. 返回参数 balance
###### history
查询历史
1. 传入参数 id和time（非必须）
2. 返回参数 list of history
	单条格式为 R/S + sp + inid/outid + sp + timestamp, outid/inid + sp + amount
	这里sp暂时设置为"\n"

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
		peer chaincode query -n ${mycc} -c "{\"Args\": [\"history\", \"id\",\"233\"]}"
	查询余额：
		peer chaincode query -n ${mycc} -c "{\"Args\": [\"getbalance\", \"id\"]}"
#####连续测试命令233
	see the test.sh
