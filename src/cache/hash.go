package cache

import (
    "math/rand"
    "../client"
)
/************************************************
Hash Function API

Initialization: 
    h = Make()

GetCaches(file string, clientID int) []int
    return which cache(s) a client should talk to for a particular file 

************************************************/

type Hash struct {
    numGroups        int
    clientIds        []int
    fileGroups       map[string]int // map of file to column group
    cacheIdToGroup   map[int]int
    groupToCacheIDs  map[int][]int
    replicaOrder     map[string]map[int][]int // map of file to client id to an ordering amongst replicas to check
    cacheIDs         []int
}

func (h *Hash) GetCaches(file string, clientID int) []int {
    return h.replicaOrder[file][clientID]
}

func MakeHash(numCaches int, filenames []string, n int, replication int, clients []*client.Client) *Hash {
    h := &Hash{}
    h.initializeClientIDs(clients)
    h.numGroups = numCaches / replication // number of "columns"
    h.fileGroups = MakeFileGroups(filenames, n, h.numGroups, SEED)
    h.CacheIdToGroupInit(numCaches, h.numGroups)
    h.MakeCacheOrderings(filenames)
    return h
}

func (h *Hash) MakeCacheOrderings(filenames []string) {
    replicaOrder := map[string]map[int][]int{} // filename --> client id --> ordering
    for _, file := range filenames {
        replicaOrder[file] = h.getCacheOrdersForFile(file)
    }
    h.replicaOrder = replicaOrder
}


func (h *Hash) initializeClientIDs(clients []*client.Client) {
    ids := make([]int, len(clients))
    for i, client := range clients {
        ids[i] = client.GetID()
    }
    h.clientIds = ids
}

func (h *Hash) getCacheOrdersForFile(file string) map[int][]int {
    caches := h.groupToCacheIDs[h.FileToGroup(file)]
    mapping := map[int][]int{}
    for i, clientId := range h.clientIds {
        shuffled := Shuffle(caches, i)
        mapping[clientId] = make([]int, len(shuffled))
        copy(mapping[clientId], shuffled)
    }
    return mapping
}


func (h *Hash) FileToGroup(filename string) int {
    return h.fileGroups[filename]
}

func (h *Hash) CacheIdToGroupInit(numCaches int, numGroups int) {
    mapping := SplitAmongstGroups(numCaches, numGroups)
    idToGroup := make(map[int]int)
    for id := 0; id < numCaches; id++ {
        idToGroup[id] = mapping[id]
    }
    groupToID := make(map[int][]int)
    for id, group := range idToGroup {
        groupToID[group] = append(groupToID[group], id)
    }
    h.cacheIdToGroup = idToGroup
    h.groupToCacheIDs = groupToID
}

func SplitAmongstGroups(n int, numGroups int) []int {
    mapping := make([]int, n)
    minpergroup := n / numGroups

    for i := 0; i < numGroups; i++ {
        for j := 0; j < minpergroup; j++ {
            mapping[minpergroup*i + j] = i
        }
    }
    for i := 0; i < n - (numGroups * minpergroup); i++ {
        mapping[minpergroup*numGroups + i] = i
    }
    return mapping

}

func Shuffle(slice []int, seed int) []int {
    rand.Seed(int64(seed))
    rand.Shuffle(len(slice), func(i, j int) {
        slice[i], slice[j] = slice[j], slice[i]
    })
    return slice
}

func MakeFileGroups(filenames []string, n int, numGroups int, seed int) map[string]int {
    mapping := SplitAmongstGroups(n, numGroups)
    mapping = Shuffle(mapping, seed)
    fileGroups := make(map[string]int)
    for i := 0; i < n; i++ {
        fileGroups[filenames[i]] = mapping[i]
    }
    return fileGroups
}


