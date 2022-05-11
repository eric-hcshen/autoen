### Steps to run this sample:
1) You need a Temporal service running. See details in README.md
2) Run the following command to start the worker
```shell script
go run parallel/worker/main.go
```
3) Run the following command to start parallel workflow
```shell script
go run parallel/starter/main.go
```

default:~/Personal_Data/work$ cat tableau.md 
wget https://downloads.tableau.com/esdalt/2022.1.0/tableau-server-container-setup-tool-2022.1.0.tar.gz
wget https://downloads.tableau.com/esdalt/2022.1.1/tableau-server-2022-1-1.x86_64.rpm

docker run -e LICENSE_KEY="TS41-299F-6260-0E46-37A0" -p 8080:8080 -e TABLEAU_USERNAME=admin -e TABLEAU_PASSWORD=password -d tableau_server_image:20221.22.0415.1144


./build-image --accepteula -i ../tableau-server-2022-1-1.x86_64.rpm 

docker exec -it c1bb60928faf tsm status -V

docker exec -it c1bb60928faf /bin/sh -c 'cat $DATA_DIR/supervisord/run-tableau-server.log'[1]+  Done                    code .
