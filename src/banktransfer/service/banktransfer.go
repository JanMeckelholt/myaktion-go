package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/JanMeckelholt/myaktion-go/src/banktransfer/grpc/banktransfer"
	"github.com/JanMeckelholt/myaktion-go/src/banktransfer/kafka"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

const retryTime = 5 * time.Second

type BankTransferService struct {
	banktransfer.BankTransferServer
	keyGenerator      *KeyGenerator
	transactionWriter kafka.TransactionWriter
}

func NewBankTransferService() *BankTransferService {
	return &BankTransferService{keyGenerator: NewKeyGenerator()}
}
func (s *BankTransferService) TransferMoney(_ context.Context, transaction *banktransfer.Transaction) (*emptypb.Empty, error) {
	entry := log.WithField("transaction", transaction)
	entry.Info("Received transaction")
	s.processTransaction(transaction)
	return &emptypb.Empty{}, nil
}

func (s *BankTransferService) ProcessTransactions(stream banktransfer.BankTransfer_ProcessTransactionsServer) error {
	return func() error {
		r := kafka.NewTransactionReader()
		for {
			err := r.Read(stream.Context(), func(transaction *banktransfer.Transaction) error {
				id := transaction.Id
				entry := log.WithField("transaction", transaction)
				if err := stream.Send(transaction); err != nil {
					return fmt.Errorf("error sending transaction: %w", err)
				}
				response, err := stream.Recv()
				if err != nil {
					entry.WithError(err).Error("Error receiving processing response")
				}
				if response.Id != id {
					return errors.New("received processing response of a different transaction")
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		if err := r.Close(); err != nil {
			return nil
		}
		return nil
	}()
}

func (s *BankTransferService) processTransaction(transaction *banktransfer.Transaction) {
	entry := log.WithField("transaction", transaction)
	go func(transaction banktransfer.Transaction) {
		entry.Info("Start processing transaction")
		transaction.Id = s.keyGenerator.getUniqueId()
		if err := s.transactionWriter.Write(&transaction); err != nil {
			entry.WithError(err).Error("Can't write transaction to transaction writer")
			return
		}
		entry.Info("Transaction forwarded to transaction writer. Processing transaction finished")
	}(*transaction)
}

func (s *BankTransferService) Start() {
	log.Info("Starting banktransfer service")
	s.transactionWriter = kafka.NewTransactionWriter()
	log.Info("Successfully created transaction writer")
}

func (s *BankTransferService) Stop() {
	log.Info("Stopping banktransfer service")
	s.transactionWriter.Close()
	log.Info("Successfully closed connection to transaction writer")
}
