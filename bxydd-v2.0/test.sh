mycc="mycc"
echo "deploy the chaincode"
peer chaincode deploy -n ${mycc} -c '{"Args": ["init","a"]}'
echo "create the a,b,c with balance 1000"
unixtime=$(date +%s)
peer chaincode invoke -n ${mycc} -c "{\"Args\": [\"create\", \"a\",\"1000\",  \"${unixtime}\"]}"
unixtime=$(date +%s)
peer chaincode invoke -n ${mycc} -c "{\"Args\": [\"create\", \"b\",\"1000\",  \"${unixtime}\"]}"
unixtime=$(date +%s)
peer chaincode invoke -n ${mycc} -c "{\"Args\": [\"create\", \"c\",\"1000\",  \"${unixtime}\"]}"
echo "get the balance of a"
peer chaincode query -n ${mycc} -c "{\"Args\": [\"getbalance\", \"a\"]}"
echo "get teh balance of b"
peer chaincode query -n ${mycc} -c "{\"Args\": [\"getbalance\", \"b\"]}"
echo "get the history of a"
peer chaincode query -n ${mycc} -c "{\"Args\": [\"history\", \"a\",\"1\"]}"
echo "transfer a to b 100"
unixtime=$(date +%s)
peer chaincode invoke -n ${mycc} -c "{\"Args\": [\"transfer\", \"a\",\"b\", \"10\", \"${unixtime}\"]}" 
echo "get the balance of a and b"
peer chaincode query -n ${mycc} -c "{\"Args\": [\"getbalance\", \"a\"]}"
peer chaincode query -n ${mycc} -c "{\"Args\": [\"getbalance\", \"b\"]}"
echo "get the txsnum of a and b"
peer chaincode query -n ${mycc} -c "{\"Args\": [\"gettxnum\", \"a\"]}"
peer chaincode query -n ${mycc} -c "{\"Args\": [\"gettxnum\", \"b\"]}"
#wait for the tx executed
sleep 2
echo "get the history of a and b "
peer chaincode query -n ${mycc} -c "{\"Args\": [\"history\", \"a\",\"2\"]}"
peer chaincode query -n ${mycc} -c "{\"Args\": [\"history\", \"b\",\"2\"]}"
