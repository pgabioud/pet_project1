import json
import os
import pandas as pd



if __name__ == "__main__":
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
                        row['case'] = case
                        row['lengths'] = lengths
                        df = df.append(row, ignore_index = True)
    print(df)