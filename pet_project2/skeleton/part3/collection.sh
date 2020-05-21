#/usr/bin/bash


python3 /client/client.py get-pk -o key-client.pub
python3 /client/client.py register -a "" -p key-client.pub -u "oss117" -o attr.cred
for i in {1..30} #number of sample per case
do
    echo cycle:  $i
	d=`date "+%H%M%S"`
	mkdir /client/part3/$d
	sleep 3
	for j in {1..100}
	do	
		echo $j
		tcpdump host 172.18.0.2 and host 46.4.88.92 -w /client/part3/$d/$j.pcap &
		python3 /client/client.py grid -p key-client.pub -c attr.cred -r "" -t $j
		sleep 1
		pkill tcpdump
	done
done
