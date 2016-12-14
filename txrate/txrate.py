import urllib.request as urlreq
import json
import time

#some var
url_blocks = "http://162.105.30.216:37050/chain/blocks/"
url_height = "http://162.105.30.216:37050/chain"
list_blocks = []
#{"num":num, "time":time, "blocksnum": i}

def get_blocks_txsnum_time(url):
    fp = urlreq.urlopen(url)
    mybytes = fp.read()
    mystr = mybytes.decode("utf8")
    fp.close()
    data = json.loads(mystr)
    #get the time
    timestamp = data['nonHashData']['localLedgerCommitTimestamp']['seconds'] + data['nonHashData']['localLedgerCommitTimestamp']['nanos']*1.0/1e9
    #print(time)
    #get tne num of txs
    try:
    	numoftxs = len(data['transactions'])
    except:
    	numoftxs = 0
    	print(data)
    #print(numoftxs)
    return (numoftxs, timestamp)
def get_height(url):
    fp = urlreq.urlopen(url)
    mybytes = fp.read()
    mystr = mybytes.decode("utf8")
    fp.close()
    data = json.loads(mystr)
    return data['height']
def cal_rate(k):
    len_list = len(list_blocks)
    endtime = list_blocks[len_list-1]['timestamp']
    count = list_blocks[len_list-1]['num']
    for i in range(len_list-2, len_list-k-1, -1):
        if i < 0:
            break
        count += list_blocks[i]['num']
    startime = list_blocks[len_list-k-1]['timestamp']
    print("txRate is:" + str(count/(endtime-startime)))
    
#init get the pre 10 blocks
pre_height = get_height(url_height)
for i in range(pre_height-10, pre_height):
    if i < 0:
        break
    (num, timestamp) = get_blocks_txsnum_time(url_blocks + str(i))
    list_blocks.append({"num":num, "timestamp":timestamp, "blocksnum": i})
    print({"num":num, "timestamp":timestamp, "blocksnum": i})

print("the pre"),
max_rate = cal_rate(5)
#loop to get the blcoks and some thing
cnt = 0
while True:
    now_height = get_height(url_height)
    if now_height == pre_height:
        if cnt%10 == 0:
            print(str(time.time()) + " : no new blocks")
    else:
        #get the new blockscee
        for i in range(pre_height, now_height):
            (num, timestamp) = get_blocks_txsnum_time(url_blocks + str(i))
            list_blocks.append({"num":num, "timestamp":timestamp, "blocksnum": i})
            print("new block:" + str({"num":num, "timestamp":timestamp, "blocksnum": i}))
        #show the txRate
        pre_height = now_height
        now_rate = cal_rate(3)
        if now_rate > max_rate:
        	max_rate = now_rate
        print("the max tx_rate is :" + str(max_rate))
        cnt = 0
    time.sleep(0.5)
    cnt += 1
        
        