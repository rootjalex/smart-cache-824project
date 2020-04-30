package cache

import (
	"./markov"
)

// master -> cache (Request param)
type ModelParamArgs struct {
}

type ModelParamReply struct {
	Chain markov.MarkovChain
}

// master -> cache (Communicate update)
type ModelParamUpdateArgs struct {
	Chain markov.MarkovChain // update
}

type ModelParamUpdateReply struct {
	Success bool
}