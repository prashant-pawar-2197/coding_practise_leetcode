#!/bin/bash

# Input JSON file
input_file="input.json"

# Read and publish JSON records to Kafka
while read -r line; do
  # Parse JSON using jq
  json_data=$(echo "$line" | jq -c .)

  # Publish JSON data to Kafka topic
  echo "$json_data" |  kubectl exec -it kafka-2 -n mtcil -- kafka-console-producer.sh --broker-list localhost:9092 --topic LOADERTEST_PR
done < "$input_file"
