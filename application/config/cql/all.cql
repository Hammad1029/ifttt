CREATE KEYSPACE 
    IF NOT EXISTS ifttt 
        WITH replication = { 'class': 'NetworkTopologyStrategy', 'replication_factor': '1' };

use ifttt;

CREATE TABLE
    IF NOT EXISTS ifttt.apis (
        group TEXT,
        name TEXT,
        method TEXT,
        type TEXT,
        path TEXT,
        description TEXT,
        request TEXT,
        dumping TEXT,
        start_rules SET<TEXT>,
        rules TEXT,
        PRIMARY KEY (group,name)
    );

CREATE INDEX
    IF NOT EXISTS ON ifttt.apis(path);

CREATE TYPE
    IF NOT EXISTS ifttt.request_data (
        req_body  TEXT,
        store    TEXT,
        response TEXT,
        query_res TEXT,
        api_res TEXT
    );

CREATE TYPE
    IF NOT EXISTS ifttt.exec_log (
        log_user TEXT,
        log_type TEXT,
        log_data TEXT
    );

CREATE TABLE
    IF NOT EXISTS ifttt.logs (
        api_group TEXT,
        api_name TEXT,
        execution_order LIST<TEXT>,
        request_data FROZEN<request_data>,
        start TIMESTAMP,
        start_partition TIMESTAMP,
        end TIMESTAMP,
        time_taken INT,
        execution_logs LIST<FROZEN<exec_log>>,
        PRIMARY KEY ((api_group,start_partition),api_name,start)
    );

CREATE TABLE
    IF NOT EXISTS ifttt.auth_strategies (
        strategy TEXT,
        enabled BOOLEAN,
        secrets MAP<TEXT,TEXT>
    );

CREATE TABLE
    IF NOT EXISTS ifttt.configurations (
        configuration_json TEXT,
        is_active BOOLEAN,
        created_at TIMESTAMP,
        PRIMARY KEY (is_active)
    );