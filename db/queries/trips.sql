-- name: GetLatestTripStart :one
SELECT toUnixTimestamp(start) as start
FROM data.trips
WHERE emuserid = ?
    AND vin IN (?)
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
    AND vin IN (?)
ORDER BY start DESC
LIMIT 1;