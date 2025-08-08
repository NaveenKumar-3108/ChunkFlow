## Sample API Requests 
1. API: To upload chunk
Method: POST
Route: /upload
Host: localhost:8080

**Sample Request** 
Content-Type: multipart/form-data;

Content-Disposition: form-data; name="user_id"
12345

Content-Disposition: form-data; name="session_id"
abc-session

Content-Disposition: form-data; name="timestamp"
2025-08-07T10:30:00Z

Content-Disposition: form-data; name="audio"; filename="sample.mp3"
Content-Type: audio/wav,audio/mp3

**Sample Response** (Success):
```json
{
    "chunkid": 1,
    "status": "S",
    "errmsg": "Audio uploaded successfully",
    "msg": ""
}

```

**Sample Response** (error):
```json
{
    "chunkid": 0,
    "status": "E",
    "errmsg": "Missing required metadata",
    "msg": ""
}

```

2. API: To upload chunk
Method: GET
Route: /chunks{chunk_id}
Host: localhost:8080


**Sample Response** (Success):
```json
{
    "MetaData": {
        "ChunkID": "2",
        "UserID": "2",
        "SessionID": "311",
        "Timestamp": "3213",
        "FileName": "Hi there this is you.mp3",
        "ContentType": "audio/mpeg",
        "Duration": "1",
        "Transcript": "Welcome ",
        "Checksum": "2af1519246104ca45f9d59b78a6d5e2662ed64213feae5ea423a8747645470f6",
        "Size": "85293"
    },
    "status": "S",
    "errmsg": ""
}

```

**Sample Response** (error):
```json
{
    "MetaData": {
        "ChunkID": "",
        "UserID": "",
        "SessionID": "",
        "Timestamp": "",
        "FileName": "",
        "ContentType": "",
        "Duration": "",
        "Transcript": " ",
        "Checksum": "",
        "Size": ""
    },
    "status": "E",
    "errmsg": "something went wrong"
}
```


2. API: To upload chunk
Method: GET
Route: /session{session_id}
Host: localhost:8080


**Sample Response** (Success):
```json
{
    "userchunks": [
        {
            "ChunkID": "2",
            "UserID": "2",
            "SessionID": "311",
            "Timestamp": "3213",
            "FileName": "Hi there this is you.mp3",
            "ContentType": "audio/mpeg",
            "Duration": "5",
            "Transcript": "Welcome ",
            "Checksum": "2af1519246104ca45f9d59b78a6d5e2662ed64213feae5ea423a8747645470f6",
            "Size": "85293"
        }
    ],
    "status": "S",
    "errmsg": ""
}

```

**Sample Response** (error):
```json
{
    "userchunks": [
        {
            "ChunkID": "",
            "UserID": "",
            "SessionID": "",
            "Timestamp": "",
            "FileName": "",
            "ContentType": "",
            "Duration": "",
            "Transcript": "",
            "Checksum": "",
            "Size": ""
        }
    ],
    "status": "E",
    "errmsg": "something went wrong"
        }

```
