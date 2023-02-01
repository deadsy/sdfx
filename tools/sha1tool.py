#!/usr/bin/env python3
"""
Output sha1 hashes for a set of matched files.
"""

import glob
import hashlib

def sha1sum(filename):
  """return the sha1 hash of a file"""
  with open(filename, mode='rb') as f:
    d = hashlib.sha1()
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
    print("%s  %s" % (sha1sum(fname), fname))

main()
