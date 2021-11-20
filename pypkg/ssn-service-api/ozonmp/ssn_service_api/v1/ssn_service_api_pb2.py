# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: ozonmp/ssn_service_api/v1/ssn_service_api.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from validate import validate_pb2 as validate_dot_validate__pb2
from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='ozonmp/ssn_service_api/v1/ssn_service_api.proto',
  package='ozonmp.ssn_service_api.v1',
  syntax='proto3',
  serialized_options=_b('ZEgithub.com/ozonmp/ssn-service-api/pkg/ssn-service-api;ssn_service_api'),
  serialized_pb=_b('\n/ozonmp/ssn_service_api/v1/ssn_service_api.proto\x12\x19ozonmp.ssn_service_api.v1\x1a\x17validate/validate.proto\x1a\x1cgoogle/api/annotations.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"\xc5\x01\n\x07Service\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\x12\x12\n\x04name\x18\x02 \x01(\tR\x04name\x12 \n\x0b\x64\x65scription\x18\x03 \x01(\tR\x0b\x64\x65scription\x12\x39\n\ncreated_at\x18\x04 \x01(\x0b\x32\x1a.google.protobuf.TimestampR\tcreatedAt\x12\x39\n\nupdated_at\x18\x05 \x01(\x0b\x32\x1a.google.protobuf.TimestampR\tupdatedAt\"e\n\x16\x43reateServiceV1Request\x12\x1d\n\x04name\x18\x01 \x01(\tB\t\xfa\x42\x06r\x04\x10\x01\x18\x64R\x04name\x12,\n\x0b\x64\x65scription\x18\x02 \x01(\tB\n\xfa\x42\x07r\x05\x10\x01\x18\xc8\x01R\x0b\x64\x65scription\"8\n\x17\x43reateServiceV1Response\x12\x1d\n\nservice_id\x18\x01 \x01(\x04R\tserviceId\"B\n\x18\x44\x65scribeServiceV1Request\x12&\n\nservice_id\x18\x01 \x01(\x04\x42\x07\xfa\x42\x04\x32\x02 \x00R\tserviceId\"Y\n\x19\x44\x65scribeServiceV1Response\x12<\n\x07service\x18\x01 \x01(\x0b\x32\".ozonmp.ssn_service_api.v1.ServiceR\x07service\"\x8d\x01\n\x16UpdateServiceV1Request\x12&\n\nservice_id\x18\x01 \x01(\x04\x42\x07\xfa\x42\x04\x32\x02 \x00R\tserviceId\x12\x1d\n\x04name\x18\x02 \x01(\tB\t\xfa\x42\x06r\x04\x10\x01\x18\x64R\x04name\x12,\n\x0b\x64\x65scription\x18\x03 \x01(\tB\n\xfa\x42\x07r\x05\x10\x01\x18\xc8\x01R\x0b\x64\x65scription\"\x19\n\x17UpdateServiceV1Response\"Q\n\x15ListServicesV1Request\x12\x16\n\x06offset\x18\x01 \x01(\x04R\x06offset\x12 \n\x05limit\x18\x02 \x01(\x04\x42\n\xfa\x42\x07\x32\x05\x18\xf4\x03 \x00R\x05limit\"X\n\x16ListServicesV1Response\x12>\n\x08services\x18\x01 \x03(\x0b\x32\".ozonmp.ssn_service_api.v1.ServiceR\x08services\"@\n\x16RemoveServiceV1Request\x12&\n\nservice_id\x18\x01 \x01(\x04\x42\x07\xfa\x42\x04\x32\x02 \x00R\tserviceId\"\x19\n\x17RemoveServiceV1Response\"\x8a\x01\n\x13ServiceEventPayload\x12&\n\nservice_id\x18\x01 \x01(\x04\x42\x07\xfa\x42\x04\x32\x02 \x00R\tserviceId\x12\x1d\n\x04name\x18\x02 \x01(\tB\t\xfa\x42\x06r\x04\x10\x01\x18\x64R\x04name\x12,\n\x0b\x64\x65scription\x18\x03 \x01(\tB\n\xfa\x42\x07r\x05\x10\x01\x18\xc8\x01R\x0b\x64\x65scription\"\xd9\x01\n\x0cServiceEvent\x12\x17\n\x02id\x18\x01 \x01(\x04\x42\x07\xfa\x42\x04\x32\x02 \x00R\x02id\x12&\n\nservice_id\x18\x02 \x01(\x04\x42\x07\xfa\x42\x04\x32\x02 \x00R\tserviceId\x12\x1b\n\x04type\x18\x03 \x01(\tB\x07\xfa\x42\x04r\x02\x10\x01R\x04type\x12!\n\x07subtype\x18\x04 \x01(\tB\x07\xfa\x42\x04r\x02\x10\x01R\x07subtype\x12H\n\x07payload\x18\x05 \x01(\x0b\x32..ozonmp.ssn_service_api.v1.ServiceEventPayloadR\x07payload2\xa2\x06\n\x14SsnServiceApiService\x12\x95\x01\n\x0f\x43reateServiceV1\x12\x31.ozonmp.ssn_service_api.v1.CreateServiceV1Request\x1a\x32.ozonmp.ssn_service_api.v1.CreateServiceV1Response\"\x1b\x82\xd3\xe4\x93\x02\x15\"\x10/api/v1/services:\x01*\x12\xa5\x01\n\x11\x44\x65scribeServiceV1\x12\x33.ozonmp.ssn_service_api.v1.DescribeServiceV1Request\x1a\x34.ozonmp.ssn_service_api.v1.DescribeServiceV1Response\"%\x82\xd3\xe4\x93\x02\x1f\x12\x1d/api/v1/services/{service_id}\x12\x95\x01\n\x0fUpdateServiceV1\x12\x31.ozonmp.ssn_service_api.v1.UpdateServiceV1Request\x1a\x32.ozonmp.ssn_service_api.v1.UpdateServiceV1Response\"\x1b\x82\xd3\xe4\x93\x02\x15\x1a\x10/api/v1/services:\x01*\x12\x8f\x01\n\x0eListServicesV1\x12\x30.ozonmp.ssn_service_api.v1.ListServicesV1Request\x1a\x31.ozonmp.ssn_service_api.v1.ListServicesV1Response\"\x18\x82\xd3\xe4\x93\x02\x12\x12\x10/api/v1/services\x12\x9f\x01\n\x0fRemoveServiceV1\x12\x31.ozonmp.ssn_service_api.v1.RemoveServiceV1Request\x1a\x32.ozonmp.ssn_service_api.v1.RemoveServiceV1Response\"%\x82\xd3\xe4\x93\x02\x1f*\x1d/api/v1/services/{service_id}BGZEgithub.com/ozonmp/ssn-service-api/pkg/ssn-service-api;ssn_service_apib\x06proto3')
  ,
  dependencies=[validate_dot_validate__pb2.DESCRIPTOR,google_dot_api_dot_annotations__pb2.DESCRIPTOR,google_dot_protobuf_dot_timestamp__pb2.DESCRIPTOR,])




_SERVICE = _descriptor.Descriptor(
  name='Service',
  full_name='ozonmp.ssn_service_api.v1.Service',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='ozonmp.ssn_service_api.v1.Service.id', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, json_name='id', file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='name', full_name='ozonmp.ssn_service_api.v1.Service.name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, json_name='name', file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='description', full_name='ozonmp.ssn_service_api.v1.Service.description', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, json_name='description', file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='created_at', full_name='ozonmp.ssn_service_api.v1.Service.created_at', index=3,
      number=4, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, json_name='createdAt', file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='updated_at', full_name='ozonmp.ssn_service_api.v1.Service.updated_at', index=4,
      number=5, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, json_name='updatedAt', file=DESCRIPTOR),
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
  serialized_start=167,
  serialized_end=364,
)


_CREATESERVICEV1REQUEST = _descriptor.Descriptor(
  name='CreateServiceV1Request',
  full_name='ozonmp.ssn_service_api.v1.CreateServiceV1Request',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='ozonmp.ssn_service_api.v1.CreateServiceV1Request.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\006r\004\020\001\030d'), json_name='name', file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='description', full_name='ozonmp.ssn_service_api.v1.CreateServiceV1Request.description', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\007r\005\020\001\030\310\001'), json_name='description', file=DESCRIPTOR),
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
  serialized_start=366,
  serialized_end=467,
)


_CREATESERVICEV1RESPONSE = _descriptor.Descriptor(
  name='CreateServiceV1Response',
  full_name='ozonmp.ssn_service_api.v1.CreateServiceV1Response',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='service_id', full_name='ozonmp.ssn_service_api.v1.CreateServiceV1Response.service_id', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, json_name='serviceId', file=DESCRIPTOR),
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
  serialized_start=469,
  serialized_end=525,
)


_DESCRIBESERVICEV1REQUEST = _descriptor.Descriptor(
  name='DescribeServiceV1Request',
  full_name='ozonmp.ssn_service_api.v1.DescribeServiceV1Request',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='service_id', full_name='ozonmp.ssn_service_api.v1.DescribeServiceV1Request.service_id', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\0042\002 \000'), json_name='serviceId', file=DESCRIPTOR),
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
  serialized_start=527,
  serialized_end=593,
)


_DESCRIBESERVICEV1RESPONSE = _descriptor.Descriptor(
  name='DescribeServiceV1Response',
  full_name='ozonmp.ssn_service_api.v1.DescribeServiceV1Response',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='service', full_name='ozonmp.ssn_service_api.v1.DescribeServiceV1Response.service', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, json_name='service', file=DESCRIPTOR),
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
  serialized_start=595,
  serialized_end=684,
)


_UPDATESERVICEV1REQUEST = _descriptor.Descriptor(
  name='UpdateServiceV1Request',
  full_name='ozonmp.ssn_service_api.v1.UpdateServiceV1Request',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='service_id', full_name='ozonmp.ssn_service_api.v1.UpdateServiceV1Request.service_id', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\0042\002 \000'), json_name='serviceId', file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='name', full_name='ozonmp.ssn_service_api.v1.UpdateServiceV1Request.name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\006r\004\020\001\030d'), json_name='name', file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='description', full_name='ozonmp.ssn_service_api.v1.UpdateServiceV1Request.description', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\007r\005\020\001\030\310\001'), json_name='description', file=DESCRIPTOR),
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
  serialized_start=687,
  serialized_end=828,
)


_UPDATESERVICEV1RESPONSE = _descriptor.Descriptor(
  name='UpdateServiceV1Response',
  full_name='ozonmp.ssn_service_api.v1.UpdateServiceV1Response',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
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
  serialized_start=830,
  serialized_end=855,
)


_LISTSERVICESV1REQUEST = _descriptor.Descriptor(
  name='ListServicesV1Request',
  full_name='ozonmp.ssn_service_api.v1.ListServicesV1Request',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='offset', full_name='ozonmp.ssn_service_api.v1.ListServicesV1Request.offset', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, json_name='offset', file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='limit', full_name='ozonmp.ssn_service_api.v1.ListServicesV1Request.limit', index=1,
      number=2, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\0072\005\030\364\003 \000'), json_name='limit', file=DESCRIPTOR),
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
  serialized_start=857,
  serialized_end=938,
)


_LISTSERVICESV1RESPONSE = _descriptor.Descriptor(
  name='ListServicesV1Response',
  full_name='ozonmp.ssn_service_api.v1.ListServicesV1Response',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='services', full_name='ozonmp.ssn_service_api.v1.ListServicesV1Response.services', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, json_name='services', file=DESCRIPTOR),
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
  serialized_start=940,
  serialized_end=1028,
)


_REMOVESERVICEV1REQUEST = _descriptor.Descriptor(
  name='RemoveServiceV1Request',
  full_name='ozonmp.ssn_service_api.v1.RemoveServiceV1Request',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='service_id', full_name='ozonmp.ssn_service_api.v1.RemoveServiceV1Request.service_id', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\0042\002 \000'), json_name='serviceId', file=DESCRIPTOR),
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
  serialized_start=1030,
  serialized_end=1094,
)


_REMOVESERVICEV1RESPONSE = _descriptor.Descriptor(
  name='RemoveServiceV1Response',
  full_name='ozonmp.ssn_service_api.v1.RemoveServiceV1Response',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
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
  serialized_start=1096,
  serialized_end=1121,
)


_SERVICEEVENTPAYLOAD = _descriptor.Descriptor(
  name='ServiceEventPayload',
  full_name='ozonmp.ssn_service_api.v1.ServiceEventPayload',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='service_id', full_name='ozonmp.ssn_service_api.v1.ServiceEventPayload.service_id', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\0042\002 \000'), json_name='serviceId', file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='name', full_name='ozonmp.ssn_service_api.v1.ServiceEventPayload.name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\006r\004\020\001\030d'), json_name='name', file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='description', full_name='ozonmp.ssn_service_api.v1.ServiceEventPayload.description', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\007r\005\020\001\030\310\001'), json_name='description', file=DESCRIPTOR),
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
  serialized_start=1124,
  serialized_end=1262,
)


_SERVICEEVENT = _descriptor.Descriptor(
  name='ServiceEvent',
  full_name='ozonmp.ssn_service_api.v1.ServiceEvent',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='ozonmp.ssn_service_api.v1.ServiceEvent.id', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\0042\002 \000'), json_name='id', file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='service_id', full_name='ozonmp.ssn_service_api.v1.ServiceEvent.service_id', index=1,
      number=2, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\0042\002 \000'), json_name='serviceId', file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='type', full_name='ozonmp.ssn_service_api.v1.ServiceEvent.type', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\004r\002\020\001'), json_name='type', file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='subtype', full_name='ozonmp.ssn_service_api.v1.ServiceEvent.subtype', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\372B\004r\002\020\001'), json_name='subtype', file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='payload', full_name='ozonmp.ssn_service_api.v1.ServiceEvent.payload', index=4,
      number=5, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, json_name='payload', file=DESCRIPTOR),
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
  serialized_start=1265,
  serialized_end=1482,
)

_SERVICE.fields_by_name['created_at'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_SERVICE.fields_by_name['updated_at'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_DESCRIBESERVICEV1RESPONSE.fields_by_name['service'].message_type = _SERVICE
_LISTSERVICESV1RESPONSE.fields_by_name['services'].message_type = _SERVICE
_SERVICEEVENT.fields_by_name['payload'].message_type = _SERVICEEVENTPAYLOAD
DESCRIPTOR.message_types_by_name['Service'] = _SERVICE
DESCRIPTOR.message_types_by_name['CreateServiceV1Request'] = _CREATESERVICEV1REQUEST
DESCRIPTOR.message_types_by_name['CreateServiceV1Response'] = _CREATESERVICEV1RESPONSE
DESCRIPTOR.message_types_by_name['DescribeServiceV1Request'] = _DESCRIBESERVICEV1REQUEST
DESCRIPTOR.message_types_by_name['DescribeServiceV1Response'] = _DESCRIBESERVICEV1RESPONSE
DESCRIPTOR.message_types_by_name['UpdateServiceV1Request'] = _UPDATESERVICEV1REQUEST
DESCRIPTOR.message_types_by_name['UpdateServiceV1Response'] = _UPDATESERVICEV1RESPONSE
DESCRIPTOR.message_types_by_name['ListServicesV1Request'] = _LISTSERVICESV1REQUEST
DESCRIPTOR.message_types_by_name['ListServicesV1Response'] = _LISTSERVICESV1RESPONSE
DESCRIPTOR.message_types_by_name['RemoveServiceV1Request'] = _REMOVESERVICEV1REQUEST
DESCRIPTOR.message_types_by_name['RemoveServiceV1Response'] = _REMOVESERVICEV1RESPONSE
DESCRIPTOR.message_types_by_name['ServiceEventPayload'] = _SERVICEEVENTPAYLOAD
DESCRIPTOR.message_types_by_name['ServiceEvent'] = _SERVICEEVENT
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Service = _reflection.GeneratedProtocolMessageType('Service', (_message.Message,), dict(
  DESCRIPTOR = _SERVICE,
  __module__ = 'ozonmp.ssn_service_api.v1.ssn_service_api_pb2'
  # @@protoc_insertion_point(class_scope:ozonmp.ssn_service_api.v1.Service)
  ))
_sym_db.RegisterMessage(Service)

CreateServiceV1Request = _reflection.GeneratedProtocolMessageType('CreateServiceV1Request', (_message.Message,), dict(
  DESCRIPTOR = _CREATESERVICEV1REQUEST,
  __module__ = 'ozonmp.ssn_service_api.v1.ssn_service_api_pb2'
  # @@protoc_insertion_point(class_scope:ozonmp.ssn_service_api.v1.CreateServiceV1Request)
  ))
_sym_db.RegisterMessage(CreateServiceV1Request)

CreateServiceV1Response = _reflection.GeneratedProtocolMessageType('CreateServiceV1Response', (_message.Message,), dict(
  DESCRIPTOR = _CREATESERVICEV1RESPONSE,
  __module__ = 'ozonmp.ssn_service_api.v1.ssn_service_api_pb2'
  # @@protoc_insertion_point(class_scope:ozonmp.ssn_service_api.v1.CreateServiceV1Response)
  ))
_sym_db.RegisterMessage(CreateServiceV1Response)

DescribeServiceV1Request = _reflection.GeneratedProtocolMessageType('DescribeServiceV1Request', (_message.Message,), dict(
  DESCRIPTOR = _DESCRIBESERVICEV1REQUEST,
  __module__ = 'ozonmp.ssn_service_api.v1.ssn_service_api_pb2'
  # @@protoc_insertion_point(class_scope:ozonmp.ssn_service_api.v1.DescribeServiceV1Request)
  ))
_sym_db.RegisterMessage(DescribeServiceV1Request)

DescribeServiceV1Response = _reflection.GeneratedProtocolMessageType('DescribeServiceV1Response', (_message.Message,), dict(
  DESCRIPTOR = _DESCRIBESERVICEV1RESPONSE,
  __module__ = 'ozonmp.ssn_service_api.v1.ssn_service_api_pb2'
  # @@protoc_insertion_point(class_scope:ozonmp.ssn_service_api.v1.DescribeServiceV1Response)
  ))
_sym_db.RegisterMessage(DescribeServiceV1Response)

UpdateServiceV1Request = _reflection.GeneratedProtocolMessageType('UpdateServiceV1Request', (_message.Message,), dict(
  DESCRIPTOR = _UPDATESERVICEV1REQUEST,
  __module__ = 'ozonmp.ssn_service_api.v1.ssn_service_api_pb2'
  # @@protoc_insertion_point(class_scope:ozonmp.ssn_service_api.v1.UpdateServiceV1Request)
  ))
_sym_db.RegisterMessage(UpdateServiceV1Request)

UpdateServiceV1Response = _reflection.GeneratedProtocolMessageType('UpdateServiceV1Response', (_message.Message,), dict(
  DESCRIPTOR = _UPDATESERVICEV1RESPONSE,
  __module__ = 'ozonmp.ssn_service_api.v1.ssn_service_api_pb2'
  # @@protoc_insertion_point(class_scope:ozonmp.ssn_service_api.v1.UpdateServiceV1Response)
  ))
_sym_db.RegisterMessage(UpdateServiceV1Response)

ListServicesV1Request = _reflection.GeneratedProtocolMessageType('ListServicesV1Request', (_message.Message,), dict(
  DESCRIPTOR = _LISTSERVICESV1REQUEST,
  __module__ = 'ozonmp.ssn_service_api.v1.ssn_service_api_pb2'
  # @@protoc_insertion_point(class_scope:ozonmp.ssn_service_api.v1.ListServicesV1Request)
  ))
_sym_db.RegisterMessage(ListServicesV1Request)

ListServicesV1Response = _reflection.GeneratedProtocolMessageType('ListServicesV1Response', (_message.Message,), dict(
  DESCRIPTOR = _LISTSERVICESV1RESPONSE,
  __module__ = 'ozonmp.ssn_service_api.v1.ssn_service_api_pb2'
  # @@protoc_insertion_point(class_scope:ozonmp.ssn_service_api.v1.ListServicesV1Response)
  ))
_sym_db.RegisterMessage(ListServicesV1Response)

RemoveServiceV1Request = _reflection.GeneratedProtocolMessageType('RemoveServiceV1Request', (_message.Message,), dict(
  DESCRIPTOR = _REMOVESERVICEV1REQUEST,
  __module__ = 'ozonmp.ssn_service_api.v1.ssn_service_api_pb2'
  # @@protoc_insertion_point(class_scope:ozonmp.ssn_service_api.v1.RemoveServiceV1Request)
  ))
_sym_db.RegisterMessage(RemoveServiceV1Request)

RemoveServiceV1Response = _reflection.GeneratedProtocolMessageType('RemoveServiceV1Response', (_message.Message,), dict(
  DESCRIPTOR = _REMOVESERVICEV1RESPONSE,
  __module__ = 'ozonmp.ssn_service_api.v1.ssn_service_api_pb2'
  # @@protoc_insertion_point(class_scope:ozonmp.ssn_service_api.v1.RemoveServiceV1Response)
  ))
_sym_db.RegisterMessage(RemoveServiceV1Response)

ServiceEventPayload = _reflection.GeneratedProtocolMessageType('ServiceEventPayload', (_message.Message,), dict(
  DESCRIPTOR = _SERVICEEVENTPAYLOAD,
  __module__ = 'ozonmp.ssn_service_api.v1.ssn_service_api_pb2'
  # @@protoc_insertion_point(class_scope:ozonmp.ssn_service_api.v1.ServiceEventPayload)
  ))
_sym_db.RegisterMessage(ServiceEventPayload)

ServiceEvent = _reflection.GeneratedProtocolMessageType('ServiceEvent', (_message.Message,), dict(
  DESCRIPTOR = _SERVICEEVENT,
  __module__ = 'ozonmp.ssn_service_api.v1.ssn_service_api_pb2'
  # @@protoc_insertion_point(class_scope:ozonmp.ssn_service_api.v1.ServiceEvent)
  ))
_sym_db.RegisterMessage(ServiceEvent)


DESCRIPTOR._options = None
_CREATESERVICEV1REQUEST.fields_by_name['name']._options = None
_CREATESERVICEV1REQUEST.fields_by_name['description']._options = None
_DESCRIBESERVICEV1REQUEST.fields_by_name['service_id']._options = None
_UPDATESERVICEV1REQUEST.fields_by_name['service_id']._options = None
_UPDATESERVICEV1REQUEST.fields_by_name['name']._options = None
_UPDATESERVICEV1REQUEST.fields_by_name['description']._options = None
_LISTSERVICESV1REQUEST.fields_by_name['limit']._options = None
_REMOVESERVICEV1REQUEST.fields_by_name['service_id']._options = None
_SERVICEEVENTPAYLOAD.fields_by_name['service_id']._options = None
_SERVICEEVENTPAYLOAD.fields_by_name['name']._options = None
_SERVICEEVENTPAYLOAD.fields_by_name['description']._options = None
_SERVICEEVENT.fields_by_name['id']._options = None
_SERVICEEVENT.fields_by_name['service_id']._options = None
_SERVICEEVENT.fields_by_name['type']._options = None
_SERVICEEVENT.fields_by_name['subtype']._options = None

_SSNSERVICEAPISERVICE = _descriptor.ServiceDescriptor(
  name='SsnServiceApiService',
  full_name='ozonmp.ssn_service_api.v1.SsnServiceApiService',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=1485,
  serialized_end=2287,
  methods=[
  _descriptor.MethodDescriptor(
    name='CreateServiceV1',
    full_name='ozonmp.ssn_service_api.v1.SsnServiceApiService.CreateServiceV1',
    index=0,
    containing_service=None,
    input_type=_CREATESERVICEV1REQUEST,
    output_type=_CREATESERVICEV1RESPONSE,
    serialized_options=_b('\202\323\344\223\002\025\"\020/api/v1/services:\001*'),
  ),
  _descriptor.MethodDescriptor(
    name='DescribeServiceV1',
    full_name='ozonmp.ssn_service_api.v1.SsnServiceApiService.DescribeServiceV1',
    index=1,
    containing_service=None,
    input_type=_DESCRIBESERVICEV1REQUEST,
    output_type=_DESCRIBESERVICEV1RESPONSE,
    serialized_options=_b('\202\323\344\223\002\037\022\035/api/v1/services/{service_id}'),
  ),
  _descriptor.MethodDescriptor(
    name='UpdateServiceV1',
    full_name='ozonmp.ssn_service_api.v1.SsnServiceApiService.UpdateServiceV1',
    index=2,
    containing_service=None,
    input_type=_UPDATESERVICEV1REQUEST,
    output_type=_UPDATESERVICEV1RESPONSE,
    serialized_options=_b('\202\323\344\223\002\025\032\020/api/v1/services:\001*'),
  ),
  _descriptor.MethodDescriptor(
    name='ListServicesV1',
    full_name='ozonmp.ssn_service_api.v1.SsnServiceApiService.ListServicesV1',
    index=3,
    containing_service=None,
    input_type=_LISTSERVICESV1REQUEST,
    output_type=_LISTSERVICESV1RESPONSE,
    serialized_options=_b('\202\323\344\223\002\022\022\020/api/v1/services'),
  ),
  _descriptor.MethodDescriptor(
    name='RemoveServiceV1',
    full_name='ozonmp.ssn_service_api.v1.SsnServiceApiService.RemoveServiceV1',
    index=4,
    containing_service=None,
    input_type=_REMOVESERVICEV1REQUEST,
    output_type=_REMOVESERVICEV1RESPONSE,
    serialized_options=_b('\202\323\344\223\002\037*\035/api/v1/services/{service_id}'),
  ),
])
_sym_db.RegisterServiceDescriptor(_SSNSERVICEAPISERVICE)

DESCRIPTOR.services_by_name['SsnServiceApiService'] = _SSNSERVICEAPISERVICE

# @@protoc_insertion_point(module_scope)
