package cache

import (

)

// master -> cache (Request param)
type ModelParamArgs struct {

}

type ModelParamReply struct {
	Chain 		MarkovChain
}

// master -> cache (Communicate update)
type ModelParamUpdateArgs struct {
	Chain 		MarkovChain // update
}

type ModelParamUpdateReply struct {
	Success 	bool
}

// client -> cache (Request a file)
type RequestFileArgs {
	Filename 	string
}

type RequestFileReply {
	File 		*os.File 
	Hit 		bool 
}