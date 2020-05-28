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
            if pkt.haslayer(TCP):
                #print(len(pkt[TCP].payload))
                init_time = packets[0].time
                if pkt[IP].src in IP_ADDR:
                    times.append(str(round(pkt.time - init_time, 3)))
                    if(len(pkt[TCP].payload)) == 0:
                        lengths.append(-1)
                    else:         
                        lengths.append(-len(pkt[TCP].payload))
                elif pkt[IP].dst in IP_ADDR:
                    times.append(str(round(pkt.time - init_time, 3)))
                    if(len(pkt[TCP].payload)) == 0:
                        lengths.append(1)
                    else:   
                        lengths.append(len(pkt[TCP].payload))
                
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
