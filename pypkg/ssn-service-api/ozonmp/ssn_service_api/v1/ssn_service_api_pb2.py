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


DESCRIPTOR = _descriptor.FileDescriptor(
  name='ozonmp/ssn_service_api/v1/ssn_service_api.proto',
  package='ozonmp.ssn_service_api.v1',
  syntax='proto3',
  serialized_options=_b('ZEgithub.com/ozonmp/ssn-service-api/pkg/ssn-service-api;ssn_service_api'),
  serialized_pb=_b('\n/ozonmp/ssn_service_api/v1/ssn_service_api.proto\x12\x19ozonmp.ssn_service_api.v1\x1a\x17validate/validate.proto\x1a\x1cgoogle/api/annotations.proto\"-\n\x07Service\x12\x0e\n\x02id\x18\x01 \x01(\x04R\x02id\x12\x12\n\x04name\x18\x02 \x01(\tR\x04name\"7\n\x16\x43reateServiceV1Request\x12\x1d\n\x04name\x18\x01 \x01(\tB\t\xfa\x42\x06r\x04\x10\x01\x18\x64R\x04name\"8\n\x17\x43reateServiceV1Response\x12\x1d\n\nservice_id\x18\x01 \x01(\x04R\tserviceId\"B\n\x18\x44\x65scribeServiceV1Request\x12&\n\nservice_id\x18\x01 \x01(\x04\x42\x07\xfa\x42\x04\x32\x02 \x00R\tserviceId\"Y\n\x19\x44\x65scribeServiceV1Response\x12<\n\x07service\x18\x01 \x01(\x0b\x32\".ozonmp.ssn_service_api.v1.ServiceR\x07service\"\x17\n\x15ListServicesV1Request\"X\n\x16ListServicesV1Response\x12>\n\x08services\x18\x01 \x03(\x0b\x32\".ozonmp.ssn_service_api.v1.ServiceR\x08services\"@\n\x16RemoveServiceV1Request\x12&\n\nservice_id\x18\x01 \x01(\x04\x42\x07\xfa\x42\x04\x32\x02 \x00R\tserviceId\"/\n\x17RemoveServiceV1Response\x12\x14\n\x05\x66ound\x18\x01 \x01(\x08R\x05\x66ound2\x8a\x05\n\x14SsnServiceApiService\x12\x95\x01\n\x0f\x43reateServiceV1\x12\x31.ozonmp.ssn_service_api.v1.CreateServiceV1Request\x1a\x32.ozonmp.ssn_service_api.v1.CreateServiceV1Response\"\x1b\x82\xd3\xe4\x93\x02\x15\"\x10/api/v1/services:\x01*\x12\xa5\x01\n\x11\x44\x65scribeServiceV1\x12\x33.ozonmp.ssn_service_api.v1.DescribeServiceV1Request\x1a\x34.ozonmp.ssn_service_api.v1.DescribeServiceV1Response\"%\x82\xd3\xe4\x93\x02\x1f\x12\x1d/api/v1/services/{service_id}\x12\x8f\x01\n\x0eListServicesV1\x12\x30.ozonmp.ssn_service_api.v1.ListServicesV1Request\x1a\x31.ozonmp.ssn_service_api.v1.ListServicesV1Response\"\x18\x82\xd3\xe4\x93\x02\x12\x12\x10/api/v1/services\x12\x9f\x01\n\x0fRemoveServiceV1\x12\x31.ozonmp.ssn_service_api.v1.RemoveServiceV1Request\x1a\x32.ozonmp.ssn_service_api.v1.RemoveServiceV1Response\"%\x82\xd3\xe4\x93\x02\x1f*\x1d/api/v1/services/{service_id}BGZEgithub.com/ozonmp/ssn-service-api/pkg/ssn-service-api;ssn_service_apib\x06proto3')
  ,
  dependencies=[validate_dot_validate__pb2.DESCRIPTOR,google_dot_api_dot_annotations__pb2.DESCRIPTOR,])




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
  serialized_start=133,
  serialized_end=178,
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
  serialized_start=180,
  serialized_end=235,
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
  serialized_start=237,
  serialized_end=293,
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
  serialized_start=295,
  serialized_end=361,
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
  serialized_start=363,
  serialized_end=452,
)


_LISTSERVICESV1REQUEST = _descriptor.Descriptor(
  name='ListServicesV1Request',
  full_name='ozonmp.ssn_service_api.v1.ListServicesV1Request',
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
  serialized_start=454,
  serialized_end=477,
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
  serialized_start=479,
  serialized_end=567,
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
  serialized_start=569,
  serialized_end=633,
)


_REMOVESERVICEV1RESPONSE = _descriptor.Descriptor(
  name='RemoveServiceV1Response',
  full_name='ozonmp.ssn_service_api.v1.RemoveServiceV1Response',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='found', full_name='ozonmp.ssn_service_api.v1.RemoveServiceV1Response.found', index=0,
      number=1, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, json_name='found', file=DESCRIPTOR),
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
  serialized_start=635,
  serialized_end=682,
)

_DESCRIBESERVICEV1RESPONSE.fields_by_name['service'].message_type = _SERVICE
_LISTSERVICESV1RESPONSE.fields_by_name['services'].message_type = _SERVICE
DESCRIPTOR.message_types_by_name['Service'] = _SERVICE
DESCRIPTOR.message_types_by_name['CreateServiceV1Request'] = _CREATESERVICEV1REQUEST
DESCRIPTOR.message_types_by_name['CreateServiceV1Response'] = _CREATESERVICEV1RESPONSE
DESCRIPTOR.message_types_by_name['DescribeServiceV1Request'] = _DESCRIBESERVICEV1REQUEST
DESCRIPTOR.message_types_by_name['DescribeServiceV1Response'] = _DESCRIBESERVICEV1RESPONSE
DESCRIPTOR.message_types_by_name['ListServicesV1Request'] = _LISTSERVICESV1REQUEST
DESCRIPTOR.message_types_by_name['ListServicesV1Response'] = _LISTSERVICESV1RESPONSE
DESCRIPTOR.message_types_by_name['RemoveServiceV1Request'] = _REMOVESERVICEV1REQUEST
DESCRIPTOR.message_types_by_name['RemoveServiceV1Response'] = _REMOVESERVICEV1RESPONSE
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


DESCRIPTOR._options = None
_CREATESERVICEV1REQUEST.fields_by_name['name']._options = None
_DESCRIBESERVICEV1REQUEST.fields_by_name['service_id']._options = None
_REMOVESERVICEV1REQUEST.fields_by_name['service_id']._options = None

_SSNSERVICEAPISERVICE = _descriptor.ServiceDescriptor(
  name='SsnServiceApiService',
  full_name='ozonmp.ssn_service_api.v1.SsnServiceApiService',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=685,
  serialized_end=1335,
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
    name='ListServicesV1',
    full_name='ozonmp.ssn_service_api.v1.SsnServiceApiService.ListServicesV1',
    index=2,
    containing_service=None,
    input_type=_LISTSERVICESV1REQUEST,
    output_type=_LISTSERVICESV1RESPONSE,
    serialized_options=_b('\202\323\344\223\002\022\022\020/api/v1/services'),
  ),
  _descriptor.MethodDescriptor(
    name='RemoveServiceV1',
    full_name='ozonmp.ssn_service_api.v1.SsnServiceApiService.RemoveServiceV1',
    index=3,
    containing_service=None,
    input_type=_REMOVESERVICEV1REQUEST,
    output_type=_REMOVESERVICEV1RESPONSE,
    serialized_options=_b('\202\323\344\223\002\037*\035/api/v1/services/{service_id}'),
  ),
])
_sym_db.RegisterServiceDescriptor(_SSNSERVICEAPISERVICE)

DESCRIPTOR.services_by_name['SsnServiceApiService'] = _SSNSERVICEAPISERVICE

# @@protoc_insertion_point(module_scope)
