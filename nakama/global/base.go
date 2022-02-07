package global

import "google.golang.org/protobuf/encoding/protojson"

var (
	Marshaler   *protojson.MarshalOptions
	Unmarshaler *protojson.UnmarshalOptions
)

func Init() {
	Marshaler = &protojson.MarshalOptions{
		UseEnumNumbers:  true,
		EmitUnpopulated: true,
	}
	Unmarshaler = &protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}
}
