from scapy.all import *
from scapy.layers.inet import IP, TCP
import json
import os


def process_packets(dirname, output):
    data_dict = {}
    flist = os.listdir(dirname)
    IP_ADDR = ['88.99.31.186']

    for pcapfile in flist:
        lengths = []
        times = []
        case = pcapfile.split('.')[0]
        data_dict[case] = {}

        pcap_path = dirname + pcapfile
        #print(pcap_path)
        packets = rdpcap(pcap_path)
        for pkt in packets:
            #pkt.show()
            #tcp_pkt = IP()/TCP()
            #print(tcp_pkt)
            #if tcp_pkt.haslayer(TCP):
                #print("here")
            #print(len(pkt[TCP].payload))
            init_time = packets[0].time
            if pkt[IP].src in IP_ADDR:
                #print("there")
                lengths.append(-len(pkt))
                times.append(str(round(pkt.time -  init_time, 3)))
            elif pkt[IP].dst in IP_ADDR:
                #print("everywhere")
                lengths.append(len(pkt))
                times.append(str(round(pkt.time -  init_time, 3)))
                
        data_dict[case]['lengths'] = lengths
        data_dict[case]['times'] = times

    with open(output, 'w') as outfile:
        json.dump(data_dict, outfile, sort_keys=True, indent=4)



if __name__ == "__main__":
    resultDirectory = "result_json"
    os.mkdir(resultDirectory)
    for root,d_names,f_names in os.walk("."):
        for d in d_names:
            if d != "result_json":
                dname = os.path.join(root,d)
                output = d + ".json"
                print ("#############################################################")
                print ("Process directory: " + dname)
                process_packets(dname + '/', resultDirectory + '/' + output)
                print ("Done\n")
