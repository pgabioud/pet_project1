import torch
import numpy as np
from torch.nn import functional as F
from torch import nn, optim
import math
import pandas as pd
import sklearn
from sklearn.utils import shuffle
from sklearn.model_selection import KFold
import time

from sklearn.metrics import accuracy_score
from sklearn.ensemble import RandomForestClassifier

# just to have some data
def load_data():
    df = pd.read_pickle("data.pkl")
    df = shuffle(df, random_state=1)

    label = torch.LongTensor(df['case'].values)
    timestamps = np.array(df['times'].values)
    
    
    label = label - 1
    size_data = []
    avg_length_data = []
    max_length_data = []
    min_length_data = []
    percent_server_data = []
    for length in df['lengths']:
        #t = length[abs(length) != 597]
        t = length
        size_data.append(len(t))
        percent_server_data.append(len(t[t < 0]) / len(t))
        avg_length_data.append(np.mean(abs(t)))
        max_length_data.append(np.max(abs(t)))
        min_length_data.append(np.min(abs(t)))
    max_size = np.max(np.array(size_data))
    
    

    datum = np.array([size_data, avg_length_data, max_length_data, min_length_data, percent_server_data, timestamps])


    data = np.zeros((len(df), max_size), np.float64)
    cnt = 0
    for i in df['lengths'].values:
        i = i[abs(i) != 597]
        
        data[cnt, :i.shape[0]] += i[:]
        
        cnt += 1
    data = torch.FloatTensor(data)
    data = torch.cat((torch.FloatTensor(timestamps).unsqueeze(dim = 1), data), dim = 1)
    print(label.size(), data.size())
    
    datum = torch.Tensor(datum.T)
    
    return datum, data, label

# train function using sgd with cross entropy loss


def train(model, train_input, train_target, learning_rate, batch_size):

    criterion = nn.CrossEntropyLoss()
    opt = torch.optim.Adam(model.parameters(), lr=learning_rate, weight_decay = 1e-4)
    epochs = 100
    for e in range(epochs):
        sum_loss = 0

        for b in range(0, train_input.size(0), batch_size):

            pred = model(train_input.narrow(0, b, batch_size))
           
            loss = criterion(pred, train_target.narrow(0, b, batch_size))
            sum_loss += loss.item()

            opt.zero_grad()
            loss.backward()
            sum_loss += loss.item()

            opt.step()
        if e % 2 == 0:
            print(e, sum_loss)


def compute_nb_errors(model, input, target, batch_size):

    nb_errors = 0
    for b in range(0, input.size(0), batch_size):

        predicted = model(input.narrow(0, b, batch_size))
        _, predicted = predicted.max(1)

        for k in range(batch_size):
            if target[b + k] != predicted[k]:
                nb_errors += 1
    print("accuracy = {:0.2f}".format(1 - nb_errors/len(target)))
    return nb_errors




datum, data, label = load_data()

print(data[0])



skf = KFold(n_splits=10)
accuracy_model = []
clf = RandomForestClassifier(random_state=1)
for train_index, test_index in skf.split(datum, label):
    
    '''
    model = nn.Sequential(nn.Linear(datum.size(1), 16),
                      nn.ReLU(),
                      nn.Linear(16, 32),
                      nn.ReLU(),
                      nn.Linear(32, 64),
                      nn.ReLU(),
                      nn.Linear(64, 256),
                      nn.ReLU(),
                      nn.Linear(256, 256),
                      nn.ReLU(),
                      nn.Linear(256, 256),
                      nn.ReLU(),
                      nn.Linear(256, 256),
                      nn.ReLU(),
                      nn.Linear(256, 256),
                      nn.ReLU(),
                      nn.Linear(256, 100),

                      nn.Softmax(dim=1))
    '''
    X_train, X_test = datum[train_index], datum[test_index]
    y_train, y_test = label[train_index], label[test_index]

    print(X_train.size())
    
    model2 = clf.fit(X_train, y_train)
    accuracy_model.append(accuracy_score(y_test, model2.predict(X_test), normalize=True)*100)
    '''
    train(model, X_train, y_train, 1e-3, 100)
    compute_nb_errors(model, X_test, y_test, 100)
    '''
print(accuracy_model)