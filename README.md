# PHOTOSHARING WEBSITE

docker compose up/down

```connect to postgres docker```
docker exec -it website-db-1 psql -U postgres -d website 

```to create user table postgres```
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    age INT,
    first_name TEXT,
    last_name TEXT,
    email TEXT UNIQUE NOT NULL,
)

```insert into table```
INSERT INTO users VALUES (1, 22, 'John','Smith','john@smith.com');
```or this way```
INSERT INTO users(id, email, password_hash)
VALUES (2, 'gav@gmail.com', 'gfdsfgsfdfsdfdsf232323');

```udpate record```
UPDATE users SET first_name = 'gavin' WHERE id = 1;

```delete record```
DELETE FROM users WHERE id=1;

```alter existing table```
ALTER TABLE sessions ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id);