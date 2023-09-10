package db

import (
	"database/sql"
	"encoding/json"
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
		return model.Experiment{}
	}
	defer experimentRow.Close()

	// There shouldn't be more than one result.
	for experimentRow.Next() {
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

	return model.Experiment{}
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
	stmt, err := e.db.Prepare(
		"INSERT INTO experiment (experiment_id, policy_type, parameters) VALUES (?, ?, ?)")
	if err != nil {
		log.Print(err)
		return model.Experiment{}
	}
	defer stmt.Close()

	jsonParams, _ := json.Marshal(create.Parameters)
	_, err = stmt.Exec(create.ExperimentId, create.PolicyType, jsonParams)
	if err != nil {
		log.Print(err)
		return model.Experiment{}
	}

	e.addArms(create.Arms, create.ExperimentId)

	return create
}

func (e *ExperimentsDAOImpl) addArms(arms []model.Arm, experimentId string) {
	for _, a := range arms {
		stmt, err := e.db.Prepare(
			"INSERT INTO arm (experiment_id, name) VALUES (?, ?)")

		if err != nil {
			log.Print(err)
			return
		}

		_, err = stmt.Exec(experimentId, a.Name)

		if err != nil {
			log.Print(err)
			return
		}
	}
}

func (e *ExperimentsDAOImpl) UpdateExperiment(update model.Experiment) model.Experiment {
	// TODO
	return model.Experiment{}
}

func (e *ExperimentsDAOImpl) Close() {
	err := e.db.Close()
	if err != nil {
		log.Print("Error closing db.")
	}
}
