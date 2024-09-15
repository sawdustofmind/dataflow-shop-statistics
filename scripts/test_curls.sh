#/bin/bash

### user & public api
USER_URL=http://0.0.0.0:8080

### data
curl -XPOST ${USER_URL}/v1/data \
-H 'Content-Type: application/json' \
--data-binary @- << EOF
{
  "product_id": "12345",
  "store_id": "6789",
  "quantity_sold": 10,
  "sale_price": 19.99,
  "sale_date": "2024-06-15T14:30:00Z"
}
EOF

curl ${USER_URL}/v1/data

### statistics
curl -XPOST ${USER_URL}/v1/calculate \
-H 'Content-Type: application/json' \
--data-binary @- << EOF
{
  "operation": "total_sales",
  "store_id": "6789",
  "start_date": "2024-06-15T14:30:00Z",
  "end_date": "2024-06-15T14:30:00Z"
}
EOF
