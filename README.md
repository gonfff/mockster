# My pet - Mockster
---
## Todo plan
- [X] Come up with a schema for yaml files for mock endpoints
- [X] Create yaml parsers and validators
- [X] Choose a DB to store data (thats should be file or in memory database, document or kv store) `I choose to try to use in memory hash map, write some interfaces and then try to use boltdb or rocksdb`
- [X] Write interfaces and repos for storage `In memory hash map is ready`
- [X] Write mock-serving logic `But with constrains. Str combination of method+path+body is used as key for searching mock. So if you have two mocks with same method, path and body, you will get only first mock. `
- [X] I will try to add routes on the fly on another Echo instance, and restart it on update/delete mock, i am not sure that it will work, but i will try `I decide to use path variable for serve mocks, register on the fly have a much more complexity and blurring logic of the project. but i save a draft of this idea in the branch "routes-the-fly"`
- [X] Write handlers for getting yaml files with http(replace)
- [X] Write handlers for management mocks by rest api
- [X] Write handlers for import and export data with yaml files
- [X] Write handlers for creating mocks by rest api
- [X] Write some tests
- [X] Write simple UI for managing mocks
- [X] Write simple UI for creating mocks
- [ ] Write CI
- [ ] Write paper about this project =)
## ever or never
- [ ] Add BoltDB or RocksDB as storage
- [ ] Mocks for different protocols (ws, grpc, etc)
- [ ] Write api for checking incoming requests
- [ ] Add some metrics
- [ ] tbd
