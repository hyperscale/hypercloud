#!/bin/bash

set -e

usage() {
    echo "Usage: $1 TAG_NAME"
    exit 1
}

PROJECT_WORKSPACE="$(dirname $0)/.."
PROJECT_WORKSPACE="$(cd $PROJECT_WORKSPACE; pwd)"

DOCKER_ORGANISATION="hypersclae"

echo -e "Building HyperPaaS projects...\n"

echo "Config:"
if [ -z "$1" ]; then
    usage "$0"
else
    #DOCKER_TAG="$1"
    CI_BUILD_VERSION="${1#v}"
    DOCKER_TAG="$CI_BUILD_VERSION"
fi

if [ -z "$TRAVIS_COMMIT"]; then
    CI_BUILD_COMMIT=$(git rev-parse HEAD)
else
    CI_BUILD_COMMIT="$TRAVIS_COMMIT"
fi

CI_BUILD_URL=$(git config --get remote.origin.url)
CI_BUILD_DATE=$(date +%Y-%m-%dT%T%z)

echo "  Docker Tag: $DOCKER_TAG"
echo "  Version: $CI_BUILD_VERSION"
echo "  VCS URL: $CI_BUILD_URL"
echo "  VCS Ref: $CI_BUILD_COMMIT"
echo "  Build Date: $CI_BUILD_DATE"
echo "  Workspace: $PROJECT_WORKSPACE"
echo ""

for image in $(find $PROJECT_WORKSPACE/cmd -name "Dockerfile"); do
    DOCKER_REPO=${image#$PROJECT_WORKSPACE/cmd/}
    DOCKER_REPO=${DOCKER_REPO%/Dockerfile}

    echo "Building $DOCKER_ORGANISATION/$DOCKER_REPO..."
    docker build --rm \
        --build-arg "VERSION=$CI_BUILD_VERSION" \
        --build-arg "VCS_URL=$CI_BUILD_URL" \
        --build-arg "VCS_REF=$CI_BUILD_COMMIT" \
        --build-arg "BUILD_DATE=$CI_BUILD_DATE" \
        -f "$PROJECT_WORKSPACE/cmd/$DOCKER_REPO/Dockerfile" \
        -t "$DOCKER_ORGANISATION/$DOCKER_REPO:$DOCKER_TAG" \
        "$PROJECT_WORKSPACE"

    # tagging latest only master branch
    if [ "$TRAVIS_BRANCH" == "master" ]; then
        echo "Tagging $DOCKER_ORGANISATION/$DOCKER_REPO:$DOCKER_TAG in latest.."
        docker tag "$DOCKER_ORGANISATION/$DOCKER_REPO" "$DOCKER_ORGANISATION/$DOCKER_REPO:latest"
    fi

    # pushing only in CI mode
    if [ "$CI" == "true" ]; then
        echo "Pushing $DOCKER_ORGANISATION/$DOCKER_REPO:$DOCKER_TAG..."
        docker push "$DOCKER_ORGANISATION/$DOCKER_REPO:$DOCKER_TAG"

        if [ "$TRAVIS_BRANCH" == "master" ]; then
            echo "Pushing $DOCKER_ORGANISATION/$DOCKER_REPO:latest..."
            docker push "$DOCKER_ORGANISATION/$DOCKER_REPO:latest"
        fi
    fi

done

