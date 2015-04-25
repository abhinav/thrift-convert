package main

import (
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
)

// Converter manages context for converting from data from one thrift protocol
// to another.
type Converter struct {
	from, to thrift.TProtocol
}

func (c Converter) convertMessage() error {
	name, fieldType, seqID, err := c.from.ReadMessageBegin()
	if err != nil {
		return err
	}

	if err := c.to.WriteMessageBegin(name, fieldType, seqID); err != nil {
		return err
	}

	if err := c.convertFields(); err != nil {
		return err
	}

	if err := c.from.ReadMessageEnd(); err != nil {
		return err
	}

	if err := c.to.WriteMessageEnd(); err != nil {
		return err
	}

	if err := c.to.Flush(); err != nil {
		return err
	}

	return nil
}

func (c Converter) convertFields() error {
	for {
		fieldType, err := c.convertField()
		if err != nil {
			return err
		}

		if fieldType == thrift.STOP {
			return c.to.WriteFieldStop()
		}
	}
}

func (c Converter) convertField() (thrift.TType, error) {
	name, fieldType, fieldID, err := c.from.ReadFieldBegin()
	if err != nil {
		return fieldType, err
	}

	if fieldType == thrift.STOP {
		return fieldType, nil
	}

	if err := c.to.WriteFieldBegin(name, fieldType, fieldID); err != nil {
		return fieldType, err
	}

	if err := c.convertType(fieldType); err != nil {
		return fieldType, err
	}

	if err := c.from.ReadFieldEnd(); err != nil {
		return fieldType, err
	}

	if err := c.to.WriteFieldEnd(); err != nil {
		return fieldType, err
	}

	return fieldType, nil
}

func (c Converter) convertType(fieldType thrift.TType) error {
	switch fieldType {
	case thrift.STOP:
		panic("convertType must not be called with STOP fields")
	case thrift.BOOL:
		value, err := c.from.ReadBool()
		if err != nil {
			return err
		}
		return c.to.WriteBool(value)
	case thrift.BYTE:
		value, err := c.from.ReadByte()
		if err != nil {
			return err
		}
		return c.to.WriteByte(value)
	case thrift.I16:
		value, err := c.from.ReadI16()
		if err != nil {
			return err
		}
		return c.to.WriteI16(value)
	case thrift.I32:
		value, err := c.from.ReadI32()
		if err != nil {
			return err
		}
		return c.to.WriteI32(value)
	case thrift.I64:
		value, err := c.from.ReadI64()
		if err != nil {
			return err
		}
		return c.to.WriteI64(value)
	case thrift.DOUBLE:
		value, err := c.from.ReadDouble()
		if err != nil {
			return err
		}
		return c.to.WriteDouble(value)
	case thrift.STRING:
		value, err := c.from.ReadString()
		if err != nil {
			return err
		}
		return c.to.WriteString(value)
	case thrift.STRUCT:
		name, err := c.from.ReadStructBegin()
		if err != nil {
			return err
		}

		if err := c.to.WriteStructBegin(name); err != nil {
			return err
		}

		if err := c.convertFields(); err != nil {
			return err
		}

		if err := c.from.ReadStructEnd(); err != nil {
			return err
		}

		if err := c.to.WriteStructEnd(); err != nil {
			return err
		}
	case thrift.MAP:
		keyType, valueType, size, err := c.from.ReadMapBegin()
		if err != nil {
			return err
		}

		if err := c.to.WriteMapBegin(keyType, valueType, size); err != nil {
			return err
		}

		for i := 0; i < size; i++ {
			if err := c.convertType(keyType); err != nil {
				return err
			}

			if err := c.convertType(valueType); err != nil {
				return err
			}
		}

		if err := c.from.ReadMapEnd(); err != nil {
			return err
		}

		if err := c.to.WriteMapEnd(); err != nil {
			return err
		}
	case thrift.SET:
		elemType, size, err := c.from.ReadSetBegin()
		if err != nil {
			return err
		}

		if err := c.to.WriteSetBegin(elemType, size); err != nil {
			return err
		}

		for i := 0; i < size; i++ {
			if err := c.convertType(elemType); err != nil {
				return err
			}
		}

		if err := c.from.ReadSetEnd(); err != nil {
			return err
		}

		if err := c.to.WriteSetEnd(); err != nil {
			return err
		}
	case thrift.LIST:
		elemType, size, err := c.from.ReadListBegin()
		if err != nil {
			return err
		}

		if err := c.to.WriteListBegin(elemType, size); err != nil {
			return err
		}

		for i := 0; i < size; i++ {
			if err := c.convertType(elemType); err != nil {
				return err
			}
		}

		if err := c.from.ReadListEnd(); err != nil {
			return err
		}

		if err := c.to.WriteListEnd(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported field type '%q'", fieldType)
	}

	return nil
}
