num=$1 #the num of txs every round
waittime=$2 # the delay of every tx
round=$3 #the round num
key=$4 #the key for this test
user=$5
source /etc/profile
echo "one round will invoke $num txs"
echo "every tx will wait $waittime s"
echo "there will have $round rounds"
echo "the key for this test is $key"
unixtime=$(date +%s)
date
CC_ID="mycc"
for ((i=0; i < $round; i++)); do
    echo "round $i"
    for ((j=0; j < $num; j++)); do
        #peer chaincode invoke -n ${CC_ID} -c "{\"Args\": [\"create\", \"cp_${i}_${key}_${j}\",\"JD\", \"PKU\", \"${unixtime}\", \"price:${i}${j}\"]}"  -u ${user} &
        curl -X POST -H 'Content-Type:application/json' --data "{\"jsonrpc\":\"2.0\",\"method\":\"invoke\",\"params\":{\"type\":1,\"chaincodeID\":{\"name\":\"${CC_ID}\"},\"ctorMsg\":{\"args\":[\"create\",\"cp_${i}_${key}_${j}\",\"JD\",\"PKU\",\"${unixtime}\",\"price:${i}${j}\"]},\"secureContext\":\"${user}\"},\"id\":1}" localhost:7050/chaincode
        #echo "the key is : ${key}${i}${j}"
    #sleep $waittime
    done
done
date
