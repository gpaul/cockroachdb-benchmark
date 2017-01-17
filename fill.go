package main

import (
	"database/sql"
	"fmt"
	"math/rand"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=26257 dbname=bench user=root sslmode=disable")
	if err != nil {
		panic(err)
	}
	structure := `
DROP TABLE IF EXISTS users;
CREATE TABLE users (id INT PRIMARY KEY, name STRING);

DROP TABLE IF EXISTS groups;
CREATE TABLE groups (id INT PRIMARY KEY, name STRING);

DROP TABLE IF EXISTS permissions;
CREATE TABLE permissions (id INT PRIMARY KEY, rid STRING);

DROP TABLE IF EXISTS groups_permissions;
CREATE TABLE groups_permissions (id INT PRIMARY KEY, group_id INT, permission_id INT);
CREATE INDEX ON groups_permissions (group_id, permission_id);
CREATE INDEX ON groups_permissions (permission_id, group_id);

DROP TABLE IF EXISTS users_permissions;
CREATE TABLE users_permissions (id INT PRIMARY KEY, user_id INT, permission_id INT);
CREATE INDEX ON users_permissions (user_id, permission_id);
CREATE INDEX ON users_permissions (permission_id, user_id);

DROP TABLE IF EXISTS users_groups;
CREATE TABLE users_groups (id INT PRIMARY KEY, user_id INT, group_id INT);
CREATE INDEX ON users_groups (user_id, group_id);
CREATE INDEX ON users_groups (group_id, user_id);
`
	_, err = db.Exec(structure)
	if err != nil {
		panic(err)
	}

	numberOfUsers := 1000
	numberOfGroups := 1000
	usersPerGroup := 100
	permissionsPerUser := 100
	permissionsPerGroup := 100

	for id := 0; id < numberOfUsers; id++ {
		fmt.Println("User ", id)
		name := fmt.Sprintf("user-%d", id)
		_, err := db.Exec("INSERT INTO users VALUES ($1, $2);", id, name)
		if err != nil {
			panic(err)
		}
	}

	mid := 0
	for id := 0; id < numberOfGroups; id++ {
		fmt.Println("Group ", id)
		name := fmt.Sprintf("group-%d", id)
		_, err := db.Exec("INSERT INTO groups VALUES ($1, $2);", id, name)
		if err != nil {
			panic(err)
		}
		choices := rand.Perm(numberOfUsers)[:usersPerGroup]
		for _, uu := range choices {
			_, err := db.Exec("INSERT INTO users_groups VALUES ($1, $2, $3);", mid, uu, id)
			if err != nil {
				panic(err)
			}
			mid += 1
		}
	}

	rid := 0
	for id := 0; id < permissionsPerUser; id++ {
		fmt.Println("User Permissions ", id)
		name := fmt.Sprintf("rid-user-%d", rid)
		_, err := db.Exec("INSERT INTO permissions VALUES ($1, $2);", id, name)
		if err != nil {
			panic(err)
		}
		rid += 1
	}

	for id := 0; id < permissionsPerGroup; id++ {
		fmt.Println("Group Permissions ", id)
		name := fmt.Sprintf("rid-group-%d", id)
		_, err := db.Exec("INSERT INTO permissions VALUES ($1, $2);", rid, name)
		if err != nil {
			panic(err)
		}
		rid += 1
	}

	id := 0
	for uid := 0; uid < numberOfUsers; uid++ {
		fmt.Println("Permissions for user ", uid)
		for rid := 0; rid < permissionsPerUser; rid++ {
			_, err := db.Exec("INSERT INTO users_permissions VALUES ($1, $2, $3);", id, uid, rid)
			if err != nil {
				panic(err)
			}
			id += 1
		}
	}

	id = 0
	for gid := 0; gid < numberOfGroups; gid++ {
		fmt.Println("Permissions for group ", gid)
		for rid := permissionsPerUser; rid < permissionsPerUser+permissionsPerGroup; rid++ {
			_, err := db.Exec("INSERT INTO groups_permissions VALUES ($1, $2, $3);", id, gid, rid)
			if err != nil {
				panic(err)
			}
			id += 1
		}
	}
}
