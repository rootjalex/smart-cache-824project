package task

type MLParams struct {
	MinFileSleep	int
	MaxFileSleep	int
	MinBatchSleep	int
	MaxBatchSleep	int
	NBatches        int
	BatchLength     int
	NIterations     int
	Name            string
}

type WebParams struct {
	MinFileSleep		int
	MaxFileSleep		int
	MinBatchSleep		int
	MaxBatchSleep		int
	NBatches			int
	NPatterns			int
	MinPatternLength	int
	MaxPatternLength	int
	MaxFileCount		int
	BatchLength         int
	Name				string
}


type RandomParams struct {
	MinFileSleep	int
	MaxFileSleep	int
	MinBatchSleep	int
	MaxBatchSleep	int
	NBatches		int
	MinBatchLength	int
	MaxBatchLength	int
	MaxFileCount	int
	BatchLength     int
	NIterations     int
	Name            string
}

