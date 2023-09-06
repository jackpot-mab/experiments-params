CREATE DATABASE IF NOT EXISTS experiment_db;
USE experiment_db;

CREATE TABLE IF NOT EXISTS experiment (
    experiment_id VARCHAR(10) PRIMARY KEY,
    policy_type VARCHAR(50) NOT NULL,
    parameters JSON
);

CREATE TABLE IF NOT EXISTS arm (
    arm_id INT AUTO_INCREMENT PRIMARY KEY,
    experiment_id VARCHAR(10),
    name VARCHAR(50) NOT NULL,
    FOREIGN KEY (experiment_id) REFERENCES experiment(experiment_id)
);