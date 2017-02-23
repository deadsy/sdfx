#!/bin/bash

# find and replace

old=( NewGearRack )
new=( GearRack2D )

total=${#old[*]}

for (( i=0; i<=$(( $total -1 )); i++ ))
do 
  oldname="${old[$i]}"
  newname="${new[$i]}"
  echo $oldname $newname
  git grep -lz $oldname | xargs -0 sed -i'' -e "s/$oldname/$newname/g"
done

