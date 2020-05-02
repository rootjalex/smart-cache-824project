package task

type Workload struct {
}

func (wkld *Workload) HasNextItemGroup() bool {
	return false
}

func (wkld *Workload) GetNextItemGroup() []string {
	return []string{}
}

func NewWorkload() Workload {
	return Workload{}
}

func NewMLWorkload() Workload {
	return Workload{}

}
