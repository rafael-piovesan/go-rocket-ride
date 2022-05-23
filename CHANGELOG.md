# Changelog

## [2.3.0](https://github.com/rafael-piovesan/go-rocket-ride/compare/v2.2.0...v2.3.0) (2022-05-23)


### Features

* add unit-of-work ([c81a6c2](https://github.com/rafael-piovesan/go-rocket-ride/commit/c81a6c2f36b77e64f1438e1f9c444ba2bae1c315))

## [2.2.0](https://github.com/rafael-piovesan/go-rocket-ride/compare/v2.1.0...v2.2.0) (2022-05-23)


### Features

* move from repo to crud store ([c136dad](https://github.com/rafael-piovesan/go-rocket-ride/commit/c136dadf3e9fcd1f8f862b9a1f5758c92813096f))

## [2.1.0](https://github.com/rafael-piovesan/go-rocket-ride/compare/v2.0.0...v2.1.0) (2022-05-23)


### Features

* **httpserver:** handle errors ([c10a916](https://github.com/rafael-piovesan/go-rocket-ride/commit/c10a9163ce172ea951747bdfcbd6573d51420708))
* **pkg:** add config pkg ([174c5a2](https://github.com/rafael-piovesan/go-rocket-ride/commit/174c5a2ea50522b7ed0323b2ab57a973a15b7b31))
* **pkg:** add httpserver pkg ([d3987b5](https://github.com/rafael-piovesan/go-rocket-ride/commit/d3987b5530c3fa3afbc5a62101bfc342149f859f))
* **pkg:** add repo & db pkg ([9085c71](https://github.com/rafael-piovesan/go-rocket-ride/commit/9085c7171531616bf6348d21b86cb698a64e9eba))


### Bug Fixes

* **testcontainer:** load pg driver ([a86e316](https://github.com/rafael-piovesan/go-rocket-ride/commit/a86e316abeb1ac4607a633eba2d066f0c43e6ea1))

## [2.0.0](https://github.com/rafael-piovesan/go-rocket-ride/compare/v1.0.0...v2.0.0) (2022-05-22)


### ⚠ BREAKING CHANGES

* simplifying datastore pkg
* bumping version to v2
* remove sqlc

### Code Refactoring

* bumping version to v2 ([ce2e66e](https://github.com/rafael-piovesan/go-rocket-ride/commit/ce2e66e690c744b92782911d1e34c560a6fddad4))
* remove sqlc ([86f2094](https://github.com/rafael-piovesan/go-rocket-ride/commit/86f20944ab3105a1c4775a45df0a72e73ebcbbc8))
* simplifying datastore pkg ([9b57e49](https://github.com/rafael-piovesan/go-rocket-ride/commit/9b57e49cf7fc37eb56c35c31a2c815470552cb63))

## 1.0.0 (2022-05-22)


### Features

* **bun:** add bun based datastore ([1cb5c2d](https://github.com/rafael-piovesan/go-rocket-ride/commit/1cb5c2d01cb5acc5746a4f1e5d6ccd789e979a70))
* **cmd:** add app main.go ([be1dd35](https://github.com/rafael-piovesan/go-rocket-ride/commit/be1dd35fc286e9828df468b3c896e022a7c1b325))
* **config:** add app config ([69679e3](https://github.com/rafael-piovesan/go-rocket-ride/commit/69679e3f019a795726076f29efde043ab98a815b))
* **config:** validating app cfg ([9d0f978](https://github.com/rafael-piovesan/go-rocket-ride/commit/9d0f9785c25796a6ba61640a9229209c21d165a1))
* **entities:** add app entities ([77dc9f2](https://github.com/rafael-piovesan/go-rocket-ride/commit/77dc9f2eb21d9e1a22e5c5c9d1dfb38b67b589f3))
* **entity:** add new errors ([d4345d6](https://github.com/rafael-piovesan/go-rocket-ride/commit/d4345d6ef30078d4a761438134a61b3671571e0b))
* **errors:** move errors to entity pkg ([70390b8](https://github.com/rafael-piovesan/go-rocket-ride/commit/70390b849226ad28354a55518d5b37ae061d422b))
* **http:** add http port ([8dafed3](https://github.com/rafael-piovesan/go-rocket-ride/commit/8dafed34507ed1731a3b000c28c2fcc9f29b99e7))
* **interfaces:** add app interfaces ([d750014](https://github.com/rafael-piovesan/go-rocket-ride/commit/d750014f27ad76415fd4e5f8e40a88a85214ac34))
* **makefile:** add makefile ([bde73ae](https://github.com/rafael-piovesan/go-rocket-ride/commit/bde73ae0b76e40b9462d3594e069590387eecb85))
* **middleware:** add origin ip middleware ([124e495](https://github.com/rafael-piovesan/go-rocket-ride/commit/124e49563b62c15fe4adc175e40535396a1b0bef))
* **migrations:** add db migrations ([cbf10e2](https://github.com/rafael-piovesan/go-rocket-ride/commit/cbf10e29d5b7a09dec0c6f180ea1584689de1fc6))
* **mocks:** add app mocks ([a9ab14e](https://github.com/rafael-piovesan/go-rocket-ride/commit/a9ab14eb893179a91c653b498a0f7765c236507a))
* **pkg:** add helper packages ([37edd25](https://github.com/rafael-piovesan/go-rocket-ride/commit/37edd2531628a2912e7b16f5e6dc91008b0b17df))
* **pkg:** add stripe mock integration ([823d2fc](https://github.com/rafael-piovesan/go-rocket-ride/commit/823d2fca588c64efef1e05dba325dc185fde9587))
* **pkg:** stripe-mock error mode ([f9923c6](https://github.com/rafael-piovesan/go-rocket-ride/commit/f9923c669db5bcae16e478fb85015dcc8e652b09))
* **ride:** add ride uc ([176d455](https://github.com/rafael-piovesan/go-rocket-ride/commit/176d455dda321bafadc895dbd7989fd5caba7019))
* **sqlc:** add generated files ([3a16ce4](https://github.com/rafael-piovesan/go-rocket-ride/commit/3a16ce473a10f7c3ded376a5c7c3211809dbfd51))
* **sqlc:** add queries statements ([77734ad](https://github.com/rafael-piovesan/go-rocket-ride/commit/77734ad78c5b8f86f7d3dde3a59f34bb13d1f04b))
* **sqlc:** add sqlc config ([cb6d576](https://github.com/rafael-piovesan/go-rocket-ride/commit/cb6d576ff57eff6be984805ec3abd35009a59bd7))
* **sqlc:** add sqlc datastore ([#2](https://github.com/rafael-piovesan/go-rocket-ride/issues/2)) ([16899a2](https://github.com/rafael-piovesan/go-rocket-ride/commit/16899a281bbb1b6589f3ff316d561ba8e9064db2))
* **stripe:** add stripe integration ([9beee35](https://github.com/rafael-piovesan/go-rocket-ride/commit/9beee350523fbe3cd935ab7fa45362ab5b2607fe))
* **stripemock:** add init check featrure toggle ([c7eea63](https://github.com/rafael-piovesan/go-rocket-ride/commit/c7eea63bb01ca4648f119ef597f190aee3fa1261))
* **testfixtures:** add db fixtures ([2ffa93c](https://github.com/rafael-piovesan/go-rocket-ride/commit/2ffa93ceb0fe322e0f6c0d726288d0b04311f8d7))
* **tests:** add datastore integration tests ([3211eb6](https://github.com/rafael-piovesan/go-rocket-ride/commit/3211eb6668ec6032743a4ddecfdf11071fdfb41f))
* **tests:** add db integration tests ([9092d81](https://github.com/rafael-piovesan/go-rocket-ride/commit/9092d810ee36c7cfdd90e63dd86d27f4b00eb072))
* **tests:** add main.go tests ([c19c61e](https://github.com/rafael-piovesan/go-rocket-ride/commit/c19c61e0ab4e39a5b2531d0117965d54dda04aaf))
* **tests:** add ride uc unit tests ([7f3c15c](https://github.com/rafael-piovesan/go-rocket-ride/commit/7f3c15c81f9eeba07eec6eda79b1b4c85857efa5))
* **tests:** add server shutdown tests ([a598bb3](https://github.com/rafael-piovesan/go-rocket-ride/commit/a598bb3d601c07dc1b463b5a779a9c3b00e6c8e4))
* **tests:** testing stripe calls ([d910be2](https://github.com/rafael-piovesan/go-rocket-ride/commit/d910be26204f024734984a160c92dededf38690f))
* **tools:** pkg to keep track of dev deps ([22d7d15](https://github.com/rafael-piovesan/go-rocket-ride/commit/22d7d15867fcea43f8ccfd7907b618604f2239c6))
* **usecase:** read ik timeout from config ([bc84740](https://github.com/rafael-piovesan/go-rocket-ride/commit/bc84740e112af1da365fe8e023fa910127352bc4))
* **user:** add user datastore methods ([a6faff6](https://github.com/rafael-piovesan/go-rocket-ride/commit/a6faff6246dae81ec6ef88cd38ff978cfaf71239))
* **user:** set & get user info from ctx ([66da7f6](https://github.com/rafael-piovesan/go-rocket-ride/commit/66da7f6f2a72b4c190c4042c7386f04dfd73829c))
* **workflows:** add lint and tests workflows ([#1](https://github.com/rafael-piovesan/go-rocket-ride/issues/1)) ([e0f9552](https://github.com/rafael-piovesan/go-rocket-ride/commit/e0f9552061926b46892af1b7d9971d666da2032a))


### Bug Fixes

* **entities:** correct created_at field ([6410c74](https://github.com/rafael-piovesan/go-rocket-ride/commit/6410c749b6986633f2fdb25a122be4086b7ab71b))
* **handler:** fix typo ([dc605dd](https://github.com/rafael-piovesan/go-rocket-ride/commit/dc605ddcce4f3bfbb2eea8dd868c7cc9fe2de2a8))
* **handler:** init input to fix check ([dbbdab7](https://github.com/rafael-piovesan/go-rocket-ride/commit/dbbdab728ce0de20736d649ab0687c5d0194c0cc))
* **handler:** standardize error return ([6e8ce2b](https://github.com/rafael-piovesan/go-rocket-ride/commit/6e8ce2b637864b53b9f87eb3498aa57c4d4a5663))
* **ride:** using stripe customer id ([1ca283d](https://github.com/rafael-piovesan/go-rocket-ride/commit/1ca283d990db9a22400cc5fcf10488025730a9c9))
* **server:** improving logs ([4afdffe](https://github.com/rafael-piovesan/go-rocket-ride/commit/4afdffea2e9aca30bc04fd9d74c01ffc8092b7fa))
* **sqlstore:** propagate panic ([331b028](https://github.com/rafael-piovesan/go-rocket-ride/commit/331b02834c09e5f61fe3d116d4ae76eaa733ef47))
* **tests:** turn off gock ([9d214b4](https://github.com/rafael-piovesan/go-rocket-ride/commit/9d214b486ca0af4efda3f6b2b4d67ce07df3ac54))
* **usecase:** add random stripe-id ([6662793](https://github.com/rafael-piovesan/go-rocket-ride/commit/6662793e8eae523f8a801710be73f7844d39e92a))
* **usecase:** correct ride ID on audti record ([b9e2ac5](https://github.com/rafael-piovesan/go-rocket-ride/commit/b9e2ac59d23972d2e8d46433f2c61983dcb05b29))
* **usecase:** idem key vs request compare structs ([6e2c57a](https://github.com/rafael-piovesan/go-rocket-ride/commit/6e2c57abd6756e267f00cc5ab9f2ac4da3c5f1e5))
* **usecase:** unlock idem key ([42f896d](https://github.com/rafael-piovesan/go-rocket-ride/commit/42f896d7b9c8508841f1118a23970e7f3d050389))
* **user:** user info on ctx ([b45e41d](https://github.com/rafael-piovesan/go-rocket-ride/commit/b45e41d529006909ec822e2783accdc3d4a8a89a))
