import json
import os
import pandas as pd
import torch
import sklearn

from sklearn.utils import shuffle

import numpy as np
def import_data():
    resultDirectory = "pickle"
    #os.mkdir(resultDirectory)
    df = pd.DataFrame()
    for _,_,f_names in os.walk("."):
        for fname in f_names:
            if(fname != "import_data.py"):
                row = {}
                row['fname'] = os.path.basename(fname)
                with open(fname) as f:
                    data_dict = json.loads(f.read())
                    for case, lengths in data_dict.items():
                        row['case'] = int(case)
                        length = []
                        for i in lengths:
                            length.append(float(i))
                        row['lengths'] = np.array(length)
                        df = df.append(row, ignore_index = True)
    return shuffle(df, random_state = 1)

if __name__ == "__main__":
    df = import_data()
    
    label = torch.tensor(df['case'].values)
    size_data = []
    for length in df['lengths']:
        size_data.append(len(length))
    max_size = np.max(np.array(size_data))

    
    data = np.zeros((3000, max_size), np.float64)
    cnt = 0
    for i in df['lengths'].values:
        
        data[cnt, :i.shape[0]] += i[:]
        cnt += 1
    
    
    data = torch.FloatTensor(data)

    print(label.size(), data.size())