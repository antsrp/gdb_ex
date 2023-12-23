-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Departments 
(
	id SERIAL PRIMARY KEY,
	name VARCHAR (50) UNIQUE
);

CREATE TABLE IF NOT EXISTS Employees 
(
	id SERIAL PRIMARY KEY,
	name VARCHAR (30),
	surname VARCHAR (30),
	age int,
	salary int,
	departmentID int,
	FOREIGN KEY (departmentID)
        REFERENCES Departments (id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS Projects 
(
	id SERIAL PRIMARY KEY,
	name VARCHAR (50)
);

CREATE TABLE IF NOT EXISTS Projects_Employees
(
	projectId int,
	employeeId int,
	FOREIGN KEY(projectId)
	REFERENCES Projects (id) ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY(employeeId)
	REFERENCES Employees (id) ON DELETE CASCADE ON UPDATE CASCADE
);;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Projects_Employees;
DROP TABLE Projects;
DROP TABLE Employees;
DROP TABLE Departments;
-- +goose StatementEnd
