import asyncio

from grpclib.client import Channel

from ozonmp.ssn_service_api.v1.ssn_service_api_grpc import SsnServiceApiServiceStub
from ozonmp.ssn_service_api.v1.ssn_service_api_pb2 import DescribeServiceV1Request

async def main():
    async with Channel('127.0.0.1', 8082) as channel:
        client = SsnServiceApiServiceStub(channel)

        req = DescribeServiceV1Request(service_id=1)
        reply = await client.DescribeServiceV1(req)
        print(reply.message)


if __name__ == '__main__':
    asyncio.run(main())
