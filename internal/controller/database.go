package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"gorch/pkg/deployment"
	"gorch/pkg/pod"
	"log/slog"
	"slices"
	"time"
)

type Collection string

const (
	Workers     Collection = "workers"
	Deployments Collection = "deployments"
	Pods        Collection = "pods"
)

type Database struct {
	db *bbolt.DB
}

func NewDatabase() (*Database, error) {
	db, err := bbolt.Open("gorch.db", 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	database := &Database{
		db: db,
	}

	if err := database.migrate(); err != nil {
		db.Close()
		return nil, err
	}

	return database, nil
}

func (d *Database) migrate() error {
	tx, err := d.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.CreateBucket([]byte(Workers))
	if err != nil {
		slog.Info("migration", slog.Any("db", Workers), slog.Any("message", err))
	}

	_, err = tx.CreateBucket([]byte(Deployments))
	if err != nil {
		slog.Info("migration", slog.Any("db", Deployments), slog.Any("message", err))
	}

	_, err = tx.CreateBucket([]byte(Pods))
	if err != nil {
		slog.Info("migration", slog.Any("db", Pods), slog.Any("message", err))
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil

}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) save(collection Collection, id string, entry any) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("mashal data: %w", err)
	}

	return d.db.Update(
		func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte(collection))
			return bucket.Put([]byte(id), data)
		},
	)
}

func (d *Database) loadAll(collection Collection) ([]any, error) {
	result := make([]any, 0)
	err := d.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(collection))
		return bucket.ForEach(func(k, v []byte) error {
			result = append(result, v)
			return nil
		})
	})

	if err != nil {
		return nil, fmt.Errorf("load entries: %w", err)
	}
	return result, nil
}

func unmarshalListToType[T any](data []any) ([]T, error) {
	result := make([]T, 0, len(data))

	for d := range slices.Values(data) {
		entry, err := unmarshalToType[T](d)
		if err != nil {
			return nil, err
		}
		result = append(result, entry)
	}

	return result, nil
}

func unmarshalToType[T any](data any) (T, error) {
	var entry T
	err := json.Unmarshal(data.([]byte), &entry)
	if err != nil {
		return *new(T), fmt.Errorf("unmarshal entry: %w", err)
	}
	return entry, nil
}

func (d *Database) SaveWorker(worker WorkerEntity) error {
	return d.save(Workers, worker.Name, worker)
}

func (d *Database) LoadWorkers() ([]WorkerEntity, error) {
	jsonData, err := d.loadAll(Workers)
	if err != nil {
		return nil, fmt.Errorf("load workers: %w", err)
	}

	return unmarshalListToType[WorkerEntity](jsonData)
}

func (d *Database) SaveDeployment(deployment deployment.Deployment) error {
	return d.save(Deployments, deployment.ID, deployment)
}

func (d *Database) LoadDeployments() ([]deployment.Deployment, error) {
	res, err := d.loadAll(Deployments)
	if err != nil {
		return nil, fmt.Errorf("load deployments: %w", err)
	}
	return unmarshalListToType[deployment.Deployment](res)
}

func (d *Database) SavePod(pod pod.Pod) error {
	pod.LastUpdated = time.Now().UTC()
	return d.save(Pods, pod.ID, pod)
}

func (d *Database) LoadPodsByDeploymentName(name []byte) ([]pod.Pod, error) {
	raw := make([]any, 0)

	err := d.db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket([]byte(Pods)).Cursor()
		for k, v := c.Seek(name); k != nil && bytes.HasPrefix(k, name); k, v = c.Next() {
			raw = append(raw, v)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("load pods by deployment %s: %w", name, err)
	}

	return unmarshalListToType[pod.Pod](raw)
}

func (d *Database) LoadPod(id string) (*pod.Pod, error) {
	var raw []byte
	err := d.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(Pods))
		raw = bucket.Get([]byte(id))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("load pod: %w", err)
	}

	p, err := unmarshalToType[pod.Pod](raw)
	if err != nil {
		return nil, fmt.Errorf("unmarshal pod: %w", err)
	}
	return &p, nil
}

func (d *Database) LoadWorker(id string) (*WorkerEntity, error) {
	var raw []byte
	err := d.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(Workers))
		raw = bucket.Get([]byte(id))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("load worker: %w", err)
	}

	w, err := unmarshalToType[WorkerEntity](raw)
	if err != nil {
		return nil, fmt.Errorf("unmarshal worker: %w", err)
	}
	return &w, nil
}
