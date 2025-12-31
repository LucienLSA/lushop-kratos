#!/usr/bin/env bash
set -euo pipefail

# prune-crpi-images.sh
# Dry-run and delete local docker images whose repository contains a given prefix.
#
# Usage:
#  ./prune-crpi-images.sh           -> dry-run with default pattern
#  ./prune-crpi-images.sh --delete  -> perform deletion
#  ./prune-crpi-images.sh <pattern> [--delete]
#
PATTERN_DEFAULT="crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com"
PATTERN="${1:-$PATTERN_DEFAULT}"
DELETE=false

if [ "${1:-}" = "--delete" ]; then
  PATTERN="$PATTERN_DEFAULT"
  DELETE=true
elif [ "${2:-}" = "--delete" ]; then
  DELETE=true
fi

echo "Pattern: $PATTERN"
if [ "$DELETE" = true ]; then
  echo "Mode: DELETE (will remove matching image tags)"
else
  echo "Mode: DRY-RUN (no deletion). Pass --delete to remove."
fi

# Gather matching images: repository:tag and ID
mapfile -t matches < <(docker images --format '{{.Repository}}:{{.Tag}} {{.ID}}' | grep -- "$PATTERN" || true)

if [ "${#matches[@]}" -eq 0 ]; then
  echo "No images found matching pattern: $PATTERN"
  exit 0
fi

echo "Found ${#matches[@]} image tag(s) matching pattern."

# Prepare arrays
declare -a to_delete_tags
declare -A seen_ids

for line in "${matches[@]}"; do
  repo_tag="${line% *}"
  id="${line##* }"
  # skip dangling <none>:<none>
  if [[ "$repo_tag" == "<none>:<none>" ]]; then
    continue
  fi
  # If we've already recorded this exact repo:tag, skip
  if printf '%s\n' "${to_delete_tags[@]}" | grep -Fxq "$repo_tag"; then
    continue
  fi
  to_delete_tags+=("$repo_tag")
  seen_ids["$id"]=1
done

if [ "$DELETE" = false ]; then
  echo
  echo "DRY-RUN: the following tags would be removed:"
  printf '%s\n' "${to_delete_tags[@]}"
  echo
  echo "If you are satisfied, run with --delete to actually remove them."
  exit 0
fi

echo
echo "Deleting ${#to_delete_tags[@]} tag(s)..."
for tag in "${to_delete_tags[@]}"; do
  echo "Removing tag: $tag"
  if docker rmi "$tag"; then
    echo "Removed $tag"
  else
    echo "Failed to remove $tag (maybe in use)."
  fi
done

echo "Deletion complete. You may still have image IDs with no tags (dangling). To remove dangling images:"
echo "  docker image prune -f"


