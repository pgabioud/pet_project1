import numpy as np

import pandas as pd
import sklearn
from sklearn.model_selection import KFold

from sklearn.metrics import accuracy_score
from sklearn.ensemble import RandomForestClassifier
import time

def load_data():

    #import pickle
    df = pd.read_pickle("data.pkl")

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
        t = length[abs(length) != 1]
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

    max_size = np.max(np.array(size_data))
    data = np.zeros((len(df), max_size), np.float64)
    cnt = 0
    for i in df['lengths'].values:
        i = i[abs(i) != 1]

        data[cnt, :i.shape[0]] += i[:]

        cnt += 1
    times = np.expand_dims(np.array(times), 1)
    data = np.concatenate((times, data), axis=1)
    print(data.shape, max_size)
    
    return datum, data, label




datum, data, label = load_data()





skf = KFold(n_splits=10)
accuracy_model = []


clf = RandomForestClassifier(n_estimators=200, max_features=0.7, max_depth = None, min_samples_leaf = 1, min_samples_split = 5, n_jobs=-1, random_state=1, bootstrap=False)
i = 0
for train_index, test_index in skf.split(datum, label):
    X_train, X_test = datum[train_index], datum[test_index]
    y_train, y_test = label[train_index], label[test_index]

    '''
    GRIDSEARCH
    if i == 0:
        for d in [10, 20, 30, 40, 50, None]:
            for f in ['auto', 'sqrt', 0.3, 0.5, 0.6, 0.7, 0.8]:
                for min_sleaf in [1, 2, 4, 5, 10]:
                    for min_split in [2, 5, 10]:
                        for boot in [True, False]:
                            clf = RandomForestClassifier(n_estimators=200, max_features=f, max_depth = d, min_samples_leaf = min_sleaf, min_samples_split = min_split, n_jobs=-1, random_state=1)
                            model = clf.fit(X_train, y_train)
                            score = accuracy_score(y_test, model.predict(X_test), normalize=True)*100

                            print(score, d, f, min_sleaf, min_split)

                            i += 10
    '''

    print(X_train.shape)
    t0 = time.time()
    model = clf.fit(X_train, y_train)

    #show feature importance for decisions
    importances = list(model.feature_importances_)
    print(importances)

    score = accuracy_score(y_test, model.predict(X_test), normalize=True)*100
    print(score, time.time()-t0)
    accuracy_model.append(score)
    print('------------------------------------------------------------------------')
   
print(accuracy_model)

skf = KFold(n_splits=10)
accuracy_model = []

clf = RandomForestClassifier(n_estimators=500, max_features='sqrt', max_depth = 20, n_jobs=-1, random_state=1, bootstrap=False)
for train_index, test_index in skf.split(data, label):
    X_train, X_test = data[train_index], data[test_index]
    y_train, y_test = label[train_index], label[test_index]
    

    print(X_train.shape)
    t0 = time.time()
    model = clf.fit(X_train, y_train)

    score = accuracy_score(y_test, model.predict(X_test), normalize=True)*100
    print(score, time.time()-t0)
    accuracy_model.append(score)
    print('------------------------------------------------------------------------')
   
print(accuracy_model)
