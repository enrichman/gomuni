#!/bin/bash

cp .env.example .env
mkdir shp-files
curl -O http://www.istat.it/storage/cartografia/confini_amministrativi/non_generalizzati/2016/Limiti_2016_WGS84.zip
unzip Limiti_2016_WGS84.zip -d shp-files
rm Limiti_2016_WGS84.zip

go get github.com/tools/godep
godep restore ./...
