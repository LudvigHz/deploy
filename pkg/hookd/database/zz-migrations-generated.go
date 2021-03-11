// file generated by go generate

package database

var migrations = []string{
"-- Run the entire migration as an atomic operation.\nSTART TRANSACTION ISOLATION LEVEL SERIALIZABLE READ WRITE;\n\n-- Table apikey holds teams' deploy API keys.\n-- A team can have many API keys, with each key having its own expiry time.\nCREATE TABLE apikey\n(\n    \"key\"           varchar primary key      not null,\n    \"team\"          varchar                  not null,\n    \"team_azure_id\" varchar                  not null,\n    \"created\"       timestamp with time zone not null,\n    \"expires\"       timestamp with time zone null\n);\n\nCREATE INDEX apikey_team_index ON apikey (team);\nCREATE INDEX apikey_team_azure_id_index ON apikey (team_azure_id);\n\n-- Each row in the deployment table represents a single deployment request.\nCREATE TABLE deployment\n(\n    \"id\"                varchar primary key      not null,\n    \"team\"              varchar                  not null,\n    \"created\"           timestamp with time zone not null,\n    \"github_id\"         int unique               null,\n    \"github_repository\" varchar                  null\n);\n\n-- A row is recorded in deployment_status for each state change in a deployment.\nCREATE TABLE deployment_status\n(\n    \"id\"            varchar primary key                not null,\n    \"deployment_id\" varchar references deployment (id) not null,\n    \"status\"        varchar                            not null,\n    \"message\"       varchar                            not null,\n    \"github_id\"     int                                null,\n    \"created\"       timestamp with time zone           not null\n);\n\n-- Database migration\nCREATE TABLE migrations\n(\n    \"version\" int primary key          not null,\n    \"created\" timestamp with time zone not null\n);\n\n-- Mark this database migration as completed.\nINSERT INTO migrations (version, created)\nVALUES (1, now());\nCOMMIT;\n",
"-- Run the entire migration as an atomic operation.\nSTART TRANSACTION ISOLATION LEVEL SERIALIZABLE READ WRITE;\n\n-- Table team_repositories holds information about which repository can deploy to which team's resources.\n-- This supports the use of the legacy version in pkg/server/github_handler.go.\nCREATE TABLE team_repositories\n(\n    \"team\"       varchar not null,\n    \"repository\" varchar not null\n);\n\nCREATE INDEX team_repositories_team ON team_repositories (team);\nCREATE INDEX team_repositories_repository ON team_repositories (repository);\nCREATE UNIQUE INDEX team_repositories_unique ON team_repositories (team, repository);\n\n-- Mark this database migration as completed.\nINSERT INTO migrations (version, created)\nVALUES (2, now());\nCOMMIT;\n",
"-- Run the entire migration as an atomic operation.\nSTART TRANSACTION ISOLATION LEVEL SERIALIZABLE READ WRITE;\n\n-- This field has never been used and we don't intend to use it anyway.\nALTER TABLE deployment_status\n    DROP github_id;\n\n-- Mark this database migration as completed.\nINSERT INTO migrations (version, created)\nVALUES (3, now());\nCOMMIT;\n",
"-- Run the entire migration as an atomic operation.\nSTART TRANSACTION ISOLATION LEVEL SERIALIZABLE READ WRITE;\n\n-- These indices are neccessary to run the deployment frontend dashboard in a timely manner.\nCREATE INDEX deployment_created ON deployment (created);\nCREATE INDEX deployment_status_deployment_id ON deployment_status (deployment_id);\n\n-- Mark this database migration as completed.\nINSERT INTO migrations (version, created)\nVALUES (4, now());\nCOMMIT;\n",
"-- Run the entire migration as an atomic operation.\nSTART TRANSACTION ISOLATION LEVEL SERIALIZABLE READ WRITE;\n\n-- Table deployment_resource holds the identifiers of resources touched by this deployment.\nCREATE TABLE deployment_resource\n(\n    \"id\"            varchar primary key                not null,\n    \"deployment_id\" varchar references deployment (id) not null,\n    \"index\"         int                                not null,\n    \"group\"         varchar                            not null,\n    \"version\"       varchar                            not null,\n    \"kind\"          varchar                            not null,\n    \"name\"          varchar                            not null,\n    \"namespace\"     varchar                            not null\n);\n\nCREATE INDEX deployment_resource_deployment_id ON deployment_resource (deployment_id);\n\n-- Mark this database migration as completed.\nINSERT INTO migrations (version, created)\nVALUES (5, now());\nCOMMIT;\n",
"-- Run the entire migration as an atomic operation.\nSTART TRANSACTION ISOLATION LEVEL SERIALIZABLE READ WRITE;\n\n-- Add cluster field to deployment table.\nALTER TABLE deployment\nADD COLUMN \"cluster\" VARCHAR NULL;\n\n-- Mark this database migration as completed.\nINSERT INTO migrations (version, created)\nVALUES (6, now());\nCOMMIT;\n",
"-- Run the entire migration as an atomic operation.\nSTART TRANSACTION ISOLATION LEVEL SERIALIZABLE READ WRITE;\n\n-- Add cluster field to deployment table.\nALTER TABLE deployment\nADD COLUMN \"state\" VARCHAR NULL;\n\n-- Enable fast lookups on cluster and state\nCREATE INDEX deployment_state ON deployment (state);\nCREATE INDEX deployment_cluster ON deployment (cluster);\n\n-- Mark this database migration as completed.\nINSERT INTO migrations (version, created)\nVALUES (7, now());\nCOMMIT;\n",

}
