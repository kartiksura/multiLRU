# multiLRU
Implements Key value cache server.
There are two implementations supported:
## Simple Store
It is implemented as a simple Least Recently used data structure. It uses golang native heap for maintaining the queue of the objects to be evicted once the store is full. 

## Concurrent Store
It is a scalable key value store, suitable for multi-core CPUs. It is a sharded store, with the number of shards to being equal to the number of CPUs.
We use CRC32 for sharding algorithm. It gives a decent spread across all the shards and is simple to implement (we use std library)

## Connectors
Connectors enable various ways for clients to use the key value server.
Right now the TCP connector has been implemented. 

## Installation Instructions
go build
./multiLRU 

The key value server runs on port 61000

netcat can be used for operations on the tcp server:
nc 127.0.0.1  "61000" -C
> set kartik 4
> sura
< OK
> get kartik
< VALUE 4
< sura

> delete kartik
< OK
> get kartik
< VALUE 0

Additional commands
stats
{Sets:1 Gets:2 Success:1 MemoryUsed:0 Capacity:100}

contents
K: kartik V: sura
K: palo V: alto
K: san V: francisco



