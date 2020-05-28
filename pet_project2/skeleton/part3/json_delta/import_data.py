import json
import os
import pandas as pd
import torch
import sklearn


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
                    for case, data in data_dict.items():
                        row['case'] = int(case)
                        length = []
                        time = []
                        for l in data['lengths']:
                            length.append(float(l))
                        for t in data['times']:
                            time.append(float(t))
                        row['lengths'] = np.array(length)
                        row['times'] = np.array(time)
                        df = df.append(row, ignore_index = True)
    with pd.option_context('display.expand_frame_repr', False):
        print(df)
    return df

if __name__ == "__main__":
    df = import_data()
    
    df.to_pickle("data.pkl")
    
