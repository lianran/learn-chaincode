该脚本用于分开启动fabric的peer，其中在一台主机上运行ca和vp0，另外一台主机上运行vp1、vp2和vp3
####运行ca和vp0(./contain ca)
	1. 运行 . setenv.sh
	2. 运行 docker-compose -f four-peer-ca-caandvp0.yaml up
####运行vp1、vp2和vp3(./just peer)
	1. 设置 setenv.sh 中
		export vp0=""
		export membersrvc=""
		为ca和vp0的ip
	2. 运行 . setenv.sh
	3. 运行 docker-compose -f four-peer-ca-vp1-3.yaml up

####Attention
	注意保持./base文件中 peer-secure-pbft-base.yaml 参数      
		- CORE_PBFT_GENERAL_MODE=batch 
      	- CORE_PBFT_GENERAL_N=4 
    的一致性