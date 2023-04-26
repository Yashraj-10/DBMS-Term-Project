# GO CLI Application

This application provides a report of the metrics of a database server. The metrics can either be of the whole database server or of a particular database.

### Building the application

To run the file, use the following command:

```bash
    cd cmd/pgmetrics
    go build main.go csv.go report.go
```
This will create an executable file in the ```cmd/pgmetrics``` directory.

### Running the application

To run the executable file, use the following commands in the ```cmd/pgmetrics``` directory:
- For the whole database server:
```bash
    ./main --username=USERNAME > OUTPUT_FILE_NAME
```
- For a particular database:
```bash
    ./main --username=USERNAME DATABASE_NAME > OUTPUT_FILE_NAME
    # or
    ./main --username=USERNAME --database=DATABASE_NAME > OUTPUT_FILE_NAME
```

The OUTPUT_FILE_NAME will contain the report of the metrics of the database server or the particular database.