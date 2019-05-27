## MeshRPC Benchmark Suite

There are two Docker images involved in benchmarks:

* `docker.direct/meshrpc/benchmark/bench_api`
* `docker.direct/meshrpc/benchmark/peer`

Running a cluster:

```
$ docker stack deploy -c docker-compose.yml meshrpc-bench
Creating network meshrpc-bench_meshnet
Creating service meshrpc-bench_bench-api
Creating service meshrpc-bench_peer-A
Creating service meshrpc-bench_peer-B
Creating service meshrpc-bench_peer-C
Creating service meshrpc-bench_peer-D
Creating service meshrpc-bench_peer-E

$ docker service logs meshrpc-bench_bench-api
meshrpc-bench_bench-api.1.bifrqj253pne@linuxkit-025000000001    | 2019/05/25 16:13:14 wait for peers is over [A B C D E] initializing client to A
```

### Non-RPC Results

Macbook Pro 2014 (2,8 GHz Intel Core i5): **33,3 µs** per call:

```
$ curl -d'{"queue": ["A","B","C","D","E"], "limit": 1000}' http://localhost:8282/staticPeer/hop
{
    "State": {
        "Start": 1558817006608054800,
        "Current": 1558817006641384200,
        "Last": "E",
        "Queue": [
            "A",
            "B",
            "C",
            "D",
            "E"
        ],
        "QueueIdx": 0,
        "Limit": 1000,
        "HopCount": 1000,
        "HopLatency": 33329
    },
    "Error": null
}
```

Linux 4.15.0-47-generic (Intel Xeon CPU E3-1270 v6 @ 3.80GHz): **0,235 µs** per call:

```
{
    "State": {
        "Start": 1558994575386261949,
        "Current": 1558994575388621468,
        "Last": "E",
        "Queue": [
            "A",
            "B",
            "C",
            "D",
            "E"
        ],
        "QueueIdx": 0,
        "Limit": 10000,
        "HopCount": 10000,
        "HopLatency": 235
    },
    "Error": null
}
```

### RPC Results

Macbook Pro 2014 (2,8 GHz Intel Core i5): **1,8 ms** per RPC call:

```
$ curl -d'{"queue": ["A","B","C","D","E"], "limit": 1000}' http://localhost:8282/peer/hop
{
    "State": {
        "Start": 1558801539512806100,
        "Current": 1558801541410891900,
        "Last": "E",
        "Queue": [
            "A",
            "B",
            "C",
            "D",
            "E"
        ],
        "QueueIdx": 0,
        "Limit": 1000,
        "HopCount": 1000,
        "HopLatency": 1898085
    },
    "Error": null
}
```

Linux 4.15.0-47-generic (Intel Xeon CPU E3-1270 v6 @ 3.80GHz): **197,4 µs** per RPC call:

```
{
    "State": {
        "Start": 1558994662029384673,
        "Current": 1558994664003940891,
        "Last": "E",
        "Queue": [
            "A",
            "B",
            "C",
            "D",
            "E"
        ],
        "QueueIdx": 0,
        "Limit": 10000,
        "HopCount": 10000,
        "HopLatency": 197455
    },
    "Error": null
}
```

