echo "plz enter the billid and the owner"
peer chaincode query -n mycc -c "{\"Args\": [\"getbill\", \"$1\",\"$2\"]}"