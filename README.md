# Instruction

### 1 Step

First you need to clone this repo. For example, with SSH link:

```sh
git clone git@github.com:dmitruk-v/polls-service.git
```

### 2 Step

To run service just enter command:

```sh
make dc-build-all
```

It will download all needed images, compiles go application and runs containers.

### 3 Step

To check running containers:

```sh
make dc-ps
```

### 4 Step

Application container exposes port 8080, so URL to endpoint is:

```sh
POST http://localhost:8080/polls
```
It expects JSON payload, such as:

```json
{   
    "survey_id": 5,
    "pre_set_values": {
        "Who let the dogs out?": "Police",
        "В чем смысл жизни?": "В Кузинатре",
        "Де одружуються бджоли?": "У ЗАГС-і"
    }
}
```

### 5 Step

To stop all running containers:

```sh
make dc-down
```
