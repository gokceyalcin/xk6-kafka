# xk6-kafka

**`description`**
The xk6-kafka project is a k6 extension that enables k6 users to load test Apache Kafka using a producer and possibly a consumer for debugging.
This documentation refers to the development version of the xk6-kafka project, which means the latest changes on `main` branch and might not be released yet, as explained in [the release process](https://github.com/mostafa/xk6-kafka#the-release-process).

**`see`** [https://github.com/mostafa/xk6-kafka](https://github.com/mostafa/xk6-kafka)

## Table of contents

### Enumerations

- [BALANCERS](enums/BALANCERS.md)
- [COMPRESSION_CODECS](enums/COMPRESSION_CODECS.md)
- [ELEMENT_TYPES](enums/ELEMENT_TYPES.md)
- [GROUP_BALANCERS](enums/GROUP_BALANCERS.md)
- [ISOLATION_LEVEL](enums/ISOLATION_LEVEL.md)
- [SASL_MECHANISMS](enums/SASL_MECHANISMS.md)
- [SCHEMA_TYPES](enums/SCHEMA_TYPES.md)
- [START_OFFSETS](enums/START_OFFSETS.md)
- [SUBJECT_NAME_STRATEGY](enums/SUBJECT_NAME_STRATEGY.md)
- [TLS_VERSIONS](enums/TLS_VERSIONS.md)

### Classes

- [Connection](classes/Connection.md)
- [Reader](classes/Reader.md)
- [SchemaRegistry](classes/SchemaRegistry.md)
- [Writer](classes/Writer.md)

### Interfaces

- [BasicAuth](interfaces/BasicAuth.md)
- [ConfigEntry](interfaces/ConfigEntry.md)
- [ConnectionConfig](interfaces/ConnectionConfig.md)
- [ConsumeConfig](interfaces/ConsumeConfig.md)
- [Container](interfaces/Container.md)
- [Message](interfaces/Message.md)
- [ProduceConfig](interfaces/ProduceConfig.md)
- [ReaderConfig](interfaces/ReaderConfig.md)
- [ReplicaAssignment](interfaces/ReplicaAssignment.md)
- [SASLConfig](interfaces/SASLConfig.md)
- [Schema](interfaces/Schema.md)
- [SchemaRegistryConfig](interfaces/SchemaRegistryConfig.md)
- [SubjectNameConfig](interfaces/SubjectNameConfig.md)
- [TLSConfig](interfaces/TLSConfig.md)
- [TopicConfig](interfaces/TopicConfig.md)
- [WriterConfig](interfaces/WriterConfig.md)
