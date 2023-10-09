package model

type Experiment struct {
	ExperimentId    string            `json:"experiment_id"`
	PolicyType      string            `json:"policy_type"`
	Arms            []Arm             `json:"arms"`
	Parameters      interface{}       `json:"parameters"`
	ModelParameters MLModelParameters `json:"model_parameters"`
}

type Arm struct {
	Name                 string                `json:"name"`
	RewardDataParameters []RewardDataParameter `json:"reward_data_parameters"`
}

type EpsilonGreedyParams struct {
	Alpha float32
}

type MLModelParameters struct {
	ModelType     string   `json:"model_type"`
	InputFeatures []string `json:"input_features"`
	OutputClasses []string `json:"output_classes"`
}

type RewardDataParameterUpsert struct {
	ExperimentId string      `json:"experiment_id"`
	ArmName      string      `json:"arm"`
	Name         string      `json:"name"`
	Value        interface{} `json:"value"`
}

type RewardDataParameter struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
