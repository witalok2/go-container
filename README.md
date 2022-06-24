# Go container

<p align="center">
<a href=""><img src="https://miro.medium.com/max/880/1*136qhXxInh44-pWrPrkLTw.png" align="center" height="321" width="326" ></a>
</p>

<a href="https://github.com/witalok2/go-container/actions/workflows/test.yaml">
    <img src="https://github.com/witalok2/go-container/workflows/Test/badge.svg?style=flat" />
</a>

Run test 
```sh
run 'make test'
```

### 📋 Prerequisites

Tools: 
- [Golang](https://golang.org/doc/install)

### Por que devo usar?
Use Docker to run your Golang integration tests on 3rd party services on Microsoft Windows, Mac OSX and Linux!

When developing applications, it is often necessary to use services that communicate with a database system or even a redis. Unit testing these services can be tricky because simulating database/DBAL is strenuous. Making small changes to the schema entails rewriting at least some if not all of the mocks. The same goes for API changes in DBAL. To avoid this, it is smarter to test these specific services against a real database which is destroyed after testing. Docker is the perfect system for running unit tests as you can launch containers in a few seconds and kill them when the test is complete. The Dockertest library provides easy-to-use commands to launch Docker containers and use them for your tests.

## Supported / Driver
Databases | Driver
----------|-----------
PostgreSQL and compatible databases (e.g. CockroachDB) | https://github.com/lib/pq
MongoDB and compatible databases | https://github.com/mongodb/mongo-go-driver
Redis | https://github.com/go-redis/redis

### Using Container
```go
// exemple
func (suite *ContainerTestSuite) TestMongoContainer() {
	dockerContainer := NewContainer(dockertest.RunOptions{
		Repository:   "mongo",
		Tag:          "3.3.12",
		ExposedPorts: []string{"3000"},
		Cmd:          []string{"mongod", "--smallfiles", "--port", "3000"},
	})

	defer dockerContainer.Pool.Purge(dockerContainer.Resource)

	port := dockerContainer.Resource.GetPort("3000/tcp")
	suite.NotEmpty(port)

	err := dockerContainer.Pool.Retry(func() error {
		response, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s", port))
		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("could not connect to resource")
		}

		defer response.Body.Close()

		return nil
	})
	suite.Require().Nil(err)
}
```