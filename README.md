## Overview

This is a command line tool to provide simple schema migrations for cassandra.  Currently only init, dryrun and offline modes are supported.  This requires an interrim dependency on cqlsh but this could be eliminated in the future with a CQL lexical parser.

Sample Commands;

```
mig31 -env=dev -init   # initialise migration metadata
mig31 -env=dev -dryrun # output the proposed schema based on what migrations have been run previously.
mig31 -env=dev -dryrun | cqlsh 10.2.0.21  # apply the proprosed schema
```

Sample Config (config.xml):

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

create keyspace release
  with placement_strategy = $${placement_strategy}
  and strategy_options = $${placement_options}
  and durable_writes = true;

CREATE TABLE release.release (keyspace_name TEXT PRIMARY KEY, ticketNumber INT, nextTicketNumber INT) 
  WITH COMPACT STORAGE AND compaction={'class': 'SizeTieredCompactionStrategy'} 
  AND compression={'sstable_compression': 'SnappyCompressor'};

-- @down

drop column family release;
drop keyspace release;
```
