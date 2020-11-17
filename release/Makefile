## This Makefile is used to publish skopeo container images with Travis CI ##
## Environment variables, used in this Makefile are specified in .travis.yml

export DOCKER_CLI_EXPERIMENTAL=enabled
GOARCH ?= $(shell go env GOARCH)

# Dereference variable $(1), return value if non-empty, otherwise raise an error.
err_if_empty = $(if $(strip $($(1))),$(strip $($(1))),$(error Required $(1) variable is undefined or empty))

# Requires two arguments: Names of the username and the password env. vars.
define quay_login
	@echo "$(call err_if_empty,$(2))" | \
		docker login quay.io -u "$(call err_if_empty,$(1))" --password-stdin
endef

# Build container image of skopeo upstream based on host architecture
build-image/upstream:
	docker build -t "${UPSTREAM_IMAGE}-${GOARCH}" contrib/skopeoimage/upstream

# Build container image of skopeo stable based on host architecture
build-image/stable:
	docker build -t "${STABLE_IMAGE}-${GOARCH}" contrib/skopeoimage/stable

# Push container image of skopeo upstream (based on host architecture) to image repository
push-image/upstream:
	$(call quay_login,SKOPEO_QUAY_USERNAME,SKOPEO_QUAY_PASSWORD)
	docker push "${UPSTREAM_IMAGE}-${GOARCH}"

# Push container image of skopeo stable (based on host architecture) to image default and extra repositories
push-image/stable:
	$(call quay_login,SKOPEO_QUAY_USERNAME,SKOPEO_QUAY_PASSWORD)
	docker push "${STABLE_IMAGE}-${GOARCH}"
	docker tag "${STABLE_IMAGE}-${GOARCH}" "${EXTRA_STABLE_IMAGE}-${GOARCH}"
	$(call quay_login,CONTAINERS_QUAY_USERNAME,CONTAINERS_QUAY_PASSWORD)
	docker push "${EXTRA_STABLE_IMAGE}-${GOARCH}"

# Create and push multiarch image manifest of skopeo upstream
push-manifest-multiarch/upstream:
	docker manifest create "${UPSTREAM_IMAGE}" $(foreach arch,${MULTIARCH_MANIFEST_ARCHITECTURES}, ${UPSTREAM_IMAGE}-${arch})
	$(call quay_login,SKOPEO_QUAY_USERNAME,SKOPEO_QUAY_PASSWORD)
	docker manifest push --purge "${UPSTREAM_IMAGE}"

# Create and push multiarch image manifest of skopeo stable
push-manifest-multiarch/stable:
	docker manifest create "${STABLE_IMAGE}" $(foreach arch,${MULTIARCH_MANIFEST_ARCHITECTURES}, ${STABLE_IMAGE}-${arch})
	$(call quay_login,SKOPEO_QUAY_USERNAME,SKOPEO_QUAY_PASSWORD)
	docker manifest push --purge "${STABLE_IMAGE}"
	# Push to extra repository
	docker manifest create "${EXTRA_STABLE_IMAGE}" $(foreach arch,${MULTIARCH_MANIFEST_ARCHITECTURES}, ${EXTRA_STABLE_IMAGE}-${arch})
	$(call quay_login,CONTAINERS_QUAY_USERNAME,CONTAINERS_QUAY_PASSWORD)
	docker manifest push --purge "${EXTRA_STABLE_IMAGE}"
