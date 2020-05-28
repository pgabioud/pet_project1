import numpy as np

import pandas as pd
import sklearn
from sklearn.model_selection import KFold

from sklearn.metrics import accuracy_score
from sklearn.ensemble import RandomForestClassifier


def load_data():

    #import pickle
    df = pd.read_pickle("data_wtime.pkl")

    #labels
    label = np.array(df['case'].values)

    #timestamps
    timestamps = np.array(df['times'].values)

    #compute time difference between each packet
    tdeltas = []
    for i in range(len(timestamps)):
        idx = []
        for j in range(len(timestamps[i])):
            if j == 0:
                last = 0
            else:
                last = timestamps[i][j-1]
            idx.append(timestamps[i][j] - last)
        tdeltas.append(idx)


    #statistics         /       10 features
    size_data = []              #number of packets exchanged
    avg_length_data = []        #average packet length
    max_length_data = []        #maximum packet length
    min_length_data = []        #minimum packet length
    third_length_data = []      #first quartile of packet length distribution 
    percent_server_data = []    #percentage of packets that were from the server side
    times = []                  #total packet exhcange time
    times_avg = []              #average difference between packets
    times_max = []              #max time difference between packets
    times_min = []              #minimum time difference

    lengths = np.array(df['lengths'].values)
    for i in range(len(lengths)):
        length = lengths[i]

        #remove acks
        t = length[abs(length) != 54]
        size_data.append(len(t))
        p = len(t)
        #incase filtering empties the row
        if p == 0:
            p = 0
            avg_length_data.append(0)
            max_length_data.append(0)
            min_length_data.append(0)
            third_length_data.append(0)
            times.append(0)
            times_avg.append(0)
            times_max.append(0)
            times_min.append(0)
        
        else:
            p = len(t[t < 0]) / len(t)
            percent_server_data.append(p)
            avg_length_data.append(np.mean(abs(t)))
            max_length_data.append(np.max(abs(t)))
            min_length_data.append(np.min(abs(t)))
            third_length_data.append(np.percentile(abs(t), 25))
            times.append(timestamps[i][-1])
            times_avg.append(np.mean(tdeltas[i]))
            times_max.append(np.max(tdeltas[i]))
            times_min.append(np.min(tdeltas[i]))


    datum = np.array([size_data, avg_length_data, max_length_data,
                     min_length_data, third_length_data, percent_server_data, times, times_max, times_min, times_avg])

    #to have (10000, 10)
    datum = datum.T

    return datum, label




datum, label = load_data()





skf = KFold(n_splits=10)
accuracy_model = []
clf = RandomForestClassifier(
    n_estimators=1000, max_features=7, max_depth = 20, n_jobs=-1, random_state=1)


for train_index, test_index in skf.split(datum, label):

    X_train, X_test = datum[train_index], datum[test_index]
    y_train, y_test = label[train_index], label[test_index]

    print(X_train.shape)

    model = clf.fit(X_train, y_train)

    #show feature importance for decisions
    importances = list(model.feature_importances_)
    print(importances)

    score = accuracy_score(y_test, model.predict(X_test), normalize=True)*100
    print(score)
    accuracy_model.append(score)
    print('------------------------------------------------------------------------')
   
print(accuracy_model)
