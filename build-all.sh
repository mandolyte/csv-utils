#!/bin/sh 
echo build all
for i in `cat build-list.txt`
do
echo Working on $i in `pwd`
cd $i 
go install $i.go 
cd ..
done
