# DECA NEWS
website for making and sharing news and leaks about the [decagear](https://deca.net) headset

#### live at [decafans.com](https://decafans.com)

## building and running
```shell
docker build --build-arg "DB_CREDS={{Firebase auth creds json contents with \" s escaped}}" \
 --build-arg D_KEY={{discord client id}} --build-arg D_SECRET={{discord app secret}} \
  --build-arg HOST_PATH={{https://decafans.com || or path to your url}} \
  --build-arg RANDOM_SECRET={{idk make something up for this its for the local cookie store}} -t decafans .
```
then
```shell
docker run -d --restart unless-stopped -p 5000:5000 decafans
```
or in testing use
```shell
docker run -p 5000:5000 -v src/templates:templates decafans
```