package controller

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"log/slog"
	"slices"
)

type Collection string

const (
	Workers Collection = "workers"
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
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil

}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) SaveWorker(worker WorkerEntity) error {
	data, err := json.Marshal(worker)
	if err != nil {
		return fmt.Errorf("mashal data: %w", err)
	}

	return d.db.Update(
		func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte(Workers))
			return bucket.Put([]byte(worker.Name), data)
		},
	)
}

func (d *Database) LoadWorkers() ([]WorkerEntity, error) {
	jsonData, err := d.loadAll(Workers)
	if err != nil {
		return nil, fmt.Errorf("load workers: %w", err)
	}

	return unmarshalListToType[WorkerEntity](jsonData)
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
