CREATE DATABASE IF NOT EXISTS experiment_db;
USE experiment_db;

CREATE TABLE IF NOT EXISTS experiment (
    experiment_id VARCHAR(250) PRIMARY KEY,
    policy_type VARCHAR(50) NOT NULL,
    parameters JSON
);

CREATE TABLE IF NOT EXISTS arm (
    arm_id INT AUTO_INCREMENT PRIMARY KEY,
    experiment_id VARCHAR(250),
    name VARCHAR(250) NOT NULL,
    FOREIGN KEY (experiment_id) REFERENCES experiment(experiment_id)
);

CREATE TABLE IF NOT EXISTS reward_data_params (
    experiment_id VARCHAR(250) ,
    arm_id INT,
    param_name VARCHAR(250),
    param_value VARCHAR(500),
    PRIMARY KEY (experiment_id, arm_id, param_name),
    FOREIGN KEY (experiment_id) REFERENCES experiment(experiment_id),
    FOREIGN KEY (arm_id) REFERENCES arm(arm_id)
);

CREATE TABLE IF NOT EXISTS model_params (
    experiment_id VARCHAR(250),
    model_type VARCHAR(250),
    input_features JSON,
    output_classes JSON,
    PRIMARY KEY (experiment_id),
    FOREIGN KEY (experiment_id) REFERENCES experiment(experiment_id)
);