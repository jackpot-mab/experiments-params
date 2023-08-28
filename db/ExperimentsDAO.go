package db

import "jackpot-mab/experiments-params/model"

type ExperimentsDAO interface{}

type ExperimentsDAOImpl struct{}

func MakeExperimentsDAO() ExperimentsDAO {
	return nil
}

func (e *ExperimentsDAOImpl) GetExperiment() model.Experiment {
	return model.Experiment{
		ExperimentId: "1-EEE-2A",
		PolicyType:   "epsilon_greedy",
		Arms:         []model.Arm{{Name: "A"}, {Name: "B"}, {Name: "C"}},
		Parameters:   model.EpsilonGreedyParams{Alpha: 0.2},
	}
}

func (e *ExperimentsDAOImpl) AddExperiment() model.Experiment {
	return model.Experiment{
		ExperimentId: "1-EEE-2A",
		PolicyType:   "epsilon_greedy",
		Arms:         []model.Arm{{Name: "A"}, {Name: "B"}, {Name: "C"}},
		Parameters:   model.EpsilonGreedyParams{Alpha: 0.2},
	}
}

func (e *ExperimentsDAOImpl) UpdateExperiment() model.Experiment {
	return model.Experiment{
		ExperimentId: "1-EEE-2A",
		PolicyType:   "epsilon_greedy",
		Arms:         []model.Arm{{Name: "A"}, {Name: "B"}, {Name: "C"}},
		Parameters:   model.EpsilonGreedyParams{Alpha: 0.2},
	}
}
