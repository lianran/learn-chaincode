##消息存证系统
该智能合约将完成消息的存证和取证，所有的聊天记录的hash值按照用户接收和发送先后顺序，存证到区块链对应的db当中，从而方便后续的查询。另外，所有的消息的hash都被记录到了区块链上，db可以从该区块链上的交易记录进行恢复。
###存储内容
--
A. 聊天记录数据 （userId1 + msgId, "F"/"T" + userId2 + timestamp + msgHash) 
	
	userId: 用户id
	msgId: 该条信息对应该用户的消息编号，采用自增的方式存储
	F/T：F表示msg From userId2; T 表示msg Send To userId2
	timestamp: 消息时间戳
	msgHash: 消息Hash值
		
B. 用户消息数目记录 （userId, msgId)

	userId： 用户id
	msgId：用户现在的小心编号，采用自增的方式进行记录，实际上是用户相关消息的记数
	
###chaincode接口设计
--
####存证 save   
存证是消息hash上链和存入是数据库的过程。  

1. 传入参数

		消息发送方fromid，接收方toid，时间戳timstamp，消息hash值msgHash  
 
2. 数据处理流程  

	a. 获取收发方用户消息数目，如果获取失败，设置msgid为0  
	
	b. 存证消息记录  
	
		存入 (fromId + fromMsgId, "T" + toId + timestamp + msgHash)   
		存入 (toId + toMsgId, "F" + fromId + timestamp + msgHash)   
		
	c. 修改用户消息数目记录  
	
		存入 (fromId, (fromMsgId+1))  
		存入 (toId, (toMsgId+1))

####取证  
取证是对已有消息真实性的验证过程。

#####取证某个用户某条消息Hash的正确性
1. 传入数据

		用户名 userId, 大致时间 timestamp， 消息哈希msgHash
		
2. 返回数据

		True: 该消息Hash正确，在链上有相应存证
		False：该消息Hash与记录不符
					
#####取证某个用户某条消息的真实性
1. 传入数据

		用户名 userId, 大致时间 timestamp， 消息内容msgContent
		
2. 返回数据

		True: 该消息真实存在，未经篡改
		False：该消息与记录不符
		
#####取证某用户某个时间点的消息hash
1. 传入数据
	
		用户名userId, 时间timestamp
		
2. 返回数据

		列表：
			格式为 “T”/"F" + toId/fromId + timestamp + msgHash

#####取证某个用户某个时间段的消息hash
1. 传入数据

		用户名 userId, 开始时间startTime， 结束时间endTime
		
2. 返回数据

		列表：
			格式为 “T”/"F" + toId/fromId + timestamp + msgHash
			
#####取证某个用户最近的消息hash
1. 存入数据

		用户名 userId, 获取消息数目（默认为是10条）
		
2. 返回数据

		列表：
			格式为 “T”/"F" + toId/fromId + timestamp + msgHash


