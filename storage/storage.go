package storage

import (
	"github.com/danielnegri/adheretech/ledger"
)

const StartTimeKey = "storage-start-time"

type Storage interface {
	ledger.Ledger
	ledger.Checker
}
