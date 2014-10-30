# Mig31

[![Build Status](https://travis-ci.org/nfisher/mig31.svg?branch=master)](https://travis-ci.org/nfisher/mig31)

![Mig 31 Fighter](http://upload.wikimedia.org/wikipedia/commons/thumb/7/7d/Russian_Air_Force_MiG-31_inflight_Pichugin.jpg/300px-Russian_Air_Force_MiG-31_inflight_Pichugin.jpg)

## Overview

Mig31 is a command line tool to aid in the management of Cassandra schemas using CQL3. It's primary difference from other migration tools is that it allows for the specification of placement strategy and strategy options on a per environment basis. The goal in using this tool is to encourage tracking all schema changes with the service they are linked to rather than manually editing in each environment via cqlsh or cassandra-cli.

The following functionality is currently enabled:

- init - initialises a migration meta data table.
- dryrun - outputs what would be run against the database based on the current state.
- offline - outputs the full set of schema changes for the target enviroment.

Sample Commands;

```
mig31 -env=dev -init   # initialise migration metadata
mig31 -env=dev -dryrun # output the proposed schema based on what migrations have been run previously.
mig31 -env=dev -dryrun | cqlsh 10.2.0.21  # apply the proprosed schema
```

Sample Config (config.json: default):

```
{
  "environments": [
      {
        "name": "dev",
        "cluster":"192.168.33.10",
        "confirmisoptional": true,
        "keyspace": "release",
        "placement": { "strategy": "SimpleStrategy", "options": "'replication_factor': 1" }
      },
      {
        "name": "lve-prem",
        "cluster": "lve-pre.local",
        "keyspace": "release",
        "placement": { "strategy": "NetworkTopologyStrategy", "options": "eu-west-1 : 3, us-east-1 : 3, ap-northeast-1 : 3" }
      }
    ]
}
```

Sample Config (config.xml:  -config=config.xml):

```
<migrations>
  <environments>
    <environment name="dev" host="10.2.0.21" keyspace="release">
      <confirmisoptional>true</confirmisoptional>
      <placement strategy="SimpleStrategy" options="'replication_factor':1"/>
    </environment>
    <environment name="lve-prem" host="lve-prem.local" keyspace="release">
      <placement strategy="NetworkTopologyStrategy" options="eu-west-1 : 3, us-east-1 : 3, ap-northeast-1 : 3"/>
    </environment>
  </environments>
</migrations>
```

Sample Migration:

```
-- @up

CREATE KEYSPACE "release"
  WITH REPLICATION = {'class': '$${placement_strategy}', $${strategy_options}}
  AND DURABLE_WRITES = true;

CREATE TABLE release.release (keyspace_name TEXT PRIMARY KEY, ticketNumber INT, nextTicketNumber INT) 
  WITH COMPACT STORAGE AND compaction={'class': 'SizeTieredCompactionStrategy'} 
  AND compression={'sstable_compression': 'SnappyCompressor'};

-- @down

drop column family release;
drop keyspace release;
```

Improvements:

- eliminate the use of cqlsh by creating a lexical parser that can hand off individual queries to gocql without error.
- report identity.
- checksum the schema after each migration so manual changes can be identified.
- snapshotting that could collapse a freshly migrated schema into one migration.

