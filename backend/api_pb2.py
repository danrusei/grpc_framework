# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: api.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='api.proto',
  package='api',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=_b('\n\tapi.proto\x12\x03\x61pi\"1\n\rClientRequest\x12\x0e\n\x06vendor\x18\x01 \x01(\t\x12\x10\n\x08prodType\x18\x02 \x01(\t\"\"\n\x0e\x43lientResponse\x12\x10\n\x08products\x18\x01 \x01(\t\".\n\nApiRequest\x12\x0e\n\x06vendor\x18\x01 \x01(\t\x12\x10\n\x08prodType\x18\x02 \x01(\t\"\x1c\n\x0b\x41piResponse\x12\r\n\x05prods\x18\x01 \x03(\t2B\n\x0bProdService\x12\x33\n\x08GetProds\x12\x12.api.ClientRequest\x1a\x13.api.ClientResponse2D\n\x0e\x42\x61\x63kendService\x12\x32\n\rRetrieveItems\x12\x0f.api.ApiRequest\x1a\x10.api.ApiResponseb\x06proto3')
)




_CLIENTREQUEST = _descriptor.Descriptor(
  name='ClientRequest',
  full_name='api.ClientRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='vendor', full_name='api.ClientRequest.vendor', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='prodType', full_name='api.ClientRequest.prodType', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
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
  serialized_start=18,
  serialized_end=67,
)


_CLIENTRESPONSE = _descriptor.Descriptor(
  name='ClientResponse',
  full_name='api.ClientResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='products', full_name='api.ClientResponse.products', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
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
  serialized_start=69,
  serialized_end=103,
)


_APIREQUEST = _descriptor.Descriptor(
  name='ApiRequest',
  full_name='api.ApiRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='vendor', full_name='api.ApiRequest.vendor', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='prodType', full_name='api.ApiRequest.prodType', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
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
  serialized_start=105,
  serialized_end=151,
)


_APIRESPONSE = _descriptor.Descriptor(
  name='ApiResponse',
  full_name='api.ApiResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='prods', full_name='api.ApiResponse.prods', index=0,
      number=1, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
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
  serialized_start=153,
  serialized_end=181,
)

DESCRIPTOR.message_types_by_name['ClientRequest'] = _CLIENTREQUEST
DESCRIPTOR.message_types_by_name['ClientResponse'] = _CLIENTRESPONSE
DESCRIPTOR.message_types_by_name['ApiRequest'] = _APIREQUEST
DESCRIPTOR.message_types_by_name['ApiResponse'] = _APIRESPONSE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

ClientRequest = _reflection.GeneratedProtocolMessageType('ClientRequest', (_message.Message,), {
  'DESCRIPTOR' : _CLIENTREQUEST,
  '__module__' : 'api_pb2'
  # @@protoc_insertion_point(class_scope:api.ClientRequest)
  })
_sym_db.RegisterMessage(ClientRequest)

ClientResponse = _reflection.GeneratedProtocolMessageType('ClientResponse', (_message.Message,), {
  'DESCRIPTOR' : _CLIENTRESPONSE,
  '__module__' : 'api_pb2'
  # @@protoc_insertion_point(class_scope:api.ClientResponse)
  })
_sym_db.RegisterMessage(ClientResponse)

ApiRequest = _reflection.GeneratedProtocolMessageType('ApiRequest', (_message.Message,), {
  'DESCRIPTOR' : _APIREQUEST,
  '__module__' : 'api_pb2'
  # @@protoc_insertion_point(class_scope:api.ApiRequest)
  })
_sym_db.RegisterMessage(ApiRequest)

ApiResponse = _reflection.GeneratedProtocolMessageType('ApiResponse', (_message.Message,), {
  'DESCRIPTOR' : _APIRESPONSE,
  '__module__' : 'api_pb2'
  # @@protoc_insertion_point(class_scope:api.ApiResponse)
  })
_sym_db.RegisterMessage(ApiResponse)



_PRODSERVICE = _descriptor.ServiceDescriptor(
  name='ProdService',
  full_name='api.ProdService',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=183,
  serialized_end=249,
  methods=[
  _descriptor.MethodDescriptor(
    name='GetProds',
    full_name='api.ProdService.GetProds',
    index=0,
    containing_service=None,
    input_type=_CLIENTREQUEST,
    output_type=_CLIENTRESPONSE,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_PRODSERVICE)

DESCRIPTOR.services_by_name['ProdService'] = _PRODSERVICE


_BACKENDSERVICE = _descriptor.ServiceDescriptor(
  name='BackendService',
  full_name='api.BackendService',
  file=DESCRIPTOR,
  index=1,
  serialized_options=None,
  serialized_start=251,
  serialized_end=319,
  methods=[
  _descriptor.MethodDescriptor(
    name='RetrieveItems',
    full_name='api.BackendService.RetrieveItems',
    index=0,
    containing_service=None,
    input_type=_APIREQUEST,
    output_type=_APIRESPONSE,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_BACKENDSERVICE)

DESCRIPTOR.services_by_name['BackendService'] = _BACKENDSERVICE

# @@protoc_insertion_point(module_scope)
