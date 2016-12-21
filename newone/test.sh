num=$1
key=$2
#node 1
ip="10.0.2.53"
port="8050"
user="admin"
pswd="Xurw3yU9zI0l"
curl -X POST -H 'Content-Type:application/json' --data "{\"enrollId\": \"${user}\",\"enrollSecret\": \"${pswd}\"}" ${ip}:${port}/registrar
for i in `seq ${num}` 
do
  ./create.sh 1000 1 1 ${i}_cp_${key}_${user} ${user} ${ip} ${port}> tmp${i}${key}.log 2>&1 &
done 

#node 2
ip="10.0.2.53"
port="9050"
user="binhn"
pswd="7avZQLwcUe9q"
curl -X POST -H 'Content-Type:application/json' --data "{\"enrollId\": \"${user}\",\"enrollSecret\": \"${pswd}\"}" ${ip}:${port}/registrar
for i in `seq ${num}` 
do
  ./create.sh 1000 1 1 ${i}_cp_${key}_${user} ${user} ${ip} ${port}> tmp${i}${key}.log 2>&1 &
done 

#node 3
ip="10.0.2.53"
port="10050"
user="diego"
pswd="DRJ23pEQl16a"
curl -X POST -H 'Content-Type:application/json' --data "{\"enrollId\": \"${user}\",\"enrollSecret\": \"${pswd}\"}" ${ip}:${port}/registrar
for i in `seq ${num}` 
do
  ./create.sh 1000 1 1 ${i}_cp_${key}_${user} ${user} ${ip} ${port}> tmp${i}${key}.log 2>&1 &
done 