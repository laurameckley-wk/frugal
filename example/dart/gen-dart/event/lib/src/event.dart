/**
 * Autogenerated by Thrift Compiler (1.0.0-dev)
 *
 * DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
 *  @generated
 */
library event.src.event;

import 'dart:typed_data' show Uint8List;
import 'package:thrift/thrift.dart';
import 'package:event/event.dart';

/// This docstring gets added to the generated code because it has
/// the @ sign.
class Event implements TBase {
  static final TStruct _STRUCT_DESC = new TStruct("Event");
  static final TField _ID_FIELD_DESC = new TField("ID", TType.I64, 1);
  static final TField _MESSAGE_FIELD_DESC = new TField("Message", TType.STRING, 2);

  /// ID is a unique identifier for an event.
  int _iD;
  static const int ID = 1;
  /// Message contains the event payload.
  String _message;
  static const int MESSAGE = 2;

  bool __isset_iD = false;

  Event() {
    this.iD = -1;

  }

  // iD
  /// ID is a unique identifier for an event.
  int get iD => this._iD;

  /// ID is a unique identifier for an event.
  set iD(int iD) {
    this._iD = iD;
    this.__isset_iD = true;
  }

  bool isSetID() => this.__isset_iD;

  unsetID() {
    this.__isset_iD = false;
  }

  // message
  /// Message contains the event payload.
  String get message => this._message;

  /// Message contains the event payload.
  set message(String message) {
    this._message = message;
  }

  bool isSetMessage() => this.message != null;

  unsetMessage() {
    this.message = null;
  }

  getFieldValue(int fieldID) {
    switch (fieldID) {
      case ID:
        return this.iD;
      case MESSAGE:
        return this.message;
      default:
        throw new ArgumentError("Field $fieldID doesn't exist!");
    }
  }

  setFieldValue(int fieldID, Object value) {
    switch (fieldID) {
      case ID:
        if (value == null) {
          unsetID();
        } else {
          this.iD = value;
        }
        break;

      case MESSAGE:
        if (value == null) {
          unsetMessage();
        } else {
          this.message = value;
        }
        break;

      default:
        throw new ArgumentError("Field $fieldID doesn't exist!");
    }
  }

  // Returns true if field corresponding to fieldID is set (has been assigned a value) and false otherwise
  bool isSet(int fieldID) {
    switch (fieldID) {
      case ID:
        return isSetID();
      case MESSAGE:
        return isSetMessage();
      default:
        throw new ArgumentError("Field $fieldID doesn't exist!");
    }
  }

  read(TProtocol iprot) {
    TField field;
    iprot.readStructBegin();
    while (true) {
      field = iprot.readFieldBegin();
      if (field.type == TType.STOP) {
        break;
      }
      switch (field.id) {
        case ID:
          if (field.type == TType.I64) {
            this.iD = iprot.readI64();
            this.__isset_iD = true;
          } else {
            TProtocolUtil.skip(iprot, field.type);
          }
          break;
        case MESSAGE:
          if (field.type == TType.STRING) {
            this.message = iprot.readString();
          } else {
            TProtocolUtil.skip(iprot, field.type);
          }
          break;
        default:
          TProtocolUtil.skip(iprot, field.type);
          break;
      }
      iprot.readFieldEnd();
    }
    iprot.readStructEnd();

    // check for required fields of primitive type, which can't be checked in the validate method
    if (!__isset_iD) {
        throw new TProtocolError(TProtocolErrorType.UNKNOWN, "Required field 'iD' was not found in serialized data! Struct: " + toString());
    }

    validate();
  }

  write(TProtocol oprot) {
    validate();

    oprot.writeStructBegin(_STRUCT_DESC);
    oprot.writeFieldBegin(_ID_FIELD_DESC);
    oprot.writeI64(this.iD);
    oprot.writeFieldEnd();
    if (this.message != null) {
      oprot.writeFieldBegin(_MESSAGE_FIELD_DESC);
      oprot.writeString(this.message);
      oprot.writeFieldEnd();
    }
    oprot.writeFieldStop();
    oprot.writeStructEnd();
  }

  String toString() {
    StringBuffer ret = new StringBuffer("Event(");

    ret.write("iD:");
    ret.write(this.iD);

    ret.write(", ");
    ret.write("message:");
    if (this.message == null) {
      ret.write("null");
    } else {
      ret.write(this.message);
    }

    ret.write(")");

    return ret.toString();
  }

  validate() {
    // check for required fields
    // alas, we cannot check 'iD' because it's a primitive and you chose the non-beans generator.
    if (message == null) {
      throw new TProtocolError(TProtocolErrorType.UNKNOWN, "Required field 'message' was not present! Struct: " + toString());
    }
    // check that fields of type enum have valid values
  }

}

