# Changelog

## 1.0.0 (2025-10-28)


### Features

* activate user route, handler & DB logic ([862da78](https://github.com/Tiskae/gopher-social/commit/862da78047de3d6140f8db1fdad006f09dcad6c9))
* added db posts and users migrations ([73ba9b7](https://github.com/Tiskae/gopher-social/commit/73ba9b7bef843cf23df406300f11616ec418ed70))
* added filtering by posts tags to user feed ([ee6e048](https://github.com/Tiskae/gopher-social/commit/ee6e0488e2ff095ac94d963f7d8d698be5455e72))
* added GitHub actions ([7061bfe](https://github.com/Tiskae/gopher-social/commit/7061bfe8f9f2b9c762e30a52480ca5f2eecd8d0f))
* added pagination and sorting to useer feed ([d071723](https://github.com/Tiskae/gopher-social/commit/d071723fd925e080b7983d238d56ef590d508c47))
* added validation to post create payload ([78ae9b2](https://github.com/Tiskae/gopher-social/commit/78ae9b2c0c75fd3224406b13cbdb38e8b4d2c60a))
* added version to posts to prevent update post race conflict ([fe2d6f7](https://github.com/Tiskae/gopher-social/commit/fe2d6f7c12801d9c96370e002c9b085c805f22a6))
* auth token route, handler and DB logic ([787521f](https://github.com/Tiskae/gopher-social/commit/787521fea0e5da14c1094abec92240c136dd2cc7))
* authentication on users, posts and feed routes ([356ac01](https://github.com/Tiskae/gopher-social/commit/356ac01b0478c397e1d9d124ac4219360b587343))
* authorization for post deletion and updates ([9d1b873](https://github.com/Tiskae/gopher-social/commit/9d1b873f30d943c2aaf95579b579a043d2589b6f))
* crate user route, handler & DB logic ([e4ddb97](https://github.com/Tiskae/gopher-social/commit/e4ddb97a78fd6ea28128616846b160e91804d472))
* create post comment route and handler ([47a2714](https://github.com/Tiskae/gopher-social/commit/47a2714699cd59610da58846b7ce40c9e7910311))
* create post handler ([95fc47c](https://github.com/Tiskae/gopher-social/commit/95fc47c8993f9395bf4b8c0ed7d29d3de9bd2408))
* create user endpoint ([24125cb](https://github.com/Tiskae/gopher-social/commit/24125cb554de118fc39b66ed43ae46b65d2fc56f))
* delete post route, handler and DB logic ([0aa6a63](https://github.com/Tiskae/gopher-social/commit/0aa6a6325868a3cc58af0d40469fb491a5e07bd1))
* dockerfile ([f625b05](https://github.com/Tiskae/gopher-social/commit/f625b0578545ce5c6f25a4ac754da1f88b64c4ca))
* follow and unfollow users routes, handlers and DB logics ([9158596](https://github.com/Tiskae/gopher-social/commit/9158596516e32060e45c82f8843373b1e3f01c55))
* get post by id route & handler ([e957a3f](https://github.com/Tiskae/gopher-social/commit/e957a3ff253eb1a501c3a995e4bfeb567d873d63))
* get post comment route, handler and DB logic ([dd5a50f](https://github.com/Tiskae/gopher-social/commit/dd5a50fe1ca6b920a9822c32a36320d7cdd7b2af))
* get user by ID route, handler and DB logic ([a6524a7](https://github.com/Tiskae/gopher-social/commit/a6524a7229bc691db7750f6cc0f37677c345c225))
* redis caching for get user endpoints ([b0b0f79](https://github.com/Tiskae/gopher-social/commit/b0b0f79af2a5fc77847cfb17447a6b3535e7d363))
* release please script ([6acfadf](https://github.com/Tiskae/gopher-social/commit/6acfadf1d72d6b1426ed6993da32473e3d2b256f))
* send comments with post ([f570a08](https://github.com/Tiskae/gopher-social/commit/f570a08857f6134ead52b0c803b613d3cccab0c0))
* setup db connection pool ([7f5e520](https://github.com/Tiskae/gopher-social/commit/7f5e5203ade516a09854c607600b5c2b57830f75))
* setup mailer ([c51e251](https://github.com/Tiskae/gopher-social/commit/c51e2518be615f3c6735a59f0de37004edc039c7))
* update post route, handler and DB logic ([ff96453](https://github.com/Tiskae/gopher-social/commit/ff96453b1f9e41b5200773fd24a0b0c329d231fd))
* user email confirmation mail setup ([841786b](https://github.com/Tiskae/gopher-social/commit/841786b467f0745068bd59602f53ef3b513b3820))
* user feed route, handler and DB logic ([281c378](https://github.com/Tiskae/gopher-social/commit/281c3789b2c736da7a58e535ed16593655da8492))


### Bug Fixes

* GH audit workflow ([46df7b2](https://github.com/Tiskae/gopher-social/commit/46df7b226aae281c565f6dc5d3734e59934f2744))
* GitHub actions ([0332f17](https://github.com/Tiskae/gopher-social/commit/0332f170dda3110c674bec179b0ceb68de84b79d))
* GitHub actions folder structure ([b4bec08](https://github.com/Tiskae/gopher-social/commit/b4bec08f1adac741d38f8902d1bc903cf2c0c44f))
* merge conflict ([0d8cff3](https://github.com/Tiskae/gopher-social/commit/0d8cff3121787bec76f1d9eb375584cbd27b1c2a))
* omit empty user fields on JSON writes ([9b4daf2](https://github.com/Tiskae/gopher-social/commit/9b4daf21a1a95c12af015b986d0422caac2c7d58))
* password check on auth token generation ([788eabf](https://github.com/Tiskae/gopher-social/commit/788eabf985e4549c1f033d1f63223de21bb27223))
* return statement ([9580627](https://github.com/Tiskae/gopher-social/commit/9580627e62e1a265938719c453a1fed689633e2b))
* switched runner for GH audit workflow ([d20e267](https://github.com/Tiskae/gopher-social/commit/d20e26773c73fe5d34fa396a2e16eb64c83b9a99))
* typo ([967a911](https://github.com/Tiskae/gopher-social/commit/967a911254fb08b0f0a859632572da6415d07736))


### Performance Improvements

* added indexes to crucial DB columns ([ee09032](https://github.com/Tiskae/gopher-social/commit/ee090328b5abf53043a8698b74bd045cac62ac76))
* rate limiting all endpoints ([d060e9b](https://github.com/Tiskae/gopher-social/commit/d060e9bbe754791c030e1c7709549eae504c7b97))
