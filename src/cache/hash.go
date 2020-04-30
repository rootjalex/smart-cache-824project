package cache

import (
)
/************************************************
Hash Function API

Initialization: 
    h = Make()

GetCache(file string, clientID int) []int
    return which cache(s) a client should talk to for a particular file 

************************************************/

type Hash struct {
    fileGroups      map[string]int // map of file to column group
    replicaOrdering map[int][]int // map of client id to an ordering amongst replicas to check
    cacheIDs        []int
}

func MakeHash(numCaches int, numFiles int, replication int) *Hash {
    h := &Hash{}
    // TODO: actual initilization
    return h
}

func (h *Hash) GetCache(file string, clientID int) []int {
    // TODO: fill this in
    return []int{1, 2}
}
