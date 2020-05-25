from scapy.all import *
from scapy.layers.inet import IP, TCP
import json
import os


def process_packets(dirname, output):
    data_dict = {}
    flist = os.listdir(dirname)
    IP_ADDR = ['46.4.88.92']

    for pcapfile in flist:
        lengths = []
        case = pcapfile.split('.')[0]
        data_dict[case] = {}

        pcap_path = dirname + pcapfile
        #print(pcap_path)
        packets = rdpcap(pcap_path)
        for pkt in packets:
            #pkt.show()
            #tcp_pkt = IP()/TCP()
            #print(tcp_pkt)
            #if tcp_pkt.haslayer(TCP) in packets:
                #print("here")
            print(len(pkt[TCP].payload))
            if pkt[IP].src in IP_ADDR:
                #print("there")
                lengths.append(-len(pkt))
            elif pkt[IP].dst in IP_ADDR:
                #print("everywhere")
                lengths.append(len(pkt))
                
        data_dict[case]['lengths'] = lengths
        data_dict[case]['time'] = str(round(packets[len(packets)-1].time - packets[0].time, 3))

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
