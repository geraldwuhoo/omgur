# omgur

Omgur is a free and open-source alternative Imgur front-end focused on privacy.

Inspired by the [Invidious](https://github.com/iv-org/invidious), [Nitter](https://github.com/zedeus/nitter), and [Teddit](https://github.com/teddit-net/teddit) projects.

- No JavaScript or ads
- All requests go through the backend, client never talks to Imgur
- Prevents Imgur from tracking your IP or JavaScript fingerprint
- Lightweight
- Self-hostable

## Current Features

- Direct image loading via https://i.imgur.com/
- Imgur album loading via https://imgur.com/a/

## Roadmap

- Redis caching for images
- Imgur post loading via https://imgur.com/
- Imgur gallery loading via https://imgur.com/gallery

## Far-future roadmap

- Render comments on Imgur posts
- Imgur frontpage loading

## Installation

### Docker

Using Docker and docker-compose:

```
docker-compose build
docker-compose up
```

Omgur should now be running at http://localhost:8080.

### Manual

1. Install Golang.
2. Clone and set up the repository.
```
git clone https://git.geraldwu.com/gerald/omgur
cd omgur
go mod init git.geraldwu.com/gerald/omgur
go mod tidy
go build -v -a ./cmd/omgur
./omgur
```

Omgur should now be running at http://localhost:8080.
