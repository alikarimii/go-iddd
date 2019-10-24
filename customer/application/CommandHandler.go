package application

import (
	"database/sql"
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"reflect"

	"github.com/cockroachdb/errors"
)

const maxCommandHandlerRetries = 10

type CommandHandler struct {
	repositorySessionFactory StartsRepositorySessions
	db                       *sql.DB
}

/*** Factory Method ***/

func NewCommandHandler(repositorySessionFactory StartsRepositorySessions, db *sql.DB) *CommandHandler {
	return &CommandHandler{
		repositorySessionFactory: repositorySessionFactory,
		db:                       db,
	}
}

/*** Implement shared.CommandHandler ***/

func (handler *CommandHandler) Handle(command shared.Command) error {
	if err := handler.assertIsValid(command); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrCommandIsInvalid), "commandHandler.Handle")
	}

	if err := handler.assertIsKnown(command); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrCommandIsUnknown), "commandHandler.Handle")
	}

	if err := handler.handleRetry(command); err != nil {
		return errors.Wrapf(err, "commandHandler.Handle: [%s]", command.CommandName())
	}

	return nil
}

/*** Chain of handler functions ***/

func (handler *CommandHandler) handleRetry(command shared.Command) error {
	var err error
	var retries uint

	for retries = 0; retries < maxCommandHandlerRetries; retries++ {
		// call next method in chain
		if err = handler.handleSession(command); err == nil {
			break // no need to retry, handling was successful
		}

		if errors.Is(err, shared.ErrConcurrencyConflict) {
			continue // retry to resolve the concurrency conflict
		} else {
			break // don't retry for different errors
		}
	}

	if err != nil {
		if retries == maxCommandHandlerRetries {
			return errors.Wrap(err, shared.ErrMaxRetriesExceeded.Error())
		}

		return err // either to many retries or a different error
	}

	return nil
}

func (handler *CommandHandler) handleSession(command shared.Command) error {
	tx, errTx := handler.db.Begin()
	if errTx != nil {
		return errors.Mark(errTx, shared.ErrTechnical)
	}

	persistableCustomers := handler.repositorySessionFactory.StartSession(tx)

	// call next method in chain
	if err := handler.handleCommand(persistableCustomers, command); err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return errors.Wrap(err, errTx.Error())
		}

		return err
	}

	if errTx := tx.Commit(); errTx != nil {
		return errors.Mark(errTx, shared.ErrTechnical)
	}

	return nil
}

func (handler *CommandHandler) handleCommand(
	persistableCustomers PersistableCustomers,
	command shared.Command,
) error {

	var err error
	var recordedEvents shared.DomainEvents

	switch actualCommand := command.(type) {
	case *commands.Register:
		err = handler.register(persistableCustomers, actualCommand)
	case *commands.ConfirmEmailAddress:
		if recordedEvents, err = handler.confirmEmailAddress(persistableCustomers, actualCommand); err == nil {
			err = handler.persist(persistableCustomers, actualCommand.ID(), recordedEvents)
		}
	case *commands.ChangeEmailAddress:
		if recordedEvents, err = handler.changeEmailAddress(persistableCustomers, actualCommand); err == nil {
			err = handler.persist(persistableCustomers, actualCommand.ID(), recordedEvents)
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) persist(customers PersistsCustomers, id *values.ID, recordedEvents shared.DomainEvents) error {
	if err := customers.Persist(id, recordedEvents); err != nil {
		return err
	}

	return nil
}

/*** Business cases ***/

func (handler *CommandHandler) register(
	customers domain.Customers,
	register *commands.Register,
) error {

	recordedEvents := domain.RegisterCustomer(register)

	if err := customers.Register(register.ID(), recordedEvents); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) confirmEmailAddress(
	customers PersistableCustomers,
	confirmEmailAddress *commands.ConfirmEmailAddress,
) (shared.DomainEvents, error) {

	customer, err := customers.Of(confirmEmailAddress.ID())
	if err != nil {
		return nil, err
	}

	recordeEvents := customer.ConfirmEmailAddress(confirmEmailAddress)
	// TODO: map error events to errors!

	return recordeEvents, nil
}

func (handler *CommandHandler) changeEmailAddress(
	customers PersistableCustomers,
	changeEmailAddress *commands.ChangeEmailAddress,
) (shared.DomainEvents, error) {

	customer, err := customers.Of(changeEmailAddress.ID())
	if err != nil {
		return nil, err
	}

	recordeEvents := customer.ChangeEmailAddress(changeEmailAddress)

	return recordeEvents, nil
}

/*** Command Assertions ***/

func (handler *CommandHandler) assertIsValid(command shared.Command) error {
	if command == nil {
		return errors.New("command is nil interface")
	}

	if reflect.ValueOf(command).IsNil() {
		return errors.Newf("[%s]: command value is nil pointer", command.CommandName())
	}

	if reflect.ValueOf(command.AggregateID()).IsNil() {
		return errors.Newf("[%s]: command was not properly created", command.CommandName())
	}

	return nil
}

func (handler *CommandHandler) assertIsKnown(command shared.Command) error {
	switch command.(type) {
	case *commands.Register, *commands.ConfirmEmailAddress, *commands.ChangeEmailAddress:
		return nil
	default:
		return errors.Newf("[%s] command is unknown", command.CommandName())
	}
}
