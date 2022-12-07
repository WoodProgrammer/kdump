# tcp-exporter

This is a project that contains important tcp metrics like Retransmission, window scale, rst, ip_defragmentation ..etc 

You are able to run this metric exporter in all the components which is interacting with network in anyhow.

This metric exporter is responsible to listen interfaces on the workload and produce metrics based on the analysis.

## Prometheus Metrics

Metrics example this metric exporter works like that;

```sh
    #!/bin/bash

    pushd tcp-exporter
        go build .
        ./tcpdump_exporter -p 9090 -i eni-123123 -f tcp
    popd
```

After you build and run this application it will start to produce metrics like that;

# List of the Metrics

```sh
# TYPE duration_metric counter
duration_metric{dstIp="<<DST_IP_ADDR>>",srcIp="<<SRC_IP_ADDR>>"} 3.0313016523275064e+18
...
..
.
# TYPE retransmission_metric counter
retransmission_metric{dstIp="<<DST_IP_ADDR>>",srcIp="<<SRC_IP_ADDR>>"} 12
window_scale_metric{dstIp="<<DST_IP_ADDR>>",srcIp="<<SRC_IP_ADDR>>"} 5.0
rst_metric{dstIp="<<DST_IP_ADDR>>",srcIp="<<SRC_IP_ADDR>>"} 3.0
ip_df_metric{dstIp="<<DST_IP_ADDR>>",srcIp="<<SRC_IP_ADDR>>"} 1.2
zerowindow_metric{dstIp="<<DST_IP_ADDR>>",srcIp="<<SRC_IP_ADDR>>"} 0.0
ip_packet_size{}
tcp_packet_size{g}
packet_count{}
```

## Grafana Dashboard

This is the example view of the Grafana `dashboard` of this project you can check the json model of the `grafana` dashboard under `dashboard` directory.


<img src="./img/dashboard.png"></img>


## TODO;

* K8S Daemonset
* Controller
* Better dashboard