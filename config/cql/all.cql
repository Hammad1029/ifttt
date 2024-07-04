CREATE KEYSPACE 
    IF NOT EXISTS ifttt 
        WITH replication = { 'class': 'NetworkTopologyStrategy', 'replication_factor': '1' };

use ifttt;

CREATE TYPE 
    IF NOT EXISTS ifttt.indexes (
        local BOOLEAN,
        index_name TEXT,
        table_name TEXT,
        columns SET<TEXT>
);

CREATE TYPE
    IF NOT EXISTS ifttt.table_column (
        name TEXT,
        data_type TEXT,
        partition_key BOOLEAN,
        clustering_key BOOLEAN
);

CREATE TABLE
    IF NOT EXISTS ifttt.tables (
        internal_name TEXT,
        name TEXT,
        description TEXT,
        columns SET<FROZEN<table_column>>,
        indexes map<text, frozen<indexes>>,
        PRIMARY KEY (internal_name, name, description)
    );


CREATE TABLE
    IF NOT EXISTS ifttt.apis (
        api_group TEXT,
        api_name TEXT,
        api_description TEXT,
        api_path TEXT,
        api_request TEXT,
        start_rules SET<TEXT>,
        rules TEXT,
        queries TEXT,
        PRIMARY KEY (api_group,api_name,api_description)
    );

CREATE INDEX
    IF NOT EXISTS ON ifttt.apis(api_path);

CREATE TYPE
    IF NOT EXISTS ifttt.request_data (
        req_body  MAP<TEXT,TEXT>,
        store    MAP<TEXT,TEXT>,
        response MAP<TEXT,TEXT>,
        query_res MAP<TEXT,FROZEN<SET<MAP<TEXT,TEXT>>>>,
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
        execution_order LIST<INT>,
        request_data FROZEN<request_data>,
        start TIMESTAMP,
        start_partition TIMESTAMP,
        end TIMESTAMP,
        time_taken INT,
        execution_logs LIST<FROZEN<exec_log>>,
        PRIMARY KEY ((api_group,start_partition),api_name,start)
    );