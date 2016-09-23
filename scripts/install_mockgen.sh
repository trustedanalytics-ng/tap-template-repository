#!/bin/bash
# Copyright (c) 2016 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

if [ -f $GOPATH/bin/mockgen ]; then
  echo "mockgen is already built"
  exit 0
fi;
echo "building mockgen..."
mkdir -p temp/src/github.com/golang/mock
cp -sfR `pwd`/vendor/github.com/golang/mock/. temp/src/github.com/golang/mock
cd temp
OLD_GOPATH=$GOPATH
GOPATH=`pwd`
go build github.com/golang/mock/mockgen
mv mockgen $OLD_GOPATH/bin/
GOPATH=$OLD_GOPATH
echo "mockgen built"
cd ..
rm -rf temp
