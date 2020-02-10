package es_test

import (
	"go-iddd/service/lib/es"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewStreamID(t *testing.T) {
	Convey("Given valid input", t, func() {
		streamIDInput := "customer-123"

		Convey("When a new StreamID is created", func() {
			streamID := es.NewStreamID(streamIDInput)

			Convey("It should succeed", func() {
				So(streamID, ShouldNotBeNil)
			})
		})
	})

	Convey("Given empty input", t, func() {
		streamIDInput := ""

		Convey("When a new StreamID is created", func() {
			newStreamIDWithEmptyInput := func() {
				es.NewStreamID(streamIDInput)
			}

			Convey("It should fail with a panic", func() {
				So(newStreamIDWithEmptyInput, ShouldPanic)
			})
		})
	})
}

func TestStreamID_String(t *testing.T) {
	Convey("Given a StreamID", t, func() {
		streamIDInput := "customer-123"
		streamID := es.NewStreamID(streamIDInput)

		Convey("It should expose the expected value", func() {
			So(streamID.String(), ShouldEqual, streamIDInput)
		})
	})
}
