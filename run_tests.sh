 #!/bin/bash

 # this is a temporary hack to be used until we create more complete
 # *_test.go files 

 # first run the tests we do have (may want to add coverage flags later)
 cd sdf
 go test
 cd -
 
 # then run all the examples (this takes a while)
 time for d in examples/*
 do 
    echo $d
    cd $d
    make || exit 1
    ./$(basename $d) || exit 1
    make clean
    cd -
 done
