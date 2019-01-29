#!/bin/bash

set -e

PROJECT_WORKSPACE="$(dirname $0)/.."
PROJECT_WORKSPACE="$(cd $PROJECT_WORKSPACE; pwd)"

DOCKER_ORGANISATION="hypersclae"

LATEST="false"

echo -e "Building Hypercloud projects...\n"

echo "Config:"
if [ -z "$TRAVIS_TAG" ]; then
    CI_BUILD_VERSION=$(git describe --match 'v[0-9]*' --dirty='-dev' --always)
else
    CI_BUILD_VERSION="${TRAVIS_TAG#v}"
    DOCKER_TAG="$CI_BUILD_VERSION"
    LATEST="true"
fi

if [ -z "$TRAVIS_COMMIT" ]; then
    CI_BUILD_COMMIT=$(git rev-parse HEAD)
else
    CI_BUILD_COMMIT="$TRAVIS_COMMIT"
fi

CI_BUILD_URL=$(git config --get remote.origin.url)
CI_BUILD_DATE=$(date +%Y-%m-%dT%T%z)

if [ -z "$TRAVIS_BRANCH" ]; then
    CI_BUILD_BRANCH=$(git rev-parse --abbrev-ref HEAD)
else
    CI_BUILD_BRANCH="$TRAVIS_BRANCH"
fi

if [ "$CI_BUILD_BRANCH" == "develop" ]; then
    DOCKER_TAG="dev"
fi

echo "  Docker Tag: $DOCKER_TAG"
echo "  Version: $CI_BUILD_VERSION"
echo "  VCS URL: $CI_BUILD_URL"
echo "  VCS Ref: $CI_BUILD_COMMIT"
echo "  VCS Branch: $CI_BUILD_BRANCH"
echo "  Build Date: $CI_BUILD_DATE"
echo "  Workspace: $PROJECT_WORKSPACE"
echo ""

for image in $(find $PROJECT_WORKSPACE/cmd -name "Dockerfile"); do
    DOCKER_REPO=${image#$PROJECT_WORKSPACE/cmd/}
    DOCKER_REPO=${DOCKER_REPO%/Dockerfile}

    echo "Building $DOCKER_ORGANISATION/$DOCKER_REPO..."
    docker build --rm \
        --cache-from "$DOCKER_ORGANISATION/$DOCKER_REPO:latest" \
        --build-arg "VERSION=$CI_BUILD_VERSION" \
        --build-arg "VCS_URL=$CI_BUILD_URL" \
        --build-arg "VCS_REF=$CI_BUILD_COMMIT" \
        --build-arg "BUILD_DATE=$CI_BUILD_DATE" \
        -f "$PROJECT_WORKSPACE/cmd/$DOCKER_REPO/Dockerfile" \
        -t "$DOCKER_ORGANISATION/$DOCKER_REPO:$DOCKER_TAG" \
        "$PROJECT_WORKSPACE"

    # tagging latest only master branch
    if [ "$LATEST" == "true" ]; then
        echo "Tagging $DOCKER_ORGANISATION/$DOCKER_REPO:latest.."
        docker tag "$DOCKER_ORGANISATION/$DOCKER_REPO:$DOCKER_TAG" "$DOCKER_ORGANISATION/$DOCKER_REPO:latest"
    fi

    # pushing only in CI mode
    if [ "$CI" == "true" ]; then
        #docker login -u $DOCKER_USER -p $DOCKER_PASS
        echo "Pushing $DOCKER_ORGANISATION/$DOCKER_REPO:$DOCKER_TAG..."
        docker push "$DOCKER_ORGANISATION/$DOCKER_REPO:$DOCKER_TAG"

        if [ "$LATEST" == "true" ]; then
            echo "Pushing $DOCKER_ORGANISATION/$DOCKER_REPO:latest..."
            docker push "$DOCKER_ORGANISATION/$DOCKER_REPO:latest"
        fi
    fi

done

