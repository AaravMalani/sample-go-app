# CVWO Assignment Sample Golang App

## Running the app
1. [Clone](https://docs.github.com/en/get-started/quickstart/fork-a-repo#cloning-your-forked-repository) the repo.
2. Open your terminal and navigate to the directory containing your cloned project.
3. Run your database instance (for example, `docker run --name postgres-test -e POSTGRES_PASSWORD=<pw> -p 5432:5432 -d postgres`)
4. Push the initial SQL data (for example, `cat gen/init.sql | docker exec -i postgres-test psql -U postgres -d postgres`)
3. Copy `.env.example` to `.env` and complete the fields 
4. Run `go run gen/main.go` to generate the models and queries with GORM
5. Run `go run cmd/server/main.go` and head over to http://localhost:8000/topics/list to view the response.


### Navigating the code
```
.
├── cmd
│   └── server
├── gen
├── internal
│   ├── api         # Encapsulates types and utilities related to the API
│   ├── dataaccess   # Data Access layer accesses data from the database
│   ├── database    # Encapsulates the types and utilities related to the database
│   ├── handlers    # Handler functions to respond to requests
│   ├── router      # Encapsulates types and utilities related to the router
│   └── routes      # Defines routes that are used in the application
├── model           # Definitions of models used in the application
├── README.md
├── go.mod
└── go.sum
```

Main directories/files to note:
* `cmd` contains the main entry point for the application
* `internal` holds most of the functional code for the project that is specific to the core logic of your application
* `README.md` is a form of documentation about the project. It is what you are reading right now.
* `go.mod` contains important metadata, for example, the dependencies in the project. See [here](https://go.dev/ref/mod) for more information
* `go.sum` See [here](https://go.dev/ref/mod) for more information
