num=$1 #the num of txs every round
waittime=$2 # the delay of every tx
round=$3 #the round num
key=$4 #the key for this test
source /etc/profile
echo "one round will invoke $num txs"
echo "every tx will wait $waittime s"
echo "there will have $round rounds"
echo "the key for this test is $key"
date
CC_ID="mycc"
for ((i=0; i < $round; i++)); do
    echo "round $i"
    for ((j=0; j < $num; j++)); do
        peer chaincode invoke -n ${CC_ID} -c "{\"Args\": [\"create\", \"${key}${i}${j}\",\"JD\", \"PKU\", \"1480835530\", \"metadata\"]}"  
        #echo "the key is : ${key}${i}${j}"
    sleep $waittime
    done
done
date
