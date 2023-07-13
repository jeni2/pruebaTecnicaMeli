#drop database databaseClasification ;
CREATE database databaseClasification ;

CREATE TABLE databaseClasification.dbschemas (
    schema_id INT AUTO_INCREMENT PRIMARY KEY,
    last_scan TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    db_id INT,
    schema_name VARCHAR(255)
);

CREATE TABLE databaseClasification.dbtables (
    table_id INT AUTO_INCREMENT PRIMARY KEY,
    schema_id INT,
    table_name VARCHAR(255),
    FOREIGN KEY (schema_id) REFERENCES dbSchemas (schema_id)
     ON DELETE CASCADE
);

CREATE TABLE databaseClasification.dbColumns (
    column_id INT AUTO_INCREMENT PRIMARY KEY,
    table_id INT,
    column_name VARCHAR(255),
    data_type VARCHAR(255),
    classification VARCHAR(255),
    FOREIGN KEY (table_id) REFERENCES dbTables (table_id)
     ON DELETE CASCADE
);
SELECT * from databaseClasification.dbschemas;

---------------------------------------------------------------
#drop database databaseCredentials ;
CREATE database databaseCredentials ;
CREATE TABLE databaseCredentials.credentials (
    id INT AUTO_INCREMENT PRIMARY KEY,
    dbhost VARCHAR(255),
    dbport INT,
    dbusername VARCHAR(255),
    dbpassword VARCHAR(255)
);

----------------------------------------------------------
#drop database prueba  ;
CREATE database prueba  ;

CREATE TABLE IF NOT EXISTS prueba.users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255),
    useremail VARCHAR(255),
    credit_card_number VARCHAR(16),
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

--------------------------------------------------------
#drop database privatedata  ;
CREATE database privatedata  ;

CREATE TABLE IF NOT EXISTS privatedata.privateword (
    word_id INT AUTO_INCREMENT PRIMARY KEY,
    word VARCHAR(255)
);
INSERT INTO privatedata.privateword (word) VALUES
    ('PASSWORD'),
    ('IP_ADDRESS'),
    ('LAST_NAME'),
    ('FIRST_NAME'),
    ('CREDIT_CARD_NUMBER'),
    ('USERNAME'),
    ('EMAIL_ADDRESS');
    
#delete from privatedata.privateword where word="USER"


