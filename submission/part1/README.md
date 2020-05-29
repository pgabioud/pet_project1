# SecretStroll

## PART 1

### A sample run of Part 1
Here we show a typical run of the system for Part 1 with the following privates attributes: "postal, tel, social_sec, sk" and public attributes "bar, gym".
A client wishing to subscribe to bar but not gym would set his public attributes to "1,0".
Here the client username is "clement"
Note that the public attributes can be of any length the server wishes to provide a service for. 

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
(client) $ python3 client.py get-pk -o key-client.pub
(client) $ python3 client.py register -a "1,0" -p key-client.pub -u "clement" -o attr.cred
(client) $ python3 client.py loc -p key-client.pub -c attr.cred -r "1,0" 46.52345 6.57890

```
Close docker:
```
$ docker-compose down
```

### Tests
To run the tests:
```
$ pytest test.py
```

### Performance evaluation
To run the performance evaluation of the computation time:
```
$ python3 perf.py
```


