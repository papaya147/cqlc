CREATE KEYSPACE data;
create TABLE data.trips(
    emuserid uuid,
    vin text,
    start timestamp,
end timestamp,
startlat double,
startlng double,
endlat double,
endlng double,
startloc text,
endloc text,
distance double,
averagepas double,
averagevehiclespeed double,
averagemotorspeed double,
maxspeed int,
minsoc int,
maxsoc int,
averageheadlamp double,
averagebrake double,
elevationgain double,
minelevation double,
maxelevation double,
calories double,
duration bigint,
effort double,
snacount bigint,
ecocount bigint,
urbancount bigint,
fitnesscount bigint,
labcount bigint,
childcount bigint,
records int,
PRIMARY KEY((emuserid, vin), start)
) WITH CUSTOM_PROPERTIES = { 'capacity_mode': { 'throughput_mode' :'PAY_PER_REQUEST' },
'encryption_specification': { 'encryption_type' :'AWS_OWNED_KMS_KEY' } }
AND CLUSTERING
ORDER BY(start ASC)