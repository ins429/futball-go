martini, go get github.com/codegangsta/martini
martini render, go get github.com/codegangsta/martini-contrib/render
postgres, go get github.com/lib/pq

run postgresql
pg_ctl -D /usr/local/var/postgres -l /usr/local/var/postgres/server.log start

CREATE TABLE users (
  id SERIAL,
  email TEXT,
  fb_id TEXT,
  firstname TEXT,
  lastname TEXT,
  encrypted_password TEXT,
  created_at DATE,
  updated_at DATE
);

CREATE TABLE cards (
  id SERIAL,
  name TEXT,
  created_at DATE,
  updated_at DATE
);

ALTER TABLE users ADD id SERIAL;

curl -d "email=ins429@gmail.com&password=pass" "localhost:8080/signup"