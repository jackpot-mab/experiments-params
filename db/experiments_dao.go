package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
	"jackpot-mab/experiments-params/model"
	"log"
)

type ExperimentsDAO interface {
	GetAllExperiments() []model.Experiment
	GetExperiment(experimentId string) model.Experiment
	AddExperiment(create model.Experiment) model.Experiment
	UpdateExperiment(update model.Experiment) model.Experiment
	AddOrUpdateRewardParameter(
		update model.RewardDataParameterUpsert) error
	DeleteExperiment(experimentId string) error
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
		var parameters interface{}
		err = experimentRow.Scan(&experiment.ExperimentId, &experiment.PolicyType, &parameters)

		if err != nil {
			log.Print(err)
			return model.Experiment{}
		} else {
			experiment.ModelParameters = e.getModelParameters(experimentId)
			experiment.Arms = e.getArms(experimentId)
		}

		jsonParameters := string(parameters.([]byte))

		err = json.Unmarshal([]byte(jsonParameters), &experiment.Parameters)

		return experiment
	}

	return model.Experiment{}
}

func (e *ExperimentsDAOImpl) getModelParameters(experimentId string) model.MLModelParameters {
	modelParams, _ := e.db.Query(
		"SELECT model_type, input_features, output_classes "+
			"FROM model_params WHERE experiment_id = ?", experimentId)
	defer modelParams.Close()

	var mlModelParams model.MLModelParameters
	for modelParams.Next() {
		var inputFeaturesJson interface{}
		var outputClassesJson interface{}
		err := modelParams.Scan(&mlModelParams.ModelType, &inputFeaturesJson, &outputClassesJson)
		err = json.Unmarshal(inputFeaturesJson.([]byte), &mlModelParams.InputFeatures)
		err = json.Unmarshal(outputClassesJson.([]byte), &mlModelParams.OutputClasses)
		if err != nil {
			return model.MLModelParameters{}
		}
	}

	return mlModelParams

}

func (e *ExperimentsDAOImpl) getArms(experimentId string) []model.Arm {
	armsRows, err := e.db.Query(
		"SELECT arm_id, name FROM arm WHERE experiment_id = ?", experimentId)
	defer armsRows.Close()

	var arms []model.Arm
	for armsRows.Next() {
		var arm model.Arm
		var armId int
		err = armsRows.Scan(&armId, &arm.Name)
		if err == nil {
			arm.RewardDataParameters = e.getRewardDataParams(experimentId, armId)
			arms = append(arms, arm)
		}
	}
	return arms
}

func (e *ExperimentsDAOImpl) getRewardDataParams(experimentId string, armId int) []model.RewardDataParameter {
	paramsRows, err := e.db.Query(
		"SELECT param_name, param_value FROM reward_data_params WHERE experiment_id = ? and arm_id = ?", experimentId, armId)
	defer paramsRows.Close()
	var rewardDataParams []model.RewardDataParameter
	for paramsRows.Next() {
		var paramName, paramValue string
		err = paramsRows.Scan(&paramName, &paramValue)
		if err == nil {
			rewardDataParams = append(rewardDataParams, model.RewardDataParameter{
				Name:  paramName,
				Value: paramValue,
			})
		}
	}
	return rewardDataParams
}

func (e *ExperimentsDAOImpl) getArmId(experimentId string, armName string) int {
	armId, err := e.db.Query(
		"SELECT arm_id FROM arm WHERE experiment_id = ? AND name = ?", experimentId, armName)
	defer armId.Close()

	var result int
	armId.Next()

	err = armId.Scan(&result)
	if err != nil {
		return -1
	}

	return result
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
	e.addModel(create)

	return create
}

func (e *ExperimentsDAOImpl) addModel(experiment model.Experiment) {

	stmt, err := e.db.Prepare(
		"INSERT INTO model_params (experiment_id, model_type, input_features, output_classes) VALUES (?, ?, ?, ?)")
	defer stmt.Close()

	if err != nil {
		return
	}

	marshallInputFeatures, _ := json.Marshal(experiment.ModelParameters.InputFeatures)
	marshalOutputClasses, _ := json.Marshal(experiment.ModelParameters.OutputClasses)
	_, err = stmt.Exec(experiment.ExperimentId,
		experiment.ModelParameters.ModelType, marshallInputFeatures,
		marshalOutputClasses)

	if err != nil {
		return
	}
}

func (e *ExperimentsDAOImpl) addArms(arms []model.Arm, experimentId string) {
	for _, a := range arms {
		insertExperimentArmsStmt, err := e.db.Prepare(
			"INSERT INTO arm (experiment_id, name) VALUES (?, ?)")

		if err != nil {
			log.Print(err)
			return
		}

		_, err = insertExperimentArmsStmt.Exec(experimentId, a.Name)

		if err != nil {
			log.Print(err)
			return
		}
	}
}

func (e *ExperimentsDAOImpl) GetAllExperiments() []model.Experiment {
	experimentIds, err := e.db.Query(
		"SELECT experiment_id FROM experiment")
	defer experimentIds.Close()

	var experiments []model.Experiment

	for experimentIds.Next() {
		var id string
		err = experimentIds.Scan(&id)
		if err == nil {
			exp := e.GetExperiment(id)
			if exp.ExperimentId != "" {
				experiments = append(experiments, e.GetExperiment(id))
			}
		}
	}
	return experiments

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

func (e *ExperimentsDAOImpl) AddOrUpdateRewardParameter(
	update model.RewardDataParameterUpsert) error {

	// Prepare the INSERT statement
	insertStmt, err := e.db.Prepare("REPLACE INTO reward_data_params (experiment_id, arm_id, param_name, param_value) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer insertStmt.Close()
	armId := e.getArmId(update.ExperimentId, update.ArmName)

	if armId == -1 {
		return errors.New("error fetching arm")
	}

	// Execute the INSERT statement
	_, err = insertStmt.Exec(update.ExperimentId, armId, update.Name, update.Value)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil

}

func (e *ExperimentsDAOImpl) DeleteExperiment(experimentId string) error {
	// DELETE PARAMETERS
	// DELETE ARMS
	// DELETE EXPERIMENT
	return nil
}
