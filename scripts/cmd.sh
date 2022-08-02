docker-compose --env-file test_influxdb.env up

curl -i -X GET localhost:8080/power-consumptions

curl -i -X POST localhost:8080/power-consumptions 
-H "Content-Type: application/json" -d '{"client_id": "656565", "location": "BANGLADESH", "current_load_in_amperes": 0.3}'