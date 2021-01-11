//args: -Egoone
package testdata

import (
	"database/sql"
	"fmt"
	"log"
)

type Person struct {
	Name  string
	JobID int
}

type Job struct {
	JobID int
	Name  string
}

func ForStmt() {

	cnn, _ := sql.Open("mysql", "user:password@tcp(host:port)/dbname")

	rows, _ := cnn.Query("SELECT name, job_id FROM persons")

	defer rows.Close()

	for rows.Next() {
		var person Person
		if err := rows.Scan(&person.Name, &person.JobID); err != nil {
			log.Fatal(err)
		}

		var job Job
		err := cnn.QueryRow("SELECT job_id, name FROM Jobs WHERE job_id = ?", person.JobID).Scan(&job.JobID, &job.Name)//ERROR"this query is called in a loop"
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(person.Name, job.Name)
	}

}

func rangeStmt() {

	cnn, _ := sql.Open("mysql", "user:password@tcp(host:port)/dbname")

	var persons []Person

	rows, _ := cnn.Query("SELECT name, job_id FROM persons")

	defer rows.Close()

	for rows.Next() {
		var person Person
		if err := rows.Scan(&person.Name, &person.JobID); err != nil {
			log.Fatal(err)
		}
		persons = append(persons, person)
	}

	for _, person := range persons {

		var job Job
		err := cnn.QueryRow("SELECT job_id, name FROM Jobs WHERE job_id = ?", person.JobID).Scan(&job.JobID, &job.Name) //ERROR "this query is called in a loop"
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(person.Name, job.Name)
	}

}
