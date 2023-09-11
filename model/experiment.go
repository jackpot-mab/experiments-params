package model

type Experiment struct {
	ExperimentId string      `json:"experiment_id"`
	PolicyType   string      `json:"policy_type"`
	Arms         []Arm       `json:"arms"`
	Parameters   interface{} `json:"parameters"`
}

type Arm struct {
	Name string `json:"name"`
}

type EpsilonGreedyParams struct {
	Alpha float32
}
