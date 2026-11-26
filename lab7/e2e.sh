#!/usr/bin/env bash
# Lab #7 - end-to-end check against a running gorse-in-one built from this fork.
#
#   lab7/e2e.sh seed  [BASE_URL]            insert fixed items/feedback and wait for the leaderboard
#   lab7/e2e.sh check [BASE_URL] [CONTRACT] run the HTTP checks (all, or only negative_n with ONLY=negative_n)
#
# CONTRACT defaults to the content of lab7/claimed-contract.txt ("legacy" or "bounded").
# Each observation is printed as "LAB7 E2E ..." and appended as JSON to $LAB7_RESULTS (if set).
set -uo pipefail

cmd="${1:?usage: e2e.sh seed|check [BASE_URL] [CONTRACT]}"
BASE="${2:-http://127.0.0.1:8088}"
HERE="$(cd "$(dirname "$0")" && pwd)"
CONTRACT="${3:-$(tr -d '[:space:]' < "$HERE/claimed-contract.txt")}"
ONLY="${ONLY:-all}"
LEADERBOARD="most_starred"
N_ITEMS=15

item() { printf 'item-%02d' "$1"; }

seed() {
  local ts items fb i u
  ts="$(date -u -d '-1 day' +%Y-%m-%dT%H:%M:%SZ)"
  # 15 items; item-i receives (15 - i) stars, so the leaderboard is item-00 (15) ... item-14 (1).
  items='['
  for i in $(seq 0 $((N_ITEMS - 1))); do
    items+="{\"ItemId\":\"$(item "$i")\",\"IsHidden\":false,\"Categories\":[],\"Timestamp\":\"$ts\",\"Labels\":{},\"Comment\":\"\"},"
  done
  items="${items%,}]"
  fb='['
  for i in $(seq 0 $((N_ITEMS - 1))); do
    for u in $(seq 0 $((N_ITEMS - 1 - i))); do
      fb+="{\"FeedbackType\":\"star\",\"UserId\":\"fan-$u\",\"ItemId\":\"$(item "$i")\",\"Timestamp\":\"$ts\"},"
    done
  done
  # alice has read item-01: it must be excluded from her leaderboard view.
  fb+="{\"FeedbackType\":\"read\",\"UserId\":\"alice\",\"ItemId\":\"item-01\",\"Timestamp\":\"$ts\"}]"

  echo "insert items:    $(curl -fsS -X POST -H 'Content-Type: application/json' -d "$items" "$BASE/api/items" | jq -c .)" || return 1
  echo "insert feedback: $(curl -fsS -X POST -H 'Content-Type: application/json' -d "$fb" "$BASE/api/feedback" | jq -c .)" || return 1

  # The master recomputes the leaderboard in its task loop (fit_period = 10s).
  local start count
  start=$(date +%s)
  for _ in $(seq 1 60); do
    count=$(curl -fsS "$BASE/api/non-personalized/$LEADERBOARD?n=100" | jq 'length' 2>/dev/null || echo 0)
    if [ "$count" = "$N_ITEMS" ]; then
      echo "leaderboard ready: $count items after $(( $(date +%s) - start ))s"
      return 0
    fi
    sleep 2
  done
  echo "leaderboard not ready after $(( $(date +%s) - start ))s (count=$count)" >&2
  return 1
}

FAILURES=0
# observe NAME QUERY -> sets STATUS, BODY, IDS (comma-separated ids, or empty)
observe() {
  local name="$1" query="$2" out
  out="$(mktemp)"
  STATUS=$(curl -sS -o "$out" -w '%{http_code}' "$BASE/api/non-personalized/$LEADERBOARD$query")
  BODY="$(cat "$out")"
  rm -f "$out"
  if [ "$STATUS" = "200" ]; then
    IDS=$(jq -r '[.[].Id] | join(",")' <<<"$BODY")
    COUNT=$(jq 'length' <<<"$BODY")
  else
    IDS=""
    COUNT=0
  fi
  local line
  line=$(jq -cn --arg case "$name" --arg query "$query" --argjson status "$STATUS" --argjson count "$COUNT" \
    --arg ids "$IDS" --arg body "$( [ "$STATUS" = 200 ] || echo "$BODY" )" --arg contract "$CONTRACT" \
    '{case:$case, query:$query, status:$status, count:$count, ids:$ids, body:$body, contract:$contract}')
  echo "LAB7 E2E $line"
  if [ -n "${LAB7_RESULTS:-}" ]; then echo "$line" >> "$LAB7_RESULTS"; fi
}

expect() { # expect DESCRIPTION ACTUAL EXPECTED
  if [ "$2" = "$3" ]; then
    echo "  PASS $1"
  else
    echo "  FAIL $1: expected '$3', got '$2'"
    FAILURES=$((FAILURES + 1))
  fi
}

check() {
  echo "contract under check: $CONTRACT (checks: $ONLY)"
  if [ "$ONLY" = "all" ]; then
    # ---- invariants: must hold before and after the Lab #7 change ----
    observe top_n "?n=3"
    expect "top_n status" "$STATUS" 200
    expect "top_n ids" "$IDS" "item-00,item-01,item-02"

    observe default_n ""
    expect "default_n status" "$STATUS" 200
    expect "default_n ids (server.default_n=10)" "$IDS" "$(for i in $(seq 0 9); do printf '%s,' "$(item "$i")"; done | sed 's/,$//')"

    observe exclude_read_items "?n=3&user-id=alice"
    expect "exclude_read_items status" "$STATUS" 200
    expect "exclude_read_items ids" "$IDS" "item-00,item-02,item-03"

    observe non_numeric_n_rejected "?n=abc"
    expect "non_numeric_n_rejected status" "$STATUS" 400
  fi

  # ---- the intentionally changed case ----
  observe negative_n "?n=-1"
  case "$CONTRACT" in
    legacy)
      expect "negative_n status (legacy)" "$STATUS" 200
      expect "negative_n count (legacy: whole leaderboard)" "$COUNT" "$N_ITEMS"
      ;;
    bounded)
      expect "negative_n status (bounded)" "$STATUS" 400
      expect "negative_n error message (bounded)" "$(grep -c 'invalid n parameter' <<<"$BODY")" 1
      ;;
    *)
      echo "unknown contract $CONTRACT"; FAILURES=$((FAILURES + 1)) ;;
  esac

  echo "failures: $FAILURES"
  [ "$FAILURES" -eq 0 ]
}

case "$cmd" in
  seed) seed ;;
  check) check ;;
  *) echo "unknown command $cmd" >&2; exit 2 ;;
esac
