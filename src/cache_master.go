package cache

/************************************************
Cache Master supports

m = Make(caches Cache)
    Initialize a cache master with a list of caches



*************************************************/


type CacheMaster {
    mu sync.Mutex
}


func (m *CacheMaster) Make(caches []Cache) {

}


