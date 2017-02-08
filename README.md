# gomuni

All the italian cities in one place!

To setup the environment run the `download.sh` script.

It will create the `shp-files` folder and it will download the [latest shapefiles](http://www.istat.it/storage/cartografia/confini_amministrativi/non_generalizzati/2016/Limiti_2016_WGS84.zip) from the ISTAT website inside it.


Then to launch the server run:

```sh
go run cmd/gomuni-server/main.go
```