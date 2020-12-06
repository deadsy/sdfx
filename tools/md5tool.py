#!/usr/bin/env python3
"""
Output md5 hashes for a set of matched files.
"""

import glob
import hashlib

def md5sum(filename):
  """return the md5 hash of a file"""
  with open(filename, mode='rb') as f:
    d = hashlib.md5()
    while True:
      buf = f.read(8192)
      if not buf:
        break
      d.update(buf)
    return d.hexdigest()

def main():
  """entry point"""
  f = []
  for ext in ("stl", "svg", "dxf", "png"):
    f.extend(glob.glob('*.%s' % ext))
  for fname in f:
    print("%s\t%s" % (md5sum(fname), fname))

main()
