# DECA NEWS
website for making and sharing news and leaks about the [decagear](https://deca.net) headset

#### live at [decafans.com](https://decafans.com)

## building and running
```shell
docker build --build-arg LOCAL_FILE={{true || [false] -- optinal}} -t decafans .
```
if `LOCAL_FILE` is enabled, you need to make a file called `.env` in the src directory 
using the data in [example.env](https://github.com/zivoy/decaNewsSite/blob/master/example.env).

otherwise, you need to supply the `LOCATION`, `SERVER_PASSWORD` and `FILE_PASSWORD` environment 
variables on runtime so that the server knows where to get the file from and decrypt it

then
```shell
docker run -d --restart unless-stopped -p 5000:5000 --env-file docker.env decafans
```
or in testing use
```shell
docker run -p 5000:5000 -it -v $(pwd)/src/templates:/app/decafans-server/templates decafans  
```