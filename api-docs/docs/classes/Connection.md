# Class: Connection

**`classdesc`** Connection connects to Kafka for working with topics.

**`example`**

```javascript
// In init context
const connection = new Connection({
  address: "localhost:9092",
});

// In VU code (default function)
const topics = connection.listTopics();

// In teardown function
connection.close();
```

## Table of contents

### Constructors

- [constructor](Connection.md#constructor)

### Methods

- [close](Connection.md#close)
- [createTopic](Connection.md#createtopic)
- [deleteTopic](Connection.md#deletetopic)
- [listTopics](Connection.md#listtopics)

## Constructors

### constructor

• **new Connection**(`connectionConfig`)

#### Parameters

| Name               | Type                                                    | Description               |
| :----------------- | :------------------------------------------------------ | :------------------------ |
| `connectionConfig` | [`ConnectionConfig`](../interfaces/ConnectionConfig.md) | Connection configuration. |

#### Defined in

[index.d.ts:361](https://github.com/mostafa/xk6-kafka/blob/main/api-docs/index.d.ts#L361)

## Methods

### close

▸ **close**(): `void`

**`destructor`**

**`description`** Close the connection.

#### Returns

`void`

- Nothing.

#### Defined in

[index.d.ts:387](https://github.com/mostafa/xk6-kafka/blob/main/api-docs/index.d.ts#L387)

---

### createTopic

▸ **createTopic**(`topicConfig`): `void`

**`method`**
Create a new topic.

#### Parameters

| Name          | Type                                          | Description          |
| :------------ | :-------------------------------------------- | :------------------- |
| `topicConfig` | [`TopicConfig`](../interfaces/TopicConfig.md) | Topic configuration. |

#### Returns

`void`

- Nothing.

#### Defined in

[index.d.ts:368](https://github.com/mostafa/xk6-kafka/blob/main/api-docs/index.d.ts#L368)

---

### deleteTopic

▸ **deleteTopic**(`topic`): `void`

**`method`**
Delete a topic.

#### Parameters

| Name    | Type     | Description |
| :------ | :------- | :---------- |
| `topic` | `string` | Topic name. |

#### Returns

`void`

- Nothing.

#### Defined in

[index.d.ts:375](https://github.com/mostafa/xk6-kafka/blob/main/api-docs/index.d.ts#L375)

---

### listTopics

▸ **listTopics**(): `string`[]

**`method`**
List topics.

#### Returns

`string`[]

- Topics.

#### Defined in

[index.d.ts:381](https://github.com/mostafa/xk6-kafka/blob/main/api-docs/index.d.ts#L381)
