## This Makefile is used to publish skopeo container images with Travis CI ##
## Environment variables, used in this Makefile are specified in .travis.yml

export DOCKER_CLI_EXPERIMENTAL=enabled
GOARCH ?= $(shell go env GOARCH)

# Build container image of skopeo upstream based on host architecture
build-image/upstream:
	docker build -t "${UPSTREAM_IMAGE}-${GOARCH}" contrib/skopeoimage/upstream

# Build container image of skopeo stable based on host architecture
build-image/stable:
	docker build -t "${STABLE_IMAGE}-${GOARCH}" contrib/skopeoimage/stable

# Push container image of skopeo upstream (based on host architecture) to image repository
push-image/upstream:
	echo "${QUAY_PASSWORD}" | docker login quay.io -u "${QUAY_USERNAME}" --password-stdin
	docker push "${UPSTREAM_IMAGE}-${GOARCH}"

# Push container image of skopeo stable (based on host architecture) to image default and extra repositories
push-image/stable:
	echo "${QUAY_PASSWORD}" | docker login quay.io -u "${QUAY_USERNAME}" --password-stdin
	docker push "${STABLE_IMAGE}-${GOARCH}"
	docker tag "${STABLE_IMAGE}-${GOARCH}" "${EXTRA_STABLE_IMAGE}-${GOARCH}"
	docker push "${EXTRA_STABLE_IMAGE}-${GOARCH}"

# Create and push multiarch image manifest of skopeo upstream
push-manifest-multiarch/upstream:
	echo "${QUAY_PASSWORD}" | docker login quay.io -u "${QUAY_USERNAME}" --password-stdin
	docker manifest create "${UPSTREAM_IMAGE}" $(foreach arch,${MULTIARCH_MANIFEST_ARCHITECTURES}, ${UPSTREAM_IMAGE}-${arch})
	docker manifest push --purge "${UPSTREAM_IMAGE}"

# Create and push multiarch image manifest of skopeo stable
push-manifest-multiarch/stable:
	echo "${QUAY_PASSWORD}" | docker login quay.io -u "${QUAY_USERNAME}" --password-stdin
	docker manifest create "${STABLE_IMAGE}" $(foreach arch,${MULTIARCH_MANIFEST_ARCHITECTURES}, ${STABLE_IMAGE}-${arch})
	docker manifest push --purge "${STABLE_IMAGE}"
	# Push to extra repository
	docker manifest create "${EXTRA_STABLE_IMAGE}" $(foreach arch,${MULTIARCH_MANIFEST_ARCHITECTURES}, ${EXTRA_STABLE_IMAGE}-${arch})
	docker manifest push --purge "${EXTRA_STABLE_IMAGE}"
