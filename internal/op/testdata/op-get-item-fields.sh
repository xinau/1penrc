#!/bin/sh

case "$3" in
  "item-fields-well-formed")
    fields='{"foo":"bar", "baz":"qux"}'
  ;;
  "item-fields-malformed")
    fields='invalid item fields'
  ;;
esac

printf "%s\n" "$fields"
exit 0
