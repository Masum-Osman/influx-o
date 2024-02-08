package main

// Define the start and end time for the query
start_time = now() - 30d
end_time = now()

// Fetch data from classes measurement
offer_live_class_data = from(bucket: "your_bucket_name")
    |> range(start: start_time, stop: end_time)
    |> filter(fn: (r) => r["_measurement"] == "classes")

// Fetch data from class_participants measurement
live_class_participants_data = from(bucket: "your_bucket_name")
    |> range(start: start_time, stop: end_time)
    |> filter(fn: (r) => r["_measurement"] == "class_participants")

// Fetch data from quiz_participants measurement
user_quiz_participants_data = from(bucket: "your_bucket_name")
    |> range(start: start_time, stop: end_time)
    |> filter(fn: (r) => r["_measurement"] == "quiz_participants")

// Fetch data from exam_participants measurement
user_exam_participants_data = from(bucket: "your_bucket_name")
    |> range(start: start_time, stop: end_time)
    |> filter(fn: (r) => r["_measurement"] == "exam_participants")

// Join data from all measurements for the given user_id
joined_data = join(
    tables: {
        offer_live_class: offer_live_class_data,
        live_class_participants: live_class_participants_data,
        user_quiz_participants: user_quiz_participants_data,
        user_exam_participants: user_exam_participants_data,
    },
    on: ["_time", "user_id"]
)

// Calculate total time by summing duration from all measurements
total_time = joined_data
    |> filter(fn: (r) => r["user_id"] == "65b0d02d60023489a81394db")
    |> group(columns: ["_time"])
    |> aggregateWindow(every: 7d, fn: sum, column: "duration")

total_time
