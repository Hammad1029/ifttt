CREATE TABLE
    IF NOT EXISTS apis (
        id serial primary key,
        name varchar(50) not null,
        path varchar(50) not null unique,
        method varchar(10) not null,
        description text default '',
        request jsonb not null default '{}'::jsonb
        preconfig jsonb not null default '{}'::jsonb
    );

CREATE TABLE
    class (id serial primary key, name varchar(50) not null);

CREATE TABLE
    IF NOT EXISTS trigger_flows (
        id serial primary key,
        name varchar(50) not null,
        description text default '',
        class_id int references class (id) not null
    ); 

CREATE TABLE
    IF NOT EXISTS api_trigger_flows (
        api_id int references apis (id) not null,
        flow_id int references trigger_flows (id) not null
    );

CREATE TABLE
    IF NOT EXISTS rules (
        id serial primary key,
        name varchar(50) not null,
        description text default '',
        conditions jsonb not null default '{"group":true,"conditionType":"and","conditions":[]}'::jsonb,
        "then" jsonb not null default '[]'::jsonb,
        "else" jsonb not null default '[]'::jsonb
    );

CREATE TABLE
    IF NOT EXISTS trigger_start_rules (
        flow_id int references trigger_flows (id) not null,
        rule_id int references rules (id) not null
    );

CREATE TABLE
    IF NOT EXISTS trigger_all_rules (
        flow_id int references trigger_flows (id) not null,
        rule_id int references rules (id) not null
    );