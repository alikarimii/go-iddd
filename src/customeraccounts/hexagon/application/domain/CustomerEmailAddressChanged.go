package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type CustomerEmailAddressChanged struct {
	customerID       value.CustomerID
	emailAddress     value.UnconfirmedEmailAddress
	confirmationHash value.ConfirmationHash
	meta             es.EventMeta
}

func BuildCustomerEmailAddressChanged(
	customerID value.CustomerID,
	emailAddress value.UnconfirmedEmailAddress,
	confirmationHash value.ConfirmationHash,
	causationID es.MessageID,
	streamVersion uint,
) CustomerEmailAddressChanged {

	event := CustomerEmailAddressChanged{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
	}

	event.meta = es.BuildEventMeta(event, causationID, streamVersion)

	return event
}

func RebuildCustomerEmailAddressChanged(
	customerID string,
	emailAddress string,
	confirmationHash string,
	meta es.EventMeta,
) CustomerEmailAddressChanged {

	event := CustomerEmailAddressChanged{
		customerID:       value.RebuildCustomerID(customerID),
		emailAddress:     value.RebuildUnconfirmedEmailAddress(emailAddress),
		confirmationHash: value.RebuildConfirmationHash(confirmationHash),
		meta:             meta,
	}

	return event
}

func (event CustomerEmailAddressChanged) CustomerID() value.CustomerID {
	return event.customerID
}

func (event CustomerEmailAddressChanged) EmailAddress() value.UnconfirmedEmailAddress {
	return event.emailAddress
}

func (event CustomerEmailAddressChanged) ConfirmationHash() value.ConfirmationHash {
	return event.confirmationHash
}

func (event CustomerEmailAddressChanged) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerEmailAddressChanged) IsFailureEvent() bool {
	return false
}

func (event CustomerEmailAddressChanged) FailureReason() error {
	return nil
}
