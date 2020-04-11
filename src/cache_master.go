package cache

/************************************************
Cache Master supports

m = Make()
    Initialize a cache master
m.Get(filename string)


*************************************************/


type CacheMaster {
    mu sync.Mutex
}


func (m *CacheMaster) Make {

}


