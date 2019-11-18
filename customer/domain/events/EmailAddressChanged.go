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
	emailAddressChangedAggregateName       = "Customer"
	EmailAddressChangedMetaTimestampFormat = time.RFC3339Nano
)

type EmailAddressChanged struct {
	customerID              *values.CustomerID
	confirmableEmailAddress *values.ConfirmableEmailAddress
	meta                    *Meta
}

/*** Factory Methods ***/

func EmailAddressWasChanged(
	customerID *values.CustomerID,
	confirmableEmailAddress *values.ConfirmableEmailAddress,
	streamVersion uint,
) *EmailAddressChanged {

	emailAddressChanged := &EmailAddressChanged{
		customerID:              customerID,
		confirmableEmailAddress: confirmableEmailAddress,
	}

	eventType := reflect.TypeOf(emailAddressChanged).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]
	eventName = strings.Title(eventName)
	fullEventName := emailAddressChangedAggregateName + eventName

	emailAddressChanged.meta = &Meta{
		identifier:    customerID.String(),
		eventName:     fullEventName,
		occurredAt:    time.Now().Format(EmailAddressChangedMetaTimestampFormat),
		streamVersion: streamVersion,
	}

	return emailAddressChanged
}

/*** Getter Methods ***/

func (emailAddressChanged *EmailAddressChanged) CustomerID() *values.CustomerID {
	return emailAddressChanged.customerID
}

func (emailAddressChanged *EmailAddressChanged) ConfirmableEmailAddress() *values.ConfirmableEmailAddress {
	return emailAddressChanged.confirmableEmailAddress
}

/*** Implement shared.DomainEvent ***/

func (emailAddressChanged *EmailAddressChanged) Identifier() string {
	return emailAddressChanged.meta.identifier
}

func (emailAddressChanged *EmailAddressChanged) EventName() string {
	return emailAddressChanged.meta.eventName
}

func (emailAddressChanged *EmailAddressChanged) OccurredAt() string {
	return emailAddressChanged.meta.occurredAt
}

func (emailAddressChanged *EmailAddressChanged) StreamVersion() uint {
	return emailAddressChanged.meta.streamVersion
}

/*** Implement json.Marshaler ***/

func (emailAddressChanged *EmailAddressChanged) MarshalJSON() ([]byte, error) {
	data := &struct {
		CustomerID              *values.CustomerID              `json:"customerID"`
		ConfirmableEmailAddress *values.ConfirmableEmailAddress `json:"confirmableEmailAddress"`
		Meta                    *Meta                           `json:"meta"`
	}{
		CustomerID:              emailAddressChanged.customerID,
		ConfirmableEmailAddress: emailAddressChanged.confirmableEmailAddress,
		Meta:                    emailAddressChanged.meta,
	}

	return jsoniter.Marshal(data)
}

/*** Implement json.Unmarshaler ***/

func (emailAddressChanged *EmailAddressChanged) UnmarshalJSON(data []byte) error {
	unmarshaledData := &struct {
		CustomerID              *values.CustomerID              `json:"customerID"`
		ConfirmableEmailAddress *values.ConfirmableEmailAddress `json:"confirmableEmailAddress"`
		Meta                    *Meta                           `json:"meta"`
	}{}

	if err := jsoniter.Unmarshal(data, unmarshaledData); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrUnmarshalingFailed), "emailAddressChanged.UnmarshalJSON")
	}

	emailAddressChanged.customerID = unmarshaledData.CustomerID
	emailAddressChanged.confirmableEmailAddress = unmarshaledData.ConfirmableEmailAddress
	emailAddressChanged.meta = unmarshaledData.Meta

	return nil
}
