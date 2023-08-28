package model

type Experiment struct {
	ExperimentId string
	PolicyType   string
	Arms         []Arm
	Parameters   interface{}
}

type Arm struct {
	Name string
}

type EpsilonGreedyParams struct {
	Alpha float32
}
