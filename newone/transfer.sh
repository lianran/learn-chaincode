num=$1 #the num of txs every round
waittime=$2 # the delay of every tx
round=$3 #the round num
key=$4 #the key for this test
source /etc/profile
echo "one round will invoke $num txs"
echo "every tx will wait $waittime s"
echo "there will have $round rounds"
echo "the key for this test is $key"
unixtime=$(date +%s)
fromid="JD"
toid="xiaoming"
date
CC_ID="0b17b27cb5d20e1a46cf65f21319d8e0c85da243b14e7a85039227e789645a8d3ce16b618faf69a955cae965b044b89dfdc83d9b2718b9c57d3065e73f491186"
for ((i=0; i < $round; i++)); do
    echo "round $i"
    for ((j=0; j < $num; j++)); do
        peer chaincode invoke -n ${CC_ID} -c "{\"Args\": [\"transfer\", \"cp_${key}_${j}\",\"${fromid}\", \"${toid}\", \"${unixtime}\"]}"  
        #echo "the key is : ${key}${i}${j}"
    sleep $waittime
    done
done
date
