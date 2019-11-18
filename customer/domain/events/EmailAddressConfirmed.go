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
	emailAddressConfirmedAggregateName       = "Customer"
	EmailAddressConfirmedMetaTimestampFormat = time.RFC3339Nano
)

type EmailAddressConfirmed struct {
	customerID   *values.CustomerID
	emailAddress *values.EmailAddress
	meta         *Meta
}

/*** Factory Methods ***/

func EmailAddressWasConfirmed(
	customerID *values.CustomerID,
	emailAddress *values.EmailAddress,
	streamVersion uint,
) *EmailAddressConfirmed {

	emailAddressConfirmed := &EmailAddressConfirmed{
		customerID:   customerID,
		emailAddress: emailAddress,
	}

	eventType := reflect.TypeOf(emailAddressConfirmed).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]
	eventName = strings.Title(eventName)
	fullEventName := emailAddressConfirmedAggregateName + eventName

	emailAddressConfirmed.meta = &Meta{
		identifier:    customerID.String(),
		eventName:     fullEventName,
		occurredAt:    time.Now().Format(EmailAddressConfirmedMetaTimestampFormat),
		streamVersion: streamVersion,
	}

	return emailAddressConfirmed
}

/*** Getter Methods ***/

func (emailAddressConfirmed *EmailAddressConfirmed) CustomerID() *values.CustomerID {
	return emailAddressConfirmed.customerID
}

func (emailAddressConfirmed *EmailAddressConfirmed) EmailAddress() *values.EmailAddress {
	return emailAddressConfirmed.emailAddress
}

/*** Implement shared.DomainEvent ***/

func (emailAddressConfirmed *EmailAddressConfirmed) Identifier() string {
	return emailAddressConfirmed.meta.identifier
}

func (emailAddressConfirmed *EmailAddressConfirmed) EventName() string {
	return emailAddressConfirmed.meta.eventName
}

func (emailAddressConfirmed *EmailAddressConfirmed) OccurredAt() string {
	return emailAddressConfirmed.meta.occurredAt
}

func (emailAddressConfirmed *EmailAddressConfirmed) StreamVersion() uint {
	return emailAddressConfirmed.meta.streamVersion
}

/*** Implement json.Marshaler ***/

func (emailAddressConfirmed *EmailAddressConfirmed) MarshalJSON() ([]byte, error) {
	data := &struct {
		CustomerID   *values.CustomerID   `json:"customerID"`
		EmailAddress *values.EmailAddress `json:"emailAddress"`
		Meta         *Meta                `json:"meta"`
	}{
		CustomerID:   emailAddressConfirmed.customerID,
		EmailAddress: emailAddressConfirmed.emailAddress,
		Meta:         emailAddressConfirmed.meta,
	}

	return jsoniter.Marshal(data)
}

/*** Implement json.Unmarshaler ***/

func (emailAddressConfirmed *EmailAddressConfirmed) UnmarshalJSON(data []byte) error {
	unmarshaledData := &struct {
		CustomerID   *values.CustomerID   `json:"customerID"`
		EmailAddress *values.EmailAddress `json:"emailAddress"`
		Meta         *Meta                `json:"meta"`
	}{}

	if err := jsoniter.Unmarshal(data, unmarshaledData); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrUnmarshalingFailed), "emailAddressConfirmed.UnmarshalJSON")
	}

	emailAddressConfirmed.customerID = unmarshaledData.CustomerID
	emailAddressConfirmed.emailAddress = unmarshaledData.EmailAddress
	emailAddressConfirmed.meta = unmarshaledData.Meta

	return nil
}
