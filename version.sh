#!/bin/bash

git name-rev --name-only HEAD >gitversion
git log -1 --pretty=%h >>gitversion
