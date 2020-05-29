# SecretStroll

## PART 3

### Data collection
First, copy the "collection.sh" file into the "part1" repository.
The following steps should create 100 distinct repository, each containing 100 pcap files.

Initialization:
```
Open a shell
$ cd skeleton
$ docker-compose build
$ docker-compose up -d
```

Server side:

```
Open a shell
$ cd skeleton
$ docker exec -it cs523-server /bin/bash
(server) $ cd /server
(server) $ python3 server.py gen-ca -a "postal,tel,social_sec,sk,bar,gym" -s key.sec -p key.pub
(server) $ python3 server.py run -s key.sec -p key.pub
```


Client side:

```
Open a shell
$ cd skeleton
$ docker exec -it cs523-client /bin/bash
(client) $ cd /client
(client) $ bash collection.sh
```

### Feature exctraction
Copy the "process_\pcap\_payload.py" into the repository with all the subfolders containing the pcap files.
Then run:

$ python3 process_\pcap\_payload.py


Copy in the created "result_json" folder the "import\_data.py" file, and then run:

$ python3 import\_data.py


Note that there is already a sample "json" repository containing few example of json files.


### Classification
Copy the created pickle file into the "classifier" repository.
Note that we provide the processed pickle file already in the classifier repository.
In order to run the classification with the 2 models run:
$ python3 classifier.py
