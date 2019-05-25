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
```
