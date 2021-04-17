# omgur

[![pipeline status](https://git.geraldwu.com/gerald/omgur/badges/master/pipeline.svg)](https://git.geraldwu.com/gerald/omgur/-/commits/master) 

Omgur is a free and open-source alternative Imgur front-end focused on privacy.

Inspired by the [Invidious](https://github.com/iv-org/invidious), [Nitter](https://github.com/zedeus/nitter), and [Teddit](https://github.com/teddit-net/teddit) projects.

- No JavaScript or ads
- All requests go through the backend, client never talks to Imgur
- Prevents Imgur from tracking your IP or JavaScript fingerprint
- Lightweight
- Self-hostable

## Features roadmap

- [x] Direct image loading via https://i.imgur.com/
- [x] Imgur album loading via https://imgur.com/a/
- [x] Imgur gallery loading via https://imgur.com/gallery/
- [x] Redis caching for images
- [ ] Imgur post loading via https://imgur.com/
- [ ] Proper embedding of videos on albums/galleries (direct loading already works)

## Far-future roadmap

- [ ] Render comments on Imgur posts
- [ ] Frontpage imgur url form

## Installation

### Docker

Using Docker and docker-compose:

```
docker-compose build
docker-compose up
```

Omgur should now be running at http://localhost:8080.

Prebuilt images are also available at `registry.geraldwu.com/gerald/omgur:latest`. See https://git.geraldwu.com/gerald/omgur/container_registry for all images and tags.

### Manual

1. Install [Golang](https://golang.org/).
1. (Optional) Install [redis-server](https://redis.io/).
Caches images from imgur – highly recommended.
1. Clone and set up the repository.
```
git clone https://git.geraldwu.com/gerald/omgur
cd omgur
go mod init git.geraldwu.com/gerald/omgur
go mod tidy
go build -v -a ./cmd/omgur
export REDIS_HOST=localhost
redis-server
./omgur
```

Omgur should now be running at http://localhost:8080.
