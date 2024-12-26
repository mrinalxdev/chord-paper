package storage

import (
    "github.com/dgraph-io/badger/v3"
)

type BadgerStore struct {
    db *badger.DB
}

func NewBadgerStore(path string) (*BadgerStore, error) {
    opts := badger.DefaultOptions(path)
    db, err := badger.Open(opts)
    if err != nil {
        return nil, err
    }
    
    return &BadgerStore{db: db}, nil
}

func (bs *BadgerStore) Put(key []byte, value []byte) error {
    return bs.db.Update(func(txn *badger.Txn) error {
        return txn.Set(key, value)
    })
}

func (bs *BadgerStore) Get(key []byte) ([]byte, error) {
    var valCopy []byte
    err := bs.db.View(func(txn *badger.Txn) error {
        item, err := txn.Get(key)
        if err != nil {
            return err
        }
        
        valCopy, err = item.ValueCopy(nil)
        return err
    })
    
    return valCopy, err
}

func (bs *BadgerStore) Delete(key []byte) error {
    return bs.db.Update(func(txn *badger.Txn) error {
        return txn.Delete(key)
    })
}

func (bs *BadgerStore) Close() error {
    return bs.db.Close()
}