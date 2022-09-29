package postgres

import (
	"database/sql/driver"
	"github.com/jackc/pgtype"
)

type timestamptz struct {
	Time pgtype.Timestamptz
}

func (dst *timestamptz) Set(src interface{}) error {

	if src == nil {
		*dst = timestamptz{pgtype.Timestamptz{Status: pgtype.Null}}
		return nil
	}

	err := dst.Time.Set(src)
	if err != nil {
		return err
	}

	return nil
}

func (dst *timestamptz) Get() interface{} {
	return dst.Time.Get()
}

func (src *timestamptz) AssignTo(dst interface{}) error {
	err := src.Time.AssignTo(dst)
	if err != nil {
		return err
	}

	return nil
}

func (dst *timestamptz) DecodeText(ci *pgtype.ConnInfo, src []byte) error {
	err := dst.Time.DecodeText(ci, src)
	if err != nil {
		return err
	}

	return nil
}

func (dst *timestamptz) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	var t pgtype.Timestamptz
	err := t.DecodeBinary(ci, src)
	if err != nil {
		return err
	}

	*dst = timestamptz{
		Time: t,
	}
	return nil
}

func (src timestamptz) EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	text, err := src.Time.EncodeText(ci, buf)
	if err != nil {
		return nil, err
	}

	return text, nil
}

func (src timestamptz) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	binary, err := src.Time.EncodeBinary(ci, buf)
	if err != nil {
		return nil, err
	}

	return binary, nil
}

func (dst *timestamptz) Scan(src interface{}) error {
	err := dst.Time.Scan(src)
	if err != nil {
		return err
	}

	return nil
}

func (src timestamptz) Value() (driver.Value, error) {
	return src.Time.Value()
}
