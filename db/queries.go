package db

import (
	"database/sql"

	"github.com/isaacmaddox/mc-project-bot/util"
)

var db *sql.DB

func Init_database() {
	db = util.Extract(sql.Open("sqlite3", "file:projects.db"))

	create_project_table := `
		DROP TABLE IF EXISTS example;
		CREATE TABLE IF NOT EXISTS project (
			id integer not null primary key,
			name string,
			description string
		);
	`

	util.Extract(db.Exec(create_project_table))

	create_resource_table := `
		CREATE TABLE IF NOT EXISTS resource (
			id integer not null primary key,
			name string,
			amount integer,
			goal integer,
			project_id integer,
			FOREIGN KEY (project_id) REFERENCES project(id)
		);
	`

	util.Extract(db.Exec(create_resource_table))
}

func (p *Project) Create(name, description string) {
	stmt := util.Extract(db.Prepare(`
		INSERT INTO project
		(name, description)
		VALUES (?, ?) RETURNING *;
	`))

	defer stmt.Close()

	util.ErrorCheck(
		stmt.QueryRow(name, description).Scan(&p.ID, &p.Name, &p.Description), "Error inserting project: %v",
	)
}

func (p *Project) Get(name string) bool {
	get_project := util.Extract(db.Prepare(`
		SELECT * FROM project
		WHERE name LIKE ?;
	`))

	defer get_project.Close()

	err := get_project.QueryRow(name).Scan(&p.ID, &p.Name, &p.Description)

	if err != nil {
		return false
	}

	get_resources := util.Extract(db.Prepare(`
		SELECT * FROM resource
		WHERE project_id = ?
	`))

	defer get_resources.Close()

	rows := util.Extract(get_resources.Query(p.ID))
	var throwaway int
	p.Resources = []*Resource{}

	for rows.Next() {
		var resource Resource
		util.ErrorCheck(rows.Scan(&resource.ID, &resource.Name, &resource.Amount, &resource.Goal, &throwaway), "bonk %v")
		p.Resources = append(p.Resources, &resource)
	}

	if err = rows.Err(); err != nil {
		return false
	}

	return true
}

func (p *Project) AddResource(name string, amount, goal int) *Resource {
	stmt := util.Extract(db.Prepare(`
		INSERT INTO resource
		(name, amount, goal, project_id)
		VALUES (?, ?, ?, ?) RETURNING *
	`))

	defer stmt.Close()

	var resource Resource

	var throwaway int

	util.ErrorCheck(stmt.QueryRow(name, amount, goal, p.ID).Scan(&resource.ID, &resource.Name, &resource.Amount, &resource.Goal, &throwaway), "Error inserting resource: %v")

	p.Resources = append(p.Resources, &resource)

	return &resource
}

func GetProjectNames() (result []string) {
	rows := util.Extract(db.Query(`SELECT name FROM project`))

	for rows.Next() {
		var name string
		util.ErrorCheck(rows.Scan(&name), "bonk %v")
		result = append(result, name)
	}

	return
}
