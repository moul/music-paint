GOPKG ?=	moul.io/music-paint
DOCKER_IMAGE ?=	moul/music-paint
GOBINS ?=	.
NPM_PACKAGES ?=	.

all: test install

include rules.mk
