package cache

import (
	"../markov"
)

// master -> cache (Request param)
type ModelParamArgs struct {

}

type ModelParamReply struct {
	Chain		markov.MarkovChain
}

// master -> cache (Communicate update)
type ModelParamUpdateArgs struct {
	Chain		markov.MarkovChain // update
}

type ModelParamUpdateReply struct {
	Success	bool
}


type GetCacheStateArgs struct {
}

type GetCacheStateReply struct {
    State    *markov.MarkovChain
}

type UpdateCacheArgs struct {
    State    *markov.MarkovChain
}

type UpdateCacheReply struct {
    Success  bool
}
