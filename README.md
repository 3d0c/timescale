Please note, this isn't a real application, just a homework (test task).  

## TimeScaleDB SELECT queries benchmark

### Prerequisites
Please follow [this](https://docs.timescale.com/v0.9/getting-started/installation/mac/installation-source) instruction to install PostgreSQL and TimeScale module.  

### Installation

```sh
go get github.com/3d0c/timescale
```

Now you can install it in your GOPATH or use right from the repository.

```sh
go install ./...
```

If you have `GOPATH/bin` in your environment, the `ts-bench` application will be available.

### Prepare sample data

```sh
psql -U postgres < fixtures/cpu_usage.sql
psql -U postgres -d homework -c "\COPY cpu_usage FROM fixtures/cpu_usage.csv CSV HEADER"
```

### Running
This application requires special `csv` file which has specification for queries generation. There are few in `fixtures/` folder.

Run:

```sh
ts-bench -qp=./fixtures/qp20.csv -wnum=10
```

Please specify your postgres connection string by providing `-dbargs` flag!

After some time it should produce output like:

```
Total queries: 61
Distribution across workers:
	worker #0  61 queries
Duration:
	total:   2.062009658s
	minimum: 33.017833ms
	maximum: 72.919263ms
	median:  33.530845ms
	average: 52.968548ms

```




