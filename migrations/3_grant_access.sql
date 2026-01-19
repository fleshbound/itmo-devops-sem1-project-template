\c "project-sem-1";

CREATE USER validator WITH PASSWORD 'val1dat0r';
GRANT ALL PRIVILEGES ON DATABASE "project-sem-1" TO validator;
GRANT ALL PRIVILEGES ON TABLE prices TO validator;