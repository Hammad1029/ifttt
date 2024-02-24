CREATE KEYSPACE 
    IF NOT EXISTS ifttt 
        WITH replication = { 'class': 'NetworkTopologyStrategy', 'replication_factor': '3' };

CREATE TABLE
    IF NOT EXISTS ifttt.tables (
        internal_name TEXT,
        name TEXT,
        description TEXT,
        partition_keys SET<TEXT>,
        clustering_keys SET<TEXT>,
        all_columns SET<TEXT>,
        mappings MAP<TEXT, TEXT>,
        PRIMARY KEY (internal_name, name, description)
    );

CREATE TABLE
    IF NOT EXISTS ifttt.newTableName (
        int1 int,
        int2 int,
        int13 int,
        PRIMARY KEY ((int1),(int2))
    );