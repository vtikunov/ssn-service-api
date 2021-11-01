import setuptools

setuptools.setup(
    name="grpc-ssn-service-api",
    version="0.3.1",
    author="vtikunov",
    author_email="vtikunov@yandex.ru",
    description="GRPC python client for ssn-service-api",
    url="https://github.com/ozonmp/ssn-service-api",
    packages=setuptools.find_packages(),
    python_requires='>=3.5',
)