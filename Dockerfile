FROM alpine

# Make the source code path
RUN mkdir -p /go/src/github.com/enrichman/gomuni
RUN mkdir /root/shp-files

# Add sources and install the app
ADD cmd/gomuni-server/gomuni-server /root

# Add shapefiles and env file to the $HOME
ADD shp-files /root/shp-files
ADD .env.docker /root/.env

EXPOSE 8080
WORKDIR /root

CMD ["./gomuni-server"]
