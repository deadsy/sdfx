#!/bin/bash
# Cross-branch STL regression check.
#
# Renders every example/*/ at BASE and HEAD revs, then compares the STL
# outputs. Reports per-file: IDENTICAL (canonical hash match), MINOR (tiny
# float drift), or MATERIAL (geometry genuinely changed). Intended as a
# sanity check that a core-library change didn't inadvertently alter
# unrelated examples.
#
# Usage:  run.sh [BASE_REF] [HEAD_REF]   (defaults: master, HEAD)

set -u

BASE_REF="${1:-master}"
HEAD_REF="${2:-HEAD}"

REPO_ROOT=$(cd "$(dirname "$0")/../.." && pwd)
TOOL_DIR="$REPO_ROOT/tools/stldiff"
BASE_TREE=$(mktemp -d -t sdfx-stldiff-base-XXXXXX)
HEAD_TREE=$(mktemp -d -t sdfx-stldiff-head-XXXXXX)
LOG=$(mktemp -t sdfx-stldiff-results-XXXXXX.txt)
TIMEOUT=120

cleanup() {
  git -C "$REPO_ROOT" worktree remove --force "$BASE_TREE" 2>/dev/null
  git -C "$REPO_ROOT" worktree remove --force "$HEAD_TREE" 2>/dev/null
}
trap cleanup EXIT

git -C "$REPO_ROOT" worktree add --detach "$BASE_TREE" "$BASE_REF" >/dev/null
git -C "$REPO_ROOT" worktree add --detach "$HEAD_TREE" "$HEAD_REF" >/dev/null

( cd "$TOOL_DIR" && go build -o "$TOOL_DIR/stldiff" . )

run_with_timeout() {
  local d=$1 n=$2
  ( cd "$d" && rm -f *.stl && make all >/dev/null 2>&1 && ./"$n" >/dev/null 2>&1 ) &
  local pid=$!
  ( sleep $TIMEOUT && kill -9 $pid 2>/dev/null ) &
  local watcher=$!
  wait $pid 2>/dev/null
  local rc=$?
  kill -9 $watcher 2>/dev/null
  wait $watcher 2>/dev/null
  return $rc
}

: > "$LOG"

for dir in "$HEAD_TREE"/examples/*/; do
  name=$(basename "$dir")
  basedir="$BASE_TREE/examples/$name"
  headdir="$HEAD_TREE/examples/$name"

  [ -d "$basedir" ] || { echo "[$name] SKIP (not in base)" | tee -a "$LOG"; continue; }
  [ -f "$basedir/main.go" ] || { echo "[$name] SKIP (no main.go)" | tee -a "$LOG"; continue; }

  run_with_timeout "$basedir" "$name"; bstatus=$?
  run_with_timeout "$headdir" "$name"; hstatus=$?
  if [ $bstatus -ne 0 ] || [ $hstatus -ne 0 ]; then
    echo "[$name] BUILD/RUN FAILED (base=$bstatus head=$hstatus)" | tee -a "$LOG"
    continue
  fi

  any=0
  for bstl in "$basedir"/*.stl; do
    [ -f "$bstl" ] || continue
    any=1
    fname=$(basename "$bstl")
    hstl="$headdir/$fname"
    if [ ! -f "$hstl" ]; then
      echo "[$name/$fname] MISSING on head" | tee -a "$LOG"
      continue
    fi
    result=$("$TOOL_DIR/stldiff" "$bstl" "$hstl" 2>&1)
    echo "[$name/$fname] $result" | tee -a "$LOG"
  done
  [ $any -eq 0 ] && echo "[$name] no STL produced" | tee -a "$LOG"
done

echo "---" | tee -a "$LOG"
echo "Summary (base=$BASE_REF head=$HEAD_REF):" | tee -a "$LOG"
echo "  IDENTICAL: $(grep -c IDENTICAL "$LOG")" | tee -a "$LOG"
echo "  MINOR:     $(grep -c 'MINOR ' "$LOG")" | tee -a "$LOG"
echo "  MATERIAL:  $(grep -c MATERIAL "$LOG")" | tee -a "$LOG"
echo "  Results log: $LOG"
