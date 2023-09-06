package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
	"jackpot-mab/experiments-params/model"
	"log"
)

type ExperimentsDAO interface {
	GetExperiment(experimentId string) model.Experiment
	AddExperiment(create model.Experiment) model.Experiment
	UpdateExperiment(update model.Experiment) model.Experiment
	Close()
}

type ExperimentsDAOImpl struct {
	db *sql.DB
}

type ConnectionParams struct {
	User     string
	Password string
	Host     string
	Port     string
	DbName   string
}

const ConnectionFormat = "%s:%s@tcp(%s:%s)/%s"

func MakeExperimentsDAO(connection ConnectionParams) ExperimentsDAO {

	connectionString := fmt.Sprintf(ConnectionFormat, connection.User, connection.Password, connection.Host,
		connection.Port, connection.DbName)

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	return &ExperimentsDAOImpl{db: db}
}

func (e *ExperimentsDAOImpl) GetExperiment(experimentId string) model.Experiment {

	experimentRow, err := e.db.Query(
		"SELECT experiment_id, policy_type, parameters FROM experiment WHERE experiment_id = ?", experimentId)

	if err != nil {
		// TODO improve error handling.
		log.Print(err)
	}
	defer experimentRow.Close()

	// There shouldn't be more than one result.
	experimentRow.Next()
	var experiment model.Experiment
	err = experimentRow.Scan(&experiment.ExperimentId, &experiment.PolicyType, &experiment.Parameters)

	if err != nil {
		log.Print(err)
		return model.Experiment{}
	} else {

		experiment.Arms = e.getArms(experimentId)
	}

	return experiment
}

func (e *ExperimentsDAOImpl) getArms(experimentId string) []model.Arm {
	armsRows, err := e.db.Query(
		"SELECT name FROM arm WHERE experiment_id = ?", experimentId)
	defer armsRows.Close()

	var arms []model.Arm
	for armsRows.Next() {
		var arm model.Arm
		err = armsRows.Scan(&arm.Name)
		if err == nil {
			arms = append(arms, arm)
		}
	}
	return arms
}

func (e *ExperimentsDAOImpl) AddExperiment(create model.Experiment) model.Experiment {
	return model.Experiment{
		ExperimentId: "1-EEE-2A",
		PolicyType:   "epsilon_greedy",
		Arms:         []model.Arm{{Name: "A"}, {Name: "B"}, {Name: "C"}},
		Parameters:   model.EpsilonGreedyParams{Alpha: 0.2},
	}
}

func (e *ExperimentsDAOImpl) UpdateExperiment(update model.Experiment) model.Experiment {
	return model.Experiment{
		ExperimentId: "1-EEE-2A",
		PolicyType:   "epsilon_greedy",
		Arms:         []model.Arm{{Name: "A"}, {Name: "B"}, {Name: "C"}},
		Parameters:   model.EpsilonGreedyParams{Alpha: 0.2},
	}
}

func (e *ExperimentsDAOImpl) Close() {
	err := e.db.Close()
	if err != nil {
		log.Print("Error closing db.")
	}
}
