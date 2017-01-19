#!/bin/sh
set -e

ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --compilers=2
