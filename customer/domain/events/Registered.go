package events

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"reflect"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

const (
	registeredAggregateName       = "Customer"
	RegisteredMetaTimestampFormat = time.RFC3339Nano
)

type Registered struct {
	customerID              *values.CustomerID
	confirmableEmailAddress *values.ConfirmableEmailAddress
	personName              *values.PersonName
	meta                    *Meta
}

/*** Factory Methods ***/

func ItWasRegistered(
	customerID *values.CustomerID,
	confirmableEmailAddress *values.ConfirmableEmailAddress,
	personName *values.PersonName,
	streamVersion uint,
) *Registered {

	registered := &Registered{
		customerID:              customerID,
		confirmableEmailAddress: confirmableEmailAddress,
		personName:              personName,
	}

	eventType := reflect.TypeOf(registered).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]
	eventName = strings.Title(eventName)
	fullEventName := registeredAggregateName + eventName

	registered.meta = &Meta{
		identifier:    customerID.String(),
		eventName:     fullEventName,
		occurredAt:    time.Now().Format(RegisteredMetaTimestampFormat),
		streamVersion: streamVersion,
	}

	return registered
}

/*** Getter Methods ***/

func (registered *Registered) CustomerID() *values.CustomerID {
	return registered.customerID
}

func (registered *Registered) ConfirmableEmailAddress() *values.ConfirmableEmailAddress {
	return registered.confirmableEmailAddress
}

func (registered *Registered) PersonName() *values.PersonName {
	return registered.personName
}

/*** Implement shared.DomainEvent ***/

func (registered *Registered) Identifier() string {
	return registered.meta.identifier
}

func (registered *Registered) EventName() string {
	return registered.meta.eventName
}

func (registered *Registered) OccurredAt() string {
	return registered.meta.occurredAt
}

func (registered *Registered) StreamVersion() uint {
	return registered.meta.streamVersion
}

/*** Implement json.Marshaler ***/

func (registered *Registered) MarshalJSON() ([]byte, error) {
	data := &struct {
		CustomerID              *values.CustomerID              `json:"customerID"`
		ConfirmableEmailAddress *values.ConfirmableEmailAddress `json:"confirmableEmailAddress"`
		PersonName              *values.PersonName              `json:"personName"`
		Meta                    *Meta                           `json:"meta"`
	}{
		CustomerID:              registered.customerID,
		ConfirmableEmailAddress: registered.confirmableEmailAddress,
		PersonName:              registered.personName,
		Meta:                    registered.meta,
	}

	return jsoniter.Marshal(data)
}

/*** Implement json.Unmarshaler ***/

func (registered *Registered) UnmarshalJSON(data []byte) error {
	unmarshaledData := &struct {
		CustomerID              *values.CustomerID              `json:"customerID"`
		ConfirmableEmailAddress *values.ConfirmableEmailAddress `json:"confirmableEmailAddress"`
		PersonName              *values.PersonName              `json:"personName"`
		Meta                    *Meta                           `json:"meta"`
	}{}

	if err := jsoniter.Unmarshal(data, unmarshaledData); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrUnmarshalingFailed), "registered.UnmarshalJSON")
	}

	registered.customerID = unmarshaledData.CustomerID
	registered.confirmableEmailAddress = unmarshaledData.ConfirmableEmailAddress
	registered.personName = unmarshaledData.PersonName
	registered.meta = unmarshaledData.Meta

	return nil
}
