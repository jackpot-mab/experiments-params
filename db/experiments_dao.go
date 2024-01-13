package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
	"jackpot-mab/experiments-params/model"
	"log"
	"time"
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

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Duration(3600) * time.Second)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		//log.Fatal(err)
	}

	return &ExperimentsDAOImpl{db: db}
}

func (e *ExperimentsDAOImpl) GetExperiment(experimentId string) model.Experiment {

	experimentRow, err := e.db.Query(`
						SELECT e.experiment_id, e.policy_type, e.parameters,
							   a.name,
							   mp.model_type, mp.input_features, mp.output_classes,
							   rdp.param_name, rdp.param_value 
						FROM experiment e 
						LEFT JOIN arm a on (a.experiment_id = e.experiment_id)
						LEFT JOIN model_params mp on (mp.experiment_id = e.experiment_id)
						LEFT JOIN reward_data_params rdp on (rdp.experiment_id = e.experiment_id and rdp.arm_id = a.arm_id)
						WHERE e.experiment_id = ?`, experimentId)

	if err != nil {
		log.Print(err)
		return model.Experiment{}
	}
	defer experimentRow.Close()

	experiment := &model.Experiment{}
	var modelParams model.MLModelParameters
	var arms []model.Arm
	rewardDataParamsByArm := make(map[string][]model.RewardDataParameter)

	for experimentRow.Next() {

		var arm model.Arm
		var rewardDataParams model.RewardDataParameter

		var paramName, paramVal sql.NullString

		var experimentParameters interface{}
		var inputFeaturesJson interface{}
		var outputClassesJson interface{}

		err := experimentRow.Scan(&experiment.ExperimentId, &experiment.PolicyType,
			&experimentParameters, &arm.Name, &modelParams.ModelType, &inputFeaturesJson, &outputClassesJson,
			&paramName, &paramVal)

		if paramName.Valid {
			rewardDataParams.Name = paramName.String
		}

		if paramVal.Valid {
			rewardDataParams.Value = paramVal.String
		}

		if err != nil {
			return model.Experiment{}
		}

		err = json.Unmarshal(inputFeaturesJson.([]byte), &modelParams.InputFeatures)
		err = json.Unmarshal(outputClassesJson.([]byte), &modelParams.OutputClasses)
		err = json.Unmarshal(experimentParameters.([]byte), &experiment.Parameters)

		arms = append(arms, arm)

		_, ok := rewardDataParamsByArm[arm.Name]

		if !ok {
			rewardDataParamsByArm[arm.Name] = []model.RewardDataParameter{{
				Name:  rewardDataParams.Name,
				Value: rewardDataParams.Value,
			}}
		} else {
			rewardDataParamsByArm[arm.Name] = append(rewardDataParamsByArm[arm.Name], model.RewardDataParameter{
				Name:  rewardDataParams.Name,
				Value: rewardDataParams.Value,
			})
		}

	}

	for _, arm := range arms {
		allParams, ok := rewardDataParamsByArm[arm.Name]
		if ok {
			arm.RewardDataParameters = allParams
		}
	}

	experiment.Arms = arms
	experiment.ModelParameters = modelParams

	return *experiment

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
