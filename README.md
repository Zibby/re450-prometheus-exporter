# re450-exporter
Prometheus endpoint for tplink re450 wifi repeater  written in go exposed on port 8089.

* Currently will not work for you as cannot generate cookies or connect to repeaters not on 192.168.0.87 *

## deploy
~~~bash
docker run -d \
  --name=re450-exporter \
  -p 8089:8089
  -e TPLINK_ADDR="${IP_ADDRESS_OF_PLUG}" # not yet implemented\
  zibby/re450-exporter
~~~

browse localhost:8089/metrics


## todo:
- [] tests
- [] jenkins file rewrite
- [] add in env var for IP addr
- [] add in env var for cookie or pw?
- [] add in cookie generation

