# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: outliers.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='outliers.proto',
  package='pb',
  syntax='proto3',
  serialized_options=b'Z!github.com/Deiklov/diplom_backend',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\x0eoutliers.proto\x12\x02pb\x1a\x1fgoogle/protobuf/timestamp.proto\"O\n\x06Metric\x12(\n\x04time\x18\x01 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\x0c\n\x04name\x18\x02 \x01(\t\x12\r\n\x05value\x18\x03 \x01(\x01\".\n\x0fOutliersRequest\x12\x1b\n\x07metrics\x18\x01 \x03(\x0b\x32\n.pb.Metric\"#\n\x10OutliersResponse\x12\x0f\n\x07indices\x18\x01 \x03(\x05\x32\x41\n\x08Outliers\x12\x35\n\x06\x44\x65tect\x12\x13.pb.OutliersRequest\x1a\x14.pb.OutliersResponse\"\x00\x42#Z!github.com/Deiklov/diplom_backendb\x06proto3'
  ,
  dependencies=[google_dot_protobuf_dot_timestamp__pb2.DESCRIPTOR,])




_METRIC = _descriptor.Descriptor(
  name='Metric',
  full_name='pb.Metric',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='time', full_name='pb.Metric.time', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='name', full_name='pb.Metric.name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='pb.Metric.value', index=2,
      number=3, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=55,
  serialized_end=134,
)


_OUTLIERSREQUEST = _descriptor.Descriptor(
  name='OutliersRequest',
  full_name='pb.OutliersRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='metrics', full_name='pb.OutliersRequest.metrics', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=136,
  serialized_end=182,
)


_OUTLIERSRESPONSE = _descriptor.Descriptor(
  name='OutliersResponse',
  full_name='pb.OutliersResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='indices', full_name='pb.OutliersResponse.indices', index=0,
      number=1, type=5, cpp_type=1, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=184,
  serialized_end=219,
)

_METRIC.fields_by_name['time'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_OUTLIERSREQUEST.fields_by_name['metrics'].message_type = _METRIC
DESCRIPTOR.message_types_by_name['Metric'] = _METRIC
DESCRIPTOR.message_types_by_name['OutliersRequest'] = _OUTLIERSREQUEST
DESCRIPTOR.message_types_by_name['OutliersResponse'] = _OUTLIERSRESPONSE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Metric = _reflection.GeneratedProtocolMessageType('Metric', (_message.Message,), {
  'DESCRIPTOR' : _METRIC,
  '__module__' : 'outliers_pb2'
  # @@protoc_insertion_point(class_scope:pb.Metric)
  })
_sym_db.RegisterMessage(Metric)

OutliersRequest = _reflection.GeneratedProtocolMessageType('OutliersRequest', (_message.Message,), {
  'DESCRIPTOR' : _OUTLIERSREQUEST,
  '__module__' : 'outliers_pb2'
  # @@protoc_insertion_point(class_scope:pb.OutliersRequest)
  })
_sym_db.RegisterMessage(OutliersRequest)

OutliersResponse = _reflection.GeneratedProtocolMessageType('OutliersResponse', (_message.Message,), {
  'DESCRIPTOR' : _OUTLIERSRESPONSE,
  '__module__' : 'outliers_pb2'
  # @@protoc_insertion_point(class_scope:pb.OutliersResponse)
  })
_sym_db.RegisterMessage(OutliersResponse)


DESCRIPTOR._options = None

_OUTLIERS = _descriptor.ServiceDescriptor(
  name='Outliers',
  full_name='pb.Outliers',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=221,
  serialized_end=286,
  methods=[
  _descriptor.MethodDescriptor(
    name='Detect',
    full_name='pb.Outliers.Detect',
    index=0,
    containing_service=None,
    input_type=_OUTLIERSREQUEST,
    output_type=_OUTLIERSRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_OUTLIERS)

DESCRIPTOR.services_by_name['Outliers'] = _OUTLIERS

# @@protoc_insertion_point(module_scope)
