## Overview

This is a command line tool to provide simple schema migrations for cassandra.

Sample Commands;

```
mig31 -env=lve-prem -config=config.xml -dryrun
mig31 -env=dev -config=config.xml -verbose
mig31 -env=dev -config=config.xml -verbose
```

Sample Config:

```
<migrations>
<environments>
  <environment name="dev" host="10.2.0.21">
    <confirmisoptional>true</confirmisoptional>
    <placement strategy="SimpleStrategy" options="{replication_factor:1}"/>
  </environment>
  <environment name="lve-prem" host="lve-prem.local">
    <placement strategy="NetworkTopologyStrategy" options="{eu-west-1 : 3, us-east-1 : 3, ap-northeast-1 : 3}"/>
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

use release;

create column family release
	AND comparator = UTF8Type
	AND key_validation_class = UTF8Type
	AND default_validation_class = UTF8Type
	AND column_metadata = [
		{column_name: created, validation_class: DateType}
		{column_name: updated, validation_class: DateType}
	];
create column family envIndex and comparator = 'UTF8Type';

-- @down

drop column family release;
drop keyspace release;
```
