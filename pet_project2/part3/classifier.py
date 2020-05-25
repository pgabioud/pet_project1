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

# just to have some data
def load_data():
    df = pd.read_pickle("data.pkl")
    df = shuffle(df, random_state=1)

    label = torch.LongTensor(df['case'].values)
    label = label - 1
    size_data = []
    for length in df['lengths']:
        t = length[abs(length) != 597]
        t = t[abs(t) != 54]
        size_data.append(len(t))
    max_size = np.max(np.array(size_data))
    

    data = np.zeros((len(df), max_size), np.float64)
    cnt = 0
    for i in df['lengths'].values:
        i = i[abs(i) != 597]
        i = i[abs(i) != 54]
        data[cnt, :i.shape[0]] += i[:]
        cnt += 1

    data = torch.FloatTensor(data)
    
    print(label.size(), data.size())

    
    return data, label

# train function using sgd with cross entropy loss


def train(model, train_input, train_target, learning_rate, batch_size):

    criterion = nn.CrossEntropyLoss()
    opt = torch.optim.Adam(model.parameters(), lr=learning_rate)
    epochs = 20
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
    print(nb_errors)
    return nb_errors




data, label = load_data()

print(data[0])
model = nn.Sequential(nn.Linear(data.size(1), 1024),
                      nn.ReLU(),
                      nn.Linear(1024, 1024),
                      nn.ReLU(),
                      nn.Linear(1024, 1024),
                      nn.ReLU(),
                      nn.Linear(1024, 1024),
                      nn.ReLU(),
                      nn.Linear(1024, 512),
                      nn.ReLU(),
                      nn.Linear(512, 512),
                      nn.ReLU(),
                      nn.Linear(512, 256),
                      nn.ReLU(),
                      nn.Linear(256, 100),

                      nn.Softmax(dim=1))

skf = KFold(n_splits=10)
t0 = time.time()
for train_index, test_index in skf.split(data, label):
    X_train, X_test = data[train_index], data[test_index]
    y_train, y_test = label[train_index], label[test_index]
    print(X_train.size())
    train(model, X_train, y_train, 1e-4, 100)
    compute_nb_errors(model, X_test, y_test, 100)
t1 = time.time()
print(t1-t0)
compute_nb_errors(model, data, label, 100)
