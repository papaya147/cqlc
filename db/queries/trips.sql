-- name: GetLatestTripStart :one
SELECT toUnixTimestamp(start) as start
FROM data.trips
WHERE emuserid = ?
    AND vin IN (%s)
ORDER BY start DESC
LIMIT 1;
-- name: GetLatestTrip :one
SELECT toUnixTimestamp(start) as start,
    toUnixTimestamp(end
) as
end,
vin,
distance,
startlat,
startlng,
endlat,
endlng,
startloc,
endloc
FROM data.trips
WHERE emuserid = ?
    AND vin IN (%s)
ORDER BY start DESC
LIMIT 1;
-- name: UpdateTripLocations :exec
UPDATE data.trips
SET startloc = ?,
    endloc = ?
WHERE emuserid = ?
    AND vin = ?
    AND start = ?;
--- name: GetTripDetails :one
SELECT toUnixTimestamp(start) as start,
    toUnixTimestamp(end
) as
end,
vin,
distance,
startlat,
startlng,
endlat,
endlng,
startloc,
endloc,
averagepas,
averagevehiclespeed,
maxspeed,
minsoc,
maxsoc,
averageheadlamp,
averagebrake,
elevationgain,
calories,
duration,
effort,
snacount,
urbancount,
ecocount,
fitnesscount,
labcount,
childcount,
records
FROM data.trips
WHERE emuserid = ?
    AND vin = ?
    AND start = ?
LIMIT 1;
-- name: GetTrips :many
SELECT toUnixTimestamp(start) as start,
    toUnixTimestamp(end
) as
end,
vin,
distance,
startlat,
startlng,
endlat,
endlng,
startloc,
endloc
FROM data.trips
WHERE emuserid = ?
    AND vin IN (%s)
    AND start <= ?
    AND start >= ?
ORDER BY start DESC
LIMIT 8;