DROP TABLE IF EXISTS job_dep;
DROP TABLE IF EXISTS job;
DROP TABLE IF EXISTS run;
DROP TABLE IF EXISTS session;
DROP TABLE IF EXISTS appkey;
DROP TABLE IF EXISTS usr;


CREATE TABLE usr (
    id SERIAL PRIMARY KEY, 
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(1024) NOT NULL, -- format: METHOD$PAYLOAD. 
                                     -- Supported methods: 1$ => bcrypt (only one so far).
    is_admin BOOLEAN DEFAULT FALSE
);

-- appkeys are replacements for username/passwords that can be stored in files for 
-- automated authentication (scripted). They can also be revoked easier than a 
-- password can be changed
CREATE TABLE appkey (
    user_id INTEGER REFERENCES usr(id),
    access_key UUID, -- randomly generated
    secret_key VARCHAR(1024), -- randomly generated, used for the initial login request
                              -- after that the session-id is used
    created_ts TIMESTAMP,
    PRIMARY KEY(user_id, access_key)
);

-- sessions can be initiated from either an access/secretkey or username/password
CREATE TABLE session (
    user_id INTEGER REFERENCES usr(id),
    session_id UUID, -- randomly generated
    secret VARCHAR(1024), -- private key used for HMAC signing of each API request
    valid_until_ts TIMESTAMP, -- after this time, the session is no longer valid
    PRIMARY KEY(user_id, session_id)
);

-- a run is a group of jobs that should be tracked together. Jobs can be related (dependencies), 
--but they don't necessarily need to be.
CREATE TABLE run (
    id UUID PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES usr(id),
    name VARCHAR(1024),
    submit_ts TIMESTAMP,
    is_done BOOLEAN DEFAULT FALSE,
    data JSONB
);

-- an individual unit of work, store meta data in JSON blob
CREATE TABLE job (
    id UUID PRIMARY KEY,
    run_id UUID REFERENCES run(id),
    job_id VARCHAR(1024),
    name VARCHAR(1024),
    is_done BOOLEAN DEFAULT FALSE,
    is_error BOOLEAN DEFAULT FALSE,
    submit_ts TIMESTAMP, -- when submitted
    start_ts TIMESTAMP, -- when started (queue time is the time the last parent job finished)
    end_ts TIMESTAMP, -- when done
    retcode BOOLEAN, -- when retcode is set, the job is done
    data JSONB
);

-- job dependencies... DAC
CREATE TABLE job_dep (
    parent_id UUID REFERENCES job(id),
    child_id UUID REFERENCES job(id),
    PRIMARY KEY (parent_id, child_id)
);
