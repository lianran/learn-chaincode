num=$1
key=$2
for i in `seq ${num}` 
do
  ./create.sh 1000 1 1 ${i}_cp_${key} > tmp${i}${key}.log 2>&1 &
done 
