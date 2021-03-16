package postgres

import (
	"database/sql"
	"math"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
)

const identityStreamPrefix = "identity"

type IdentityEventStore struct {
	db                   *sql.DB
	tx                   *sql.Tx
	retrieveEventStream  forRetrievingEventStreams
	appendEventsToStream forAppendingEventsToStreams
	purgeEventStream     forPurgingEventStreams
	marshalDomainEvent   es.MarshalDomainEvent
	unmarshalDomainEvent es.UnmarshalDomainEvent
}

func NewIdentityEventStore(
	db *sql.DB,
	retrieveEventStream forRetrievingEventStreams,
	appendEventsToStream forAppendingEventsToStreams,
	purgeEventStream forPurgingEventStreams,
	marshalDomainEvent es.MarshalDomainEvent,
	unmarshalDomainEvent es.UnmarshalDomainEvent,
) *IdentityEventStore {

	return &IdentityEventStore{
		db:                   db,
		retrieveEventStream:  retrieveEventStream,
		appendEventsToStream: appendEventsToStream,
		purgeEventStream:     purgeEventStream,
		marshalDomainEvent:   marshalDomainEvent,
		unmarshalDomainEvent: unmarshalDomainEvent,
	}
}

func (s *IdentityEventStore) WithTx(tx *sql.Tx) application.ForStoringIdentityEventStreams {
	return &IdentityEventStore{
		db:                   s.db,
		tx:                   tx,
		retrieveEventStream:  s.retrieveEventStream,
		appendEventsToStream: s.appendEventsToStream,
		purgeEventStream:     s.purgeEventStream,
		marshalDomainEvent:   s.marshalDomainEvent,
		unmarshalDomainEvent: s.unmarshalDomainEvent,
	}
}

func (s *IdentityEventStore) RetrieveEventStream(id value.IdentityID) (es.EventStream, error) {
	wrapWithMsg := "identityEventStore.RetrieveEventStream"

	eventStream, err := s.retrieveEventStream(s.streamID(id), 0, math.MaxUint32, s.db, s.unmarshalDomainEvent)
	if err != nil {
		return nil, errors.Wrap(err, wrapWithMsg)
	}

	if len(eventStream) == 0 {
		err := errors.New("identity not found")
		return nil, shared.MarkAndWrapError(err, shared.ErrNotFound, wrapWithMsg)
	}

	return eventStream, nil
}

func (s *IdentityEventStore) StartEventStream(identityRegistered domain.IdentityRegistered) error {
	var err error
	wrapWithMsg := "identityEventStore.StartEventStream"

	recordedEvents := []es.DomainEvent{identityRegistered}

	streamID := s.streamID(identityRegistered.IdentityID())

	if err = s.appendEventsToStream(streamID, recordedEvents, s.marshalDomainEvent, s.tx); err != nil {
		if errors.Is(err, shared.ErrConcurrencyConflict) {
			return shared.MarkAndWrapError(errors.New("found duplicate identity"), shared.ErrDuplicate, wrapWithMsg)
		}

		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}

func (s *IdentityEventStore) AppendToEventStream(recordedEvents es.RecordedEvents, id value.IdentityID) error {
	var err error
	wrapWithMsg := "identityEventStore.AppendToEventStream"

	if err = s.appendEventsToStream(s.streamID(id), recordedEvents, s.marshalDomainEvent, s.tx); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}

func (s *IdentityEventStore) PurgeEventStream(id value.IdentityID) error {
	var err error
	wrapWithMsg := "identityEventStore.PurgeEventStream"

	tx, err := s.db.Begin()
	if err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	if err = s.purgeEventStream(s.streamID(id), tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = tx.Commit(); err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	return nil
}

func (s *IdentityEventStore) streamID(id value.IdentityID) es.StreamID {
	return es.BuildStreamID(identityStreamPrefix + "-" + id.String())
}
