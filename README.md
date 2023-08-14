# My pet - Mockster
---
## Todo plan
- [X] Come up with a schema for yaml files for mock endpoints
- [X] Create yaml parsers and validators
- [X] Choose a DB to store data (thats should be file or in memory database, document or kv store) `I choose to try to use in memory hash map, write some interfaces and then try to use boltdb or rocksdb`
- [X] Write interfaces and repos for storage `In memory hash map is ready`
- [ ] Write mock-serving logic
- [ ] Write handlers for getting yaml files with http(with append and replace options)
- [ ] Write handlers for management mocks by rest api
- [ ] Write handlers for import and export data with yaml files
- [ ] Write handlers for creating mocks by rest api
- [ ] Write some tests
- [ ] Write simple UI for managing mocks
- [ ] Write simple UI for creating mocks
- [ ] Write CI
- [ ] Write paper about this project =)
## ever or never
- [ ] Add BoltDB or RocksDB as storage
- [ ] Mocks for different protocols (ws, grpc, etc)
- [ ] Write api for checking incoming requests
- [ ] Add some metrics
- [ ] tbd

