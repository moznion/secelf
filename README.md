secelf
==

A secret shelf that uses Google Drive as backend storage.

Overview
--

```
+--------+                +----------+                        +----------------+
|        |  upload file   |          |  encrypt file by AES   |                |
|        | -------------> |          | ---------------------> |                |
|  User  |                |  secelf  |                        |  Google Drive  |
|        | <------------- |          | <--------------------- |                |
|        |   serve file   |          |      decrypt file      |                |
+--------+                +----------+                        +----------------+
```

How to setup
--

1. setup SQLite3 database
    - `sqlite3 YOUR_DB.sqlite3 < ./sql/000-file.sql`
2. get a credential and a token from Google Drive
3. run secelf!

How to run
--

```
Usage of secelf
  -credential-json string
        [mandatory] credential of Google Drive as JSON string
  -key string
        [mandatory] key for file encryption (must be 128bit, 192bit or 256bit)
  -port int
        [mandatory] port for listen (default -1)
  -root-dir-id string
        [mandatory] identifier fo root directory for storing files
  -sqlite-db-path string
        [mandatory] path to SQLite DB file
  -token-json string
        [mandatory] token for accessing to Google Drive as JSON string
```

Example:

```
secelf \
  --token-json="$(cat token.json)" \
  --credential-json="$(cat credentials.json)" \
  --key='this-aes-key-has-32-charactersss' \
  --root-dir-id='your-google-drive-dir-ID' \
  --sqlite-db-path='/path/to/YOUR_DB.sqlite3' \
  --port=8888
```

👍 then you can access the application as `127.0.0.1:8888`.

![image](https://user-images.githubusercontent.com/1422834/48304468-478de180-e55d-11e8-8ee9-6d1e7aa0c570.png)

Example of `token.json`:

```json
{"access_token":"your-token","token_type":"Bearer","refresh_token":"your-refresh-token","expiry":"2018-11-02T22:25:25.321592+09:00"}
```

Example of `credentials.json`:

```json
{"installed":{"client_id":"your-client-id","project_id":"your-project-id","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://www.googleapis.com/oauth2/v3/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_secret":"your-client-secret","redirect_uris":["urn:ietf:wg:oauth:2.0:oob"]}}
```

License
--

```
Copyright 2018 moznion, http://moznion.net/ <moznion@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

