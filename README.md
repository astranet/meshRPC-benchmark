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

```
$ curl -d'{"queue": ["A","B","C","D","E"], "limit": 1000}' http://localhost:8282/peer/hop
```

**1.8 ms** per call is the current latency.

```json
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

To compare with just a function invocation (no network and RPC involved):

```
curl -d'{"queue": ["A","B","C","D","E"], "limit": 1000}' http://localhost:8282/staticPeer/hop
```

Thas is **33µs**:

```
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

