package postgres

import (
	"context"
	"io/ioutil"
	"net/url"
	"path"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	logger "github.com/sirupsen/logrus"
	container "github.com/witalok2/go-container/container"
)

type PostgresContainer struct {
	url        *url.URL
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	postgresDB *pgxpool.Pool
}

func NewContainerPostgres() *PostgresContainer {
	url := &url.URL{
		Scheme:      "postgres",
		User:        url.UserPassword("myuser", "secret"),
		Path:        "mydatabase",
		RawFragment: "sslmode=disable",
	}

	dockerContainer := container.NewContainer(dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "9.5",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=" + url.User.Username(),
			"POSTGRES_DB=" + url.Path,
		},
		ExposedPorts: []string{"5432"},
	})

	err := dockerContainer.Resource.Expire(30 * 60)
	if err != nil {
		logger.WithError(err).Fatalf("could not expire dockerContainer")
	}

	url.Host = dockerContainer.Resource.Container.NetworkSettings.IPAddress

	var db *pgxpool.Pool
	err = dockerContainer.Pool.Retry(func() error {
		db, err = pgxpool.Connect(context.Background(), url.String())
		if err != nil {
			return err
		}

		return db.Ping(context.Background())
	})
	if err != nil {
		logger.WithError(err).Fatalf("could not connect to postgres server")
	}

	return &PostgresContainer{
		url:        url,
		pool:       dockerContainer.Pool,
		resource:   dockerContainer.Resource,
		postgresDB: db,
	}
}

func (c *PostgresContainer) RunMigrations(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logger.WithError(err).WithField("dir", dir).Fatal("unable to list migrations")
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, f := range files {
		filename := f.Name()
		if filepath.Ext(filename) != ".sql" {
			continue
		}
		if filename[len(filename)-7:] != ".up.sql" {
			continue
		}
		ctx := context.Background()
		row := c.postgresDB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM migrations WHERE name=$1)", filename)
		var applied int
		if err := row.Scan(&applied); err == nil && applied == 1 {
			logger.WithField("filename", filename).Info("migration already applied")
			continue
		}
		c.ExecuteSQLFile(path.Join(dir, filename))
	}
}

func (c *PostgresContainer) ExecuteSQLFile(filename string) {
	sql, _ := ioutil.ReadFile(filename)
	ctx := context.Background()
	if _, err := c.postgresDB.Exec(ctx, string(sql)); err != nil {
		logger.WithError(err).WithField("filename", filename).Fatal("could not run SQL from file")
	}
}

func (c *PostgresContainer) Purge() {
	if err := c.pool.Purge(c.resource); err != nil {
		logger.WithError(err).Fatal("Could not purge resource")
	}
}
