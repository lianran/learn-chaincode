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
id=0
date
CC_ID="db4fa76606a419673dfc30c929b3950f37c9525ce9739d510b4d4e6529a1986e898435ce5cb36725dd48eb9b48fa8d9b5b18417070c1079db9c8045698ff760c"
for ((i=0; i < $round; i++)); do
    echo "round $i"
    for ((j=0; j < $num; j++)); do
        #peer chaincode invoke -n ${CC_ID} -c "{\"Args\": [\"create\", \"cp_${i}_${key}_${j}\",\"JD\", \"PKU\", \"${unixtime}\", \"price:${i}${j}\"]}"  -u ${user} &
        curl -X POST -H 'Content-Type:application/json' --data "{\"jsonrpc\":\"2.0\",\"method\":\"invoke\",\"params\":{\"type\":1,\"chaincodeID\":{\"name\":\"${CC_ID}\"},\"ctorMsg\":{\"args\":[\"create\",\"cp_${i}_${key}_${j}\",\"JD\",\"PKU\",\"${unixtime}\",\"price:${i}${j}\"]},\"secureContext\":\"${user}\"},\"id\":${id}}" localhost:7050/chaincode
        #echo "the key is : ${key}${i}${j}"
    #sleep $waittime
    id++
    done
done
date
