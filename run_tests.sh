 #!/bin/bash
 
 nice time for d in examples/*
 do 
    echo $d
    cd $d
    make || break
    ./$(basename $d) || break
    make clean
    cd -
 done