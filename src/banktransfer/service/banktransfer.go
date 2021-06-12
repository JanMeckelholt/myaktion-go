package service

import (
	"context"
	"github.com/JanMeckelholt/myaktion-go/src/banktransfer/grpc/banktransfer"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
	"sync/atomic"
	"time"
)

const retryTime = 5 * time.Second

type BankTransferService struct {
	banktransfer.BankTransferServer
	counter    int32
	queue      chan banktransfer.Transaction
	retryQueue chan banktransfer.Transaction
	stop       chan struct{}
}

func NewBankTransferService() *BankTransferService {
	rand.Seed(time.Now().UnixNano())
	return &BankTransferService{counter: 0,
		queue:      make(chan banktransfer.Transaction),
		retryQueue: make(chan banktransfer.Transaction),
		stop:       make(chan struct{}),
	}
}
func (s *BankTransferService) TransferMoney(_ context.Context, transaction *banktransfer.Transaction) (*emptypb.Empty, error) {
	entry := log.WithField("transaction", transaction)
	entry.Info("Received transaction")
	s.processTransaction(transaction)
	return &emptypb.Empty{}, nil
}

func (s *BankTransferService) ProcessTransactions(stream banktransfer.BankTransfer_ProcessTransactionsServer) error {
	//ticker := time.NewTicker(2 * time.Second)
	//defer ticker.Stop()
	return func() error {
		for {
			select {
			case <-stream.Context().Done():
				log.Info("Watching transactions cancelled from the client side")
				return nil
			case transaction := <-s.queue:
				id := transaction.Id
				entry := log.WithField("transaction", transaction)
				entry.Info("Sending transaction")
				if err := stream.Send(&transaction); err != nil {
					s.requeuTransaction(&transaction)
					entry.WithError(err).Error("Error sending transaction")
					return err
				}
				entry.Info("Transaction send. Waiting for processing response")
				response, err := stream.Recv()
				if err != nil {
					s.requeuTransaction(&transaction)
					entry.WithError(err).Error("Error receiving processing response")
					return err
				}
				if response.Id != id {
					// NOTE: this is just a guard and not happening as transaction is local per connection
					entry.Error("Received processing response of a different transaction")
				} else {
					entry.Info("Processing response received")
				}
				//case _ = <-ticker.C:
				//	transaction := &banktransfer.Transaction{Id: s.counter, Amount: 20}
				//	entry := log.WithField("transaction", transaction)
				//	entry.Info("Sending transaction")
				//	if err := stream.Send(transaction); err != nil {
				//		entry.WithError(err).Error("Error sending transaction")
				//		return err
				//	}
				//	entry.Info("Transaction sent. Waiting for processing response")
				//	response, err := stream.Recv()
				//	if err != nil {
				//		entry.WithError(err).Error("Error receiving processing response")
				//		return err
				//	}
				//	if response.Id != s.counter {
				//		// NOTE: this is just a guard and not happening as transaction is local per connection
				//		entry.Error("Received processing response of a different transaction")
				//	} else {
				//		entry.Info("Processing response received")
				//		s.counter++
				//	}

			}
		}
	}()
}

func (s *BankTransferService) processTransaction(transaction *banktransfer.Transaction) {
	entry := log.WithField("transaction", transaction)
	go func(transaction banktransfer.Transaction) {
		entry.Info("Start processing transaction")
		transaction.Id = s.getUniqueId()
		time.Sleep(time.Duration(rand.Intn(9)+1) * time.Second)
		s.queue <- transaction
		entry.Info("Processing transaction finished")
	}(*transaction)
}

func (s *BankTransferService) requeuTransaction(transaction *banktransfer.Transaction) {
	entry := log.WithField("transaction", transaction)
	go func(transaction banktransfer.Transaction) {
		entry.Infof("Requeuing transaction. Wait for %d seconds", retryTime.Seconds())
		time.Sleep(retryTime)
		s.retryQueue <- transaction
		entry.Info("Transaction requeued")
	}(*transaction)
}

func (s *BankTransferService) getUniqueId() int32 {
	return atomic.AddInt32(&s.counter, 1)
}

func (s *BankTransferService) Start() {
	log.Infoln("Starting banktransfer service")
	go func() {
		for {
			select {
			case <-s.stop:
				break
			case transaction := <-s.retryQueue:
				s.queue <- transaction
			}
		}
	}()
}

func (s *BankTransferService) Stop() {
	log.Infoln("Stopping banktransfer service")
	close(s.stop)
}
