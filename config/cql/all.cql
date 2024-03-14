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

CREATE TABLE
    IF NOT EXISTS ifttt.tables (
        internal_name TEXT,
        name TEXT,
        description TEXT,
        partition_keys SET<TEXT>,
        clustering_keys SET<TEXT>,
        all_columns SET<TEXT>,
        mappings MAP<TEXT, TEXT>,
        indexes map<text, frozen<indexes>>,
        PRIMARY KEY (internal_name, name, description)
    );

CREATE TABLE
    IF NOT EXISTS ifttt.newTableName (
        created_at timestamp,
        int1 int,
        int2 int,
        int13 int,
        PRIMARY KEY ((int1),(int2))
    );

time_series_partition:
1-59s
1-59m
1-23h
1-27d
1-11M
1y

CREATE TYPE 
    IF NOT EXISTS action (
        type TEXT,
        data TEXT
);

CREATE TYPE 
    IF NOT EXISTS rule (
        id UUID,
        op1 TEXT,
        opnd TEXT, 
        op2 TEXT,
        then SET<FROZEN<ACTION>>,
        else SET<FROZEN<ACTION>>
);

CREATE TABLE
    IF NOT EXISTS ifttt.apis (
        api_group TEXT,
        api_name TEXT,
        api_description TEXT,
        api_path TEXT,
        rules SET<FROZEN<RULE>>,
        PRIMARY KEY (api_group,api_name,api_description)
    );

CREATE INDEX ON ifttt.apis(api_path);