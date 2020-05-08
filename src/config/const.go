package config

import (
	"time"
)

const CACHE_SIZE = 50
const SYNC_MS = 100
const PREFETCH_SIZE = 10

//const SEED = time.Now().UnixNano()
const SEED = 1

type CacheType int

const (
	LRU            CacheType = 0
	MarkovPrefetch CacheType = 1
	MarkovEviction CacheType = 2
	MarkovBoth     CacheType = 3
)

type DataType string

// DATA_DEFAULT must be initialized to an empty instance of the above default
const DATA_DEFAULT = ""

const TIME_MULTIPLIER = 100
const DATA_FETCH_TIME = time.Duration(1*TIME_MULTIPLIER) * time.Millisecond

const CLIENT_COMPUTATION_TIME = time.Duration(1*TIME_MULTIPLIER) * time.Millisecond

// latency stuff?

// web pattern configs
const MIN_PATTERN_LENGTH = 3
const MAX_PATTERN_LENGTH = 6
const MIN_PATTERN_WAIT = 10  // ms
const MAX_PATTERN_WAIT = 100 // ms
const NUM_PATTERNS_SMALL = NFILES_SMALL / MAX_PATTERN_LENGTH / 3
const NUM_PATTERNS_MED = NFILES_MED / MAX_PATTERN_LENGTH / 10
const PATTERN_REPLICATION_SMALL = NFILES_SMALL / ((MAX_PATTERN_LENGTH + MIN_PATTERN_LENGTH) / 2)
const PATTERN_REPLICATION_MED = NFILES_MED / ((MAX_PATTERN_LENGTH + MIN_PATTERN_LENGTH) / 2)

// benchmarking constants -- SMALL
const NFILES_SMALL = 200
const BATCH_SMALL = 16
const ITERS_SMALL = 10
const NCLIENTS_SMALL = 5
const NCACHES_SMALL = 2
const RFACTOR_SMALL = 1

// benchmarking constants -- MEDIUM
const NFILES_MED = 1000
const BATCH_MED = 32
const ITERS_MED = 20
const NCLIENTS_MED = 10
const NCACHES_MED = 10
const RFACTOR_MED = 2
