package main

import (
	"context"
	"github.com/JanMeckelholt/myaktion-go/src/myaktion/client"
	"github.com/JanMeckelholt/myaktion-go/src/myaktion/client/banktransfer"
	log "github.com/sirupsen/logrus"
	"time"
)

func monitortransactions() {
	for {
		connectandmonitor()
		time.Sleep(time.Second)
	}
}

func connectandmonitor() {
	conn, err := client.GetBankTransferConnection()
	if err != nil {
		log.WithError(err).Fatal("error connecting to the banktransfer service")
	}
	defer conn.Close()
	// TODO: here we force a deadline after 10 seconds to test the re-connecting logic
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	banktransferClient := banktransfer.NewBankTransferClient(conn)
	watcher, err := banktransferClient.ProcessTransactions(ctx)
	if err != nil {
		log.WithError(err).Fatal("error watching transactions")
	}
	log.Info("Successfully connected to banktransfer service for processing transactions")
	for {
		transaction, err := watcher.Recv()
		if err != nil {
			if _, deadline := ctx.Deadline(); deadline {
				log.Info("deadline reached. reconnect client")
				break
			}
			log.WithError(err).Error("error receiving transaction")
			continue
		}
		entry := log.WithField("transaction", transaction)
		entry.Info("Received transaction. Sending processing response")
		err = watcher.Send(&banktransfer.ProcessingResponse{Id: transaction.Id})
		if err != nil {
			entry.WithError(err).Error("error sending processing response")
			continue
		}
		entry.Info("Processing response sent")
	}
}
