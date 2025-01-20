#!/bin/sh
<<<<<<< HEAD
bazel build //proto/cel/expr/conformance/...
files=($(bazel aquery 'kind(proto, //proto/cel/expr/conformance/...)' | grep Outputs | grep "[.]pb[.]go" | sed 's/Outputs: \[//' | sed 's/\]//' | tr "," "\n"))
for src in ${files[@]};
do
  dst=$(echo $src | sed 's/\(.*\/cel.dev\/expr\/\(.*\)\)/\2/')
=======
bazel build //proto/test/...
files=($(bazel aquery 'kind(proto, //proto/...)' | grep Outputs | grep "[.]pb[.]go" | sed 's/Outputs: \[//' | sed 's/\]//' | tr "," "\n"))
for src in ${files[@]};
do
  dst=$(echo $src | sed 's/\(.*\%\/github.com\/google\/cel-spec\/\(.*\)\)/\2/')
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
  echo "copying $dst"
  $(cp $src $dst)
done
