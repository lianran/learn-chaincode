num=$1
key=$2
user=$3
pswd=$4
#peer network login ${user} -p ${pswd}
curl -X POST -H 'Content-Type:application/json' --data "{\"enrollId\": \"${user}\",\"enrollSecret\": \"${pswd}\"}" localhost:8050/registrar
for i in `seq ${num}` 
do
  ./create.sh 1000 1 1 ${i}_cp_${key} ${user}> tmp${i}${key}.log 2>&1 &
done 
